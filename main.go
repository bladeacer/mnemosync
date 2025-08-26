/*
Copyright © 2025 bladeacer <wg.nick.exe@gmail.com>
*/

package main

import (
	"mmsync/cmd"
	"mmsync/config"
	"fmt"
	"os"
)

func main() {
	euid := os.Geteuid()

	// Check if the effective user ID is 0 (which indicates root).
	if euid == 0 {
		fmt.Println("Warning: mnemosync should not be run as root.")
		os.Exit(1) // Exit with an error code
	}

	// Load the configuration file, falling back to defaults if it doesn't exist.
	appConfig, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading configuration: %v", err)
	}

	cmd.Execute(appConfig)
}
