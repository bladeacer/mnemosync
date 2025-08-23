/*
Copyright © 2025 bladeacer <wg.nick.exe@gmail.com>

*/
package cmd

import (
	"fmt"
	"os/exec"
	"github.com/spf13/cobra"
	"strings"
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

		fmt.Println("\tRunning Health Check")
		fmt.Printf("\t%s\n\n", repeatedSeparator)
		check_bin_wrapper("git", false)
		check_bin_wrapper("rsync", false)
		check_bin_wrapper("tar", false)
		check_bin_wrapper("zip", true)
		fmt.Printf("\t%s", repeatedSeparator)
		fmt.Println("\n\n\tHealth Check Complete")
	},
}

func check_bin_wrapper(binaryName string, isOptional bool) {
	path, err := exec.LookPath(binaryName)
	if err != nil && !isOptional {
		fmt.Printf("\tWarning: Binary '%s' not found or not executable in PATH:\n\t\t%v\n", binaryName, err)
		return
	}

	if err != nil && isOptional {
		fmt.Printf("\tWarning: Optional Binary '%s' not found or not executable in PATH:\n\t\t%v\n", binaryName, err)
		return
	}

	fmt.Printf("\tBinary '%s' found at: %s\n", binaryName, path)	
}

func init() {
	rootCmd.AddCommand(healthCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// healthCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// healthCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
