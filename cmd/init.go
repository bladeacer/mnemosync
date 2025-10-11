package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"mmsync/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Global variable to hold the path passed via flag
var repoPathFlag string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new configuration file with default values.",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := config.ResolveConfigPath()
		dbPath := config.ResolveDbPath()
		_, confErr := os.Stat(configPath)
		_, dbErr := os.Stat(dbPath)

		if confErr == nil || dbErr == nil {
			fmt.Fprintf(os.Stderr, "Error: Cannot run init. The following files already exist:\n")
			if confErr == nil {
				fmt.Fprintf(os.Stderr, "- Configuration file at %s\n", configPath)
			}
			if dbErr == nil {
				fmt.Fprintf(os.Stderr, "- Database file at %s\n", dbPath)
			}
			fmt.Fprintf(os.Stderr, "Please remove the existing files before running 'init'.\n")
			os.Exit(1)
		} else {
			if !os.IsNotExist(confErr) {
				fmt.Fprintf(os.Stderr, "Error checking for config file at %s: %v\n", configPath, confErr)
			}
			if !os.IsNotExist(dbErr) {
				fmt.Fprintf(os.Stderr, "Error checking for database file at %s: %v\n", dbPath, dbErr)
			}
		}

		var finalRepoPath string
		var err error

		if repoPathFlag != "" {
			finalRepoPath, err = processRepoPath(repoPathFlag)
		} else {
			finalRepoPath, err = getRepoPathInteractive()
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "\nInitialization aborted: %v\n", err)
			os.Exit(1)
		}

		defaultConfig := config.GetMnemoConf()
		defaultConfig.ConfigSchema.IsInit = true
		defaultConfig.ConfigSchema.RepoPath = finalRepoPath

		exists, _ := config.GitDirExists(finalRepoPath)
		
		if exists {
			fmt.Printf("\nRepository path validated: '%s/.git' exists.\n", finalRepoPath)
			writeYAML(defaultConfig, configPath)
			config.GetDataStore().SaveData(dbPath) 
			fmt.Printf("\nDatabase created at: '%s'.\n", dbPath)
		} else {
			fmt.Printf("\nDirectory '%s/.git' does not exist.\n", finalRepoPath)
			fmt.Printf("Aborting configuration write.\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&repoPathFlag, "repo-path", "r", "", "Specify the path to the target Git repository.")
}

func processRepoPath(inputPath string) (string, error) {
	if strings.HasPrefix(inputPath, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		inputPath = filepath.Join(homeDir, inputPath[2:])
	}

	absPath, err := filepath.Abs(inputPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path for '%s': %w", inputPath, err)
	}

	info, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("directory '%s' does not exist", absPath)
	} else if err != nil {
		return "", fmt.Errorf("error checking path '%s': %w", absPath, err)
	}

	if !info.IsDir() {
		return "", fmt.Errorf("'%s' is not a directory", absPath)
	}

	return absPath, nil
}


func getRepoPathInteractive() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Ensure that the target repository path is correct and does not contain other important files.\nDatabase for storing directories and their aliases would use the same parent directory.")
	
	for {
		fmt.Println("Enter a valid path to the target repository to archive files to (e.g., /path/to/repo or ~/myrepo): ")
		
		inputPath, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			continue
		}
		
		inputPath = strings.TrimSpace(inputPath)
		if inputPath == "" {
			continue
		}

		finalRepoPath, err := processRepoPath(inputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}
		
		fmt.Printf("Path accepted: %s\n", finalRepoPath)
		return finalRepoPath, nil
	}
}

func writeYAML(defaultConfig *config.MnemoConf, configPath string) {
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
