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
	"bufio"
	"strings"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new configuration file with default values.",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := config.ResolveConfigPath()

		if _, err := os.Stat(configPath); err == nil {
			fmt.Fprintf(os.Stderr, "Error: Configuration file already exists at %s\n", configPath)
			os.Exit(1)
		} else if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error checking for config file at %s: %v\n", configPath, err)
			os.Exit(1)
		}

		defaultConfig := config.GetMnemoConf()
		repoPath := get_repo_path()
		defaultConfig.ConfigSchema.IsInit = true
		defaultConfig.ConfigSchema.RepoPath = repoPath

		exists, err := config.DirExists(repoPath)
		if exists {
			fmt.Printf("\nDirectory '%s/.git' exists.\n", repoPath)
			write_yaml(defaultConfig, configPath)
		} else if err != nil {
			fmt.Printf("\nDirectory '%s/.git' does not exist.\n", repoPath)
			fmt.Printf("Aborting write\n")
		}

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func get_repo_path() string {
	reader := bufio.NewReader(os.Stdin)
	inputPathBuf := ""
	fmt.Println("Ensure that the target repository path is correct and does not contain other important files.")
	for {
		fmt.Println("Enter a valid absolute path to the target repository to archive files to: ")
		inputPath, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		inputPath = strings.TrimSpace(inputPath)

		info, err := os.Stat(inputPath)
		if os.IsNotExist(err) {
			fmt.Printf("Error: Directory '%s' does not exist.\n", inputPath)
			continue
		} else if err != nil {
			fmt.Println("Error checking path:", err)
			continue
		}

		if !info.IsDir() {
			fmt.Printf("Error: '%s' is not a directory.\n", inputPath)
			continue
		}

		absPath, err := filepath.Abs(inputPath)
		if err != nil {
			fmt.Println("Error getting absolute path:", err)
			continue
		}

		fmt.Printf("You entered the directory: %s\n", absPath)
		inputPathBuf = inputPath
		break
	}

	exists, err := config.DirExists(inputPathBuf)
	if exists {
		return inputPathBuf
	} else if err != nil {
		fmt.Printf("Directory '%s/.git' does not exist.\n", inputPathBuf)
		fmt.Printf("Aborting write\n")
		return ""
	} else {
		return ""
	}
}

func write_yaml (defaultConfig *config.MnemoConf, configPath string) {
	data, err := yaml.Marshal(defaultConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating default config:", err)
		return
	}

	dir := filepath.Dir(configPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintln(os.Stderr, "Error creating config directory:", err)
			return
		}
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		fmt.Fprintln(os.Stderr, "Error writing config file:", err)
		return
	}

	fmt.Printf("Initialized default configuration file at %s\n", configPath)
}
