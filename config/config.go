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
