/*
Copyright © 2025 bladeacer <wg.nick.exe@gmail.com>
*/

package main

import (
	"mmsync/cmd"
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

	cmd.Execute()
}
