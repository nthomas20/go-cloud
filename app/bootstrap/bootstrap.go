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
	// ConfigFilename : Configuration Filename
	ConfigFilename = "/config.yaml"
)

// SetupConfiguration : Make sure our main configuration directory is set
func SetupConfiguration() {
	// Check if there is a config file in the current directory, if so -- it wins!
	if dir, err := os.Getwd(); err == nil {
		if _, err := os.Stat(dir + ConfigFilename); err == nil {
			ConfigDirectory = dir

			return
		}
	}

	// Get the home directory
	if homeDirectory, err := os.UserHomeDir(); err == nil {
		ConfigDirectory = homeDirectory + "/.go-cloud"

		_, err := os.Stat(ConfigDirectory)
		if os.IsNotExist(err) {
			// Create directory
			if err := os.Mkdir(ConfigDirectory, os.ModePerm); err != nil {
				log.Fatalln("Could not create configuration directory:", ConfigDirectory)
			}
		}
	} else {
		log.Fatalln("Could not determine user's home directory")

	}
}
