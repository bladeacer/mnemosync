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
			fmt.Printf("\tConfiguration file not found at: %s\n", configPath)
			fmt.Printf("\n\t%s\n", repeatedSeparator)
			fmt.Printf("\n\tConfiguration file not found. Run 'mmsync init' to start.")
			fmt.Printf("\n\t%s\n", repeatedSeparator)
		} else {
			// File exists, so print that it's found.
			fmt.Printf("\tConfiguration file exists:\n\t%s\n", configPath)
			fmt.Printf("\n\t%s\n", repeatedSeparator)
		}

		fmt.Println("\n\tHealth Check Complete")
	},
}

// Renamed helper functions to follow Go conventions
func checkBinWrapper(binaryName string, isOptional bool) {
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
}

// Here you will define your flags and configuration settings.

// Cobra supports Persistent Flags which will work for this command
// and all subcommands, e.g.:
// healthCmd.PersistentFlags().String("foo", "", "A help for foo")

// Cobra supports local flags which will only run when this command
// is called directly, e.g.:
// healthCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
