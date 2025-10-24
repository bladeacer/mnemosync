package config

import (
	"os"
	"io"
	"path/filepath"
	"gopkg.in/yaml.v3"
	"fmt"
)
/*
Copyright © 2025 bladeacer <wg.nick.exe@gmail.com>
*/

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

const (
    DefaultConfigDir = ".config/mmsync"
    DefaultConfigFile = "config.yaml"
    DefaultDbFile = "mmsync-state.json"
)

func ResolveConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(DefaultConfigDir, DefaultConfigFile)
	}

	if envPath := os.Getenv("MMSYNC_CONF"); envPath != "" {
		resolvedPath := envPath
		
		if !filepath.IsAbs(envPath) {
			resolvedPath = filepath.Join(homeDir, envPath)
		}

		if filepath.Base(resolvedPath) != DefaultConfigFile {
			resolvedPath = filepath.Join(resolvedPath, DefaultConfigFile)
		}
		
		return resolvedPath
	}

	return filepath.Join(homeDir, DefaultConfigDir, DefaultConfigFile)
}

func ResolveDbPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(DefaultConfigDir, DefaultDbFile)
	}

	if envPath := os.Getenv("MMSYNC_CONF"); envPath != "" {
		configPath := ResolveConfigPath()
		
		configDir := filepath.Dir(configPath)
		return filepath.Join(configDir, DefaultDbFile)
	}

	return filepath.Join(homeDir, DefaultConfigDir, DefaultDbFile)
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer sourceFile.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory for %s: %w", dst, err)
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy content from %s to %s: %w", src, dst, err)
	}

	if err := destFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}

// Copies configuration and database files when new MMSYNC_CONF set
func migrateConfigData(newConfigPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	oldConfigDir := filepath.Join(homeDir, DefaultConfigDir)
	oldConfigFile := filepath.Join(oldConfigDir, DefaultConfigFile)
	oldDbFile := filepath.Join(oldConfigDir, DefaultDbFile)
	
	newConfigDir := filepath.Dir(newConfigPath)
	newDbFile := filepath.Join(newConfigDir, DefaultDbFile)

	if os.Getenv("MMSYNC_CONF") != "" && oldConfigDir != newConfigDir {

		if _, err := os.Stat(oldConfigFile); os.IsNotExist(err) {
			return nil
		} else if err != nil {
			return fmt.Errorf("error checking old configuration file %s: %w", oldConfigFile, err)
		}

		if _, err := os.Stat(newConfigPath); err == nil {
			return fmt.Errorf("cannot migrate: new configuration file already exists at %s", newConfigPath)
		}

		fmt.Fprintf(os.Stderr, "Migrating configuration files from %s to %s...\n", oldConfigDir, newConfigDir)
		
		if err := copyFile(oldConfigFile, newConfigPath); err != nil {
			return fmt.Errorf("failed to copy configuration file: %w", err)
		}
		
		if _, err := os.Stat(oldDbFile); err == nil {
			if err := copyFile(oldDbFile, newDbFile); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to copy database file: %v\n", err)
			}
		}

		if err := os.RemoveAll(oldConfigDir); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to clean up old configuration directory %s: %v\n", oldConfigDir, err)
		}
		
		fmt.Fprintf(os.Stderr, "Configuration migration complete.\n")
	}

	return nil
}

func LoadConfig() (*MnemoConf, error) {
	configPath := ResolveConfigPath()
	
	if err := migrateConfigData(configPath); err != nil {
        	return nil, fmt.Errorf("Configuration migration failed: %w", err)
    	}
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
