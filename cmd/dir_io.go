/*
Copyright © 2025 bladeacer wg.nick.exe@gmail.com

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// TODO: This command helps add directory paths to be staged before performing backup. Have CRUD in this.
// Somehow rsync directories to the target directory and then tar archive all of them when push is called
// Save stuff in viewable local db instead of config file. They should be separate. Probably store it beside where the config file is located at. Default: ~/.config/mmsync/mmsync-state.db
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a target path to be staged before backing up to target repository",
	Long: ` Add a target path to be staged before backing up to target repository.
For example:

mmsync add ./

Adds the current directory recursively to be staged`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
