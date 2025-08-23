/*
Copyright © 2025 bladeacer <wg.nick.exe@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Gets the version of mnemosync running",
	Long: `Gets the version of mneomsync running. Help text should be more detailed. :D
Code is currently very WIP.
Improve this text in the long run.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mnemosync Version: 0.0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
