/*
 * Filename: configuration.go
 * Author: Nathaniel Thomas
 */

package cmd

import (
	"github.com/urfave/cli/v2"
)

// Commands : Return the full set of registered commands
func Commands() []*cli.Command {
	flags := map[string]cli.Flag{
		"admin": &cli.BoolFlag{
			Name:  "admin",
			Usage: "Specify administrator privileges [optional]",
			Value: false,
		},
		"description": &cli.StringFlag{
			Name:  "description",
			Value: "",
			Usage: "Specify description [optional]",
		},
		"directory": &cli.StringFlag{
			Name:  "directory",
			Usage: "Specify directory path (e.g. /var/data/webdav)",
		},
		"email": &cli.StringFlag{
			Name:  "email",
			Usage: "Specify email",
		},
		"password": &cli.StringFlag{
			Name:  "password",
			Usage: "Specify password",
		},
		"username": &cli.StringFlag{
			Name:  "username",
			Usage: "Specify username",
		},
	}

	return []*cli.Command{
		// Account Management
		{
			Name:    "account",
			Usage:   "Account management command",
			Aliases: []string{"a"},
			Subcommands: []*cli.Command{
				{
					Name:    "add",
					Usage:   "Add an account",
					Aliases: []string{"a"},
					Action: func(c *cli.Context) error {
						return nil
					},
					Flags: []cli.Flag{
						flags["username"],
						flags["email"],
						flags["directory"],
						flags["admin"],
					},
				},
				{
					Name:    "delete",
					Usage:   "Delete an account",
					Aliases: []string{"d"},
					Action: func(c *cli.Context) error {
						return nil
					},
					Flags: []cli.Flag{
						flags["username"],
					},
				},
			},
		},
		// Password Management
		{
			Name:    "password",
			Usage:   "Account password management command",
			Aliases: []string{"p"},
			Subcommands: []*cli.Command{
				{
					Name:    "add",
					Usage:   "Add an account password",
					Aliases: []string{"a"},
					Action: func(c *cli.Context) error {
						return nil
					},
					Flags: []cli.Flag{
						flags["username"],
						flags["password"],
						flags["description"],
					},
				},
				{
					Name:    "delete",
					Usage:   "Delete an account password",
					Aliases: []string{"d"},
					Action: func(c *cli.Context) error {
						return nil
					},
					Flags: []cli.Flag{
						flags["username"],
						flags["password"],
					},
				},
			},
		}}
}
