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

// Add a variable to hold the open flag value
var openConfigFlag bool

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage the mnemosync configuration file",
	Long:  "Provides commands and flags to manage the application's configuration file.",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := config.ResolveConfigPath()
		isInit := appConf.ConfigSchema.IsInit

		// Handle the --open flag
		if openConfigFlag {
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
			return
		}

		// Default behavior (original get-config-path functionality)
		fmt.Printf("\nConfiguration file path:\n%s\n", configPath)
		if !isInit {
			fmt.Printf("\nConfiguration file not found at expected path\n%s\nRun mmsync init to start.\n", configPath)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Add the --open flag to the config command
	configCmd.Flags().BoolVarP(&openConfigFlag, "open", "o", false, "Opens the configuration file with the user's $EDITOR")

	// Add the get-config-path as a subcommand (optional, as its logic is now the default behavior of `config`)
	// If you want to keep it as a subcommand:
	// configCmd.AddCommand(getConfigPathCmd)
}
