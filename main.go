/*
 * Filename: main.go
 * Author: Nathaniel Thomas
 */

package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	version        string
	buildDate      string
	appName        = "go-cloud"
	appDescription = "WebDav Server"
)

func registerCLI() ([]*cli.Command, []cli.Flag) {
	return []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Output version",
				Action: func(c *cli.Context) error {
					fmt.Println("Version:   ", version)
					fmt.Println("Build Date:", buildDate)

					return nil
				},
			},
		},
		[]cli.Flag{}
}

func launchApp(c *cli.Context) error {
	return nil
}

func main() {
	// Manage CLI switches
	// Setup command routes
	commands, flags := registerCLI()

	// Configure application
	mainApp := &cli.App{
		Name:  appName,
		Usage: appDescription,
		Action: func(c *cli.Context) error {
			return launchApp(c)
		},
		Commands: commands,
		Flags:    flags,
	}

	// Run the app
	err := mainApp.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
