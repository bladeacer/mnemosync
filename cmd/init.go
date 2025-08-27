/*
Copyright © 2025 bladeacer <wg.nick.exe@gmail.com>
*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"mmsync/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// initCmd is the command for creating a default config file.
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new configuration file with default values.",
	Run: func(cmd *cobra.Command, args []string) {
		// Use the ResolveConfigPath helper to get the path.
		configPath := config.ResolveConfigPath()

		// Get the default configuration.
		defaultConfig := config.GetMnemoConf()
		defaultConfig.ConfigSchema.IsInit = true

		// Marshal the default config into bytes.
		data, err := yaml.Marshal(defaultConfig)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating default config:", err)
			return
		}

		// Check if the directory exists, and create it if not.
		dir := filepath.Dir(configPath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintln(os.Stderr, "Error creating config directory:", err)
				return
			}
		}

		// Write the default config to the determined path.
		if err := os.WriteFile(configPath, data, 0644); err != nil {
			fmt.Fprintln(os.Stderr, "Error writing config file:", err)
			return
		}
		
		fmt.Printf("Initialized default configuration file at %s\n", configPath)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
