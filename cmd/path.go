/*
Copyright © 2025 bladeacer <wg.nick.exe@gmail.com>
*/

package cmd

import (
	"fmt"
	"mmsync/config"
	"os"
	"os/exec"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage the mnemosync configuration file",
	Long:  "Provides commands to manage the application's configuration file.",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := config.ResolveConfigPath()
		isInit := appConf.ConfigSchema.IsInit

		// Default behavior (original get-config-path functionality)
		fmt.Printf("\nConfiguration file path:\n%s\n", configPath)
		if !isInit {
			fmt.Printf("\nConfiguration file not found at expected path\n%s\nRun mmsync init to start.\n", configPath)
		}
	},
}

// openCmd represents the open subcommand
var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Opens the configuration file with the user's $EDITOR",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := config.ResolveConfigPath()
		isInit := appConf.ConfigSchema.IsInit
		editor := os.Getenv("EDITOR")

		if editor == "" {
			fmt.Println("Error: $EDITOR environment variable not set. Please set it to your preferred text editor (e.g., 'vim', 'code').")
			os.Exit(1)
		}
		
		// Check if the config file exists before trying to open it
		if !isInit {
			fmt.Printf("\nConfiguration file not found at expected path\n%s\nRun mmsync init to start.\n", configPath)
			os.Exit(1)
		}

		// Open the file with the user's editor
		editorCmd := exec.Command(editor, configPath)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr

		err := editorCmd.Run()
		if err != nil {
			fmt.Printf("Error: failed to open config file with %s: %v\n", editor, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	// Add the open subcommand to the config command
	configCmd.AddCommand(openCmd)
}
