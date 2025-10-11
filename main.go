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

	appConfig, err := config.LoadConfig()
	dataStore, err2 := config.LoadDataStore()
	if err != nil {
		fmt.Printf("Error loading configuration: %v", err)
	}
	if err2 != nil {
		fmt.Printf("Error loading database: %v", err)
	}

	cmd.Execute(appConfig, dataStore)
}
