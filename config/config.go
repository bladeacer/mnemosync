/*
Copyright © 2025 bladeacer <wg.nick.exe@gmail.com>
*/

package config

import (
	"os"
	"path/filepath"
	"gopkg.in/yaml.v3"
	"fmt"
)

type ConfigSchema struct {
	ConfigPath string `yaml:"config_path"`
	AppVersion string `yaml:"app_version"`
	IsInit bool `yaml:"is_init"`
	RepoPath string `yaml:"repo_path"`
	DbPath string `yaml:"db_path"`
}

type MnemoConf struct {
	ConfigSchema ConfigSchema `yaml:"config_schema"`
}

func GetMnemoConf() *MnemoConf {
	return &MnemoConf{
		ConfigSchema {
			ConfigPath: ResolveConfigPath(),
			AppVersion: "Version 0.0.1",
			IsInit: false,
			RepoPath: "",
			DbPath: ResolveDbPath(),
		},
	}
}

func ResolveConfigPath() string {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return ".config/mmsync/config.yaml" 
	}

	if envPath := os.Getenv("MMSYNC_CONF"); envPath != "" {
		return filepath.Join(homeDir, envPath)
	}

	return filepath.Join(homeDir, ".config/mmsync", "config.yaml")
}

func ResolveDbPath() string {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return ".config/mmsync/mmsync-state.json" 
	}

	if envPath := os.Getenv("MMSYNC_CONF"); envPath != "" {
		return filepath.Join(homeDir, envPath)
	}

	return filepath.Join(homeDir, ".config/mmsync", "mmsync-state.json")
}

func LoadConfig() (*MnemoConf, error) {
	configPath := ResolveConfigPath()
	
	defaultCfg := GetMnemoConf()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultCfg, nil
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	
	tempCfg := GetMnemoConf() 

	if err := yaml.Unmarshal(data, tempCfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML data. File may be invalid: %w", err)
	}
    
	warnings := healConfigSchema(tempCfg, defaultCfg)

	if len(warnings) > 0 {
		fmt.Fprintf(os.Stderr, "--- Configuration Healing Performed ---\n")
		for _, w := range warnings {
		    fmt.Fprintf(os.Stderr, "Config Warning: %v\n", w)
		}
		fmt.Fprintf(os.Stderr, "--- Saving Repaired Configuration ---\n\n")

		if saveErr := saveConfig(tempCfg, configPath); saveErr != nil {
		    return nil, fmt.Errorf("critical error: failed to save repaired configuration: %w", saveErr)
		}
	}
    
	return tempCfg, nil
}

func saveConfig(cfg *MnemoConf, targetPath string) error {
    jsonData, err := yaml.Marshal(cfg)
    if err != nil {
        return fmt.Errorf("failed to marshal MnemoConf to YAML: %w", err)
    }

    dir := filepath.Dir(targetPath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("failed to create directory structure for %s: %w", targetPath, err)
    }
    if err := os.WriteFile(targetPath, jsonData, 0644); err != nil {
        return fmt.Errorf("failed to write YAML data to file %s: %w", targetPath, err)
    }
    return nil
}

func GitDirExists(path string) (bool, error) {
	info, err := os.Stat(filepath.Join(path, ".git"))
	if err == nil {
		return info.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func healConfigSchema(loadedCfg *MnemoConf, defaultCfg *MnemoConf) []error {
	warnings := make([]error, 0)
	
	loadedSchema := &loadedCfg.ConfigSchema
	defaultSchema := defaultCfg.ConfigSchema
    
	replaceField := func(field *string, defaultVal string, fieldName string, reason string) {
		*field = defaultVal
		warnings = append(warnings, fmt.Errorf("invalid or empty field '%s': %s Overridden with default: '%s'", fieldName, reason, defaultVal))
	}

    if loadedSchema.AppVersion != defaultSchema.AppVersion {
        replaceField(&loadedSchema.AppVersion, defaultSchema.AppVersion, "AppVersion", "")
    }

    if !loadedSchema.IsInit {
        warnings = append(warnings, fmt.Errorf("found configuration file marked IsInit=false. Resetting RepoPath/DbPath."))

        loadedSchema.RepoPath = defaultSchema.RepoPath
        loadedSchema.DbPath = defaultSchema.DbPath
        return warnings 
    }
    
    if loadedSchema.RepoPath == "" {
        replaceField(&loadedSchema.RepoPath, defaultSchema.RepoPath, "RepoPath", "Cannot be empty when initialized")
    } else if _, err := os.Stat(loadedSchema.RepoPath); os.IsNotExist(err) {
        replaceField(&loadedSchema.RepoPath, defaultSchema.RepoPath, "RepoPath", fmt.Sprintf("Path does not exist on disk: %s", loadedSchema.RepoPath))
    }

    if loadedSchema.DbPath == "" {
        replaceField(&loadedSchema.DbPath, defaultSchema.DbPath, "DbPath", "Cannot be empty when initialized")
    }
    
    if loadedSchema.ConfigPath == "" {
        replaceField(&loadedSchema.ConfigPath, defaultSchema.ConfigPath, "ConfigPath", fmt.Sprintf("File path mismatch: %s", loadedSchema.ConfigPath))
    }

    return warnings
}
