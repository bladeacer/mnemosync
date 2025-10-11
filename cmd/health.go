/*
Copyright © 2025 bladeacer <wg.nick.exe@gmail.com>
*/

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"github.com/spf13/cobra"
)

// healthCmd represents the health command
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Checks the health of mnemosync",
	Long: `Checks the health of mnemosync
Checks if the required system binaries are installed

Also checks if the mnemosync configuration files have been created.`,
	Run: func(cmd *cobra.Command, args []string) {
		separator := "_"
		repeatedSeparator := strings.Repeat(separator, 72)
		
		// The configPath from appConf will be the resolved path.
		configPath := appConf.ConfigSchema.ConfigPath
		repoPath := appConf.ConfigSchema.RepoPath
		dbPath := appConf.ConfigSchema.DbPath
		
		fmt.Println("\n\tRunning Health Check")
		fmt.Printf("\t%s\n\n", repeatedSeparator)

		checkBinWrapper("git", false)
		checkBinWrapper("rsync", false)
		checkBinWrapper("tar", false)
		checkBinWrapper("zip", true)

		fmt.Printf("\t%s\n\n", repeatedSeparator)

		// Check for the existence of the specific config file.
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			// File does not exist, so print a warning and exit cleanly.
			fmt.Printf("\tConfiguration file not found at:\n\t%s\n\tRun 'mmsync init' to start.", configPath)
			fmt.Printf("\n\t%s\n", repeatedSeparator)
		} else {
			// File exists, so print that it's found.
			fmt.Printf("\tConfiguration file exists:\n\t%s\n", configPath)
			fmt.Printf("\n\t%s\n", repeatedSeparator)
		}

		if repoPath == "" {
			fmt.Printf("\n\tRepository Path is not defined.\n\tRun 'mmsync init' to start.")
			fmt.Printf("\n\t%s\n", repeatedSeparator)
		} else {
			fmt.Printf("\tRepository exists:\n\t%s\n", repoPath)
			fmt.Printf("\n\t%s\n", repeatedSeparator)
		}
		if dbPath == "" {
			fmt.Printf("\n\tDatabase Path is not defined.\n\tRun 'mmsync init' to start.")
			fmt.Printf("\n\t%s\n", repeatedSeparator)
		} else {
			fmt.Printf("\n\tDatabase exists:\n\t%s\n", dbPath)
			fmt.Printf("\n\t%s\n", repeatedSeparator)
		}


		fmt.Println("\n\tHealth Check Complete")
	},
}

func checkBinWrapper(binaryName string, isOptional bool) {
	path, err := exec.LookPath(binaryName)

	if err != nil {
		if !isOptional {
			fmt.Printf("\t[FAIL] Required Binary '%s' not found in PATH.\n", binaryName)
		} else {
			fmt.Printf("\t[WARN] Optional Binary '%s' not found in PATH.\n", binaryName)
		}
		return
	}

	fmt.Printf("\t[PASS] Binary '%s' found at: %s\n", binaryName, path)

	cmd := exec.Command(binaryName, "--version")
	output, versionErr := cmd.CombinedOutput()

	if versionErr != nil {
		// Log a specific message if the command failed, even if the binary was found
		if exitError, ok := versionErr.(*exec.ExitError); ok {
			fmt.Printf("\t\t[WARN] Version check failed (Exit Code %d). Output:\n\t\t%s", exitError.ExitCode(), strings.TrimSpace(string(output)))
		} else {
			fmt.Printf("\t\t[WARN] Version check failed to execute: %v\n", versionErr)
		}
		return
	}

	// Print the version output, trimmed for cleaner output
	versionLine := strings.SplitN(string(output), "\n", 2)[0]
	fmt.Printf("\t\tVersion: %s\n", strings.TrimSpace(versionLine))
}


func init() {
	rootCmd.AddCommand(healthCmd)
}

// Here you will define your flags and configuration settings.

// Cobra supports Persistent Flags which will work for this command
// and all subcommands, e.g.:
// healthCmd.PersistentFlags().String("foo", "", "A help for foo")

// Cobra supports local flags which will only run when this command
// is called directly, e.g.:
// healthCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
