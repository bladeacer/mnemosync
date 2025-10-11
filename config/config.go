/*
Copyright © 2025 bladeacer <wg.nick.exe@gmail.com>
*/

package config

import (
	"os"
	"path/filepath"
	"gopkg.in/yaml.v3"
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
	
	cfg := GetMnemoConf()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}
	
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
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
