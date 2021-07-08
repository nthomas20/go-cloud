/*
 * Filename: bootstrap.go
 * Author: Nathaniel Thomas
 */

package bootstrap

import (
	"log"
	"os"
)

var (
	// HomeDirectory : User's Home Directory
	HomeDirectory string
	// ConfigDirectory : Main Configuration Directory
	ConfigDirectory string
)

// SetupConfiguration : Make sure our main configuration directory is set
func SetupConfiguration() {
	// Get the home directory
	homeDirectory, _ := os.UserHomeDir()

	ConfigDirectory = homeDirectory + "/.go-cloud"

	_, err := os.Stat(ConfigDirectory)
	if os.IsNotExist(err) {
		// Create directory
		if err := os.Mkdir(ConfigDirectory, os.ModePerm); err != nil {
			log.Fatal("Could not create configuration directory:", ConfigDirectory)
		}
	}
}
