/*
Copyright © 2025 bladeacer wg.nick.exe@gmail.com

*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"github.com/bladeacer/mmsync/config"
	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage the Git repository path used for archiving",
	Long:  "Provides commands to view and open the configured Git repository path.",
}

var repoGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Prints the configured repository path to stdout",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if the config is initialized
		if appConf == nil || !appConf.ConfigSchema.IsInit {
			configPath := config.ResolveConfigPath()
			fmt.Printf("Error: Configuration file not found or not initialized at expected path:\n%s\nRun mmsync init to start.\n", configPath)
			os.Exit(1)
		}

		repoPath := appConf.ConfigSchema.RepoPath

		if repoPath == "" {
			fmt.Println("Error: Repository path is not set in the configuration file.")
			os.Exit(1)
		}

		fmt.Println(repoPath)
	},
}

var repoOpenCmd = &cobra.Command{
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
		
		if !isInit {
			fmt.Printf("\nConfiguration file not found at expected path\n%s\nRun mmsync init to start.\n", configPath)
			os.Exit(1)
		}

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
	rootCmd.AddCommand(repoCmd)
	
	repoCmd.AddCommand(repoGetCmd)
	repoCmd.AddCommand(repoOpenCmd)
}
