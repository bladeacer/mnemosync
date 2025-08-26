package cmd

import (
	"fmt"
	"mmsync/config"
	"github.com/spf13/cobra"
)

// getConfigPathCmd represents the get-config-path command
var getConfigPathCmd = &cobra.Command{
	Use:   "get-config-path",
	Short: "Prints the configuration file path",
	Long: `Prints the path to the configuration file that mmsync will use.
This respects the MMSYNC_CONF environment variable if it is set, 
otherwise it returns the default path.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath := config.ResolveConfigPath()
		isInit := appConf.ConfigSchema.IsInit

		fmt.Printf("\nConfiguration file path:\n%s\n", configPath)

		if !isInit {
			fmt.Printf("\nConfiguration file not found at expected path %s\nRun mmsync init to start.\n", configPath)
		}

	},
}

func init() {
	rootCmd.AddCommand(getConfigPathCmd)
}
