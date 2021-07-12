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
		"active": &cli.BoolFlag{
			Name:  "active",
			Usage: "Specify active state",
			Value: false,
		},
		"inactive": &cli.BoolFlag{
			Name:  "inactive",
			Usage: "Specify active state",
			Value: false,
		},
		"admin": &cli.BoolFlag{
			Name:  "admin",
			Usage: "Specify administrator privileges",
			Value: false,
		},
		"noadmin": &cli.BoolFlag{
			Name:  "noadmin",
			Usage: "Specify no administrator privileges",
			Value: false,
		},
		"description": &cli.StringFlag{
			Name:    "description",
			Aliases: []string{"desc"},
			Value:   "",
			Usage:   "Specify description [optional]",
		},
		"directory": &cli.StringFlag{
			Name:     "directory",
			Aliases:  []string{"dir"},
			Usage:    "Specify directory path (e.g. /var/data/webdav)",
			Required: true,
		},
		"directory-optional": &cli.StringFlag{
			Name:    "directory",
			Aliases: []string{"dir"},
			Usage:   "Specify directory path (e.g. /var/data/webdav)",
		},
		"email": &cli.StringFlag{
			Name:     "email",
			Aliases:  []string{"e"},
			Usage:    "Specify email",
			Required: true,
		},
		"email-optional": &cli.StringFlag{
			Name:    "email",
			Aliases: []string{"e"},
			Usage:   "Specify email",
		},
		"password": &cli.StringFlag{
			Name:     "password",
			Aliases:  []string{"p"},
			Usage:    "Specify password",
			Required: true,
		},
		"username": &cli.StringFlag{
			Name:     "username",
			Aliases:  []string{"u"},
			Usage:    "Specify username",
			Required: true,
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
					Action:  addAccount,
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
					Action:  deleteAccount,
					Flags: []cli.Flag{
						flags["username"],
					},
				},
				{
					Name:    "update",
					Usage:   "Update account",
					Aliases: []string{"u"},
					Action:  updateAccount,
					Flags: []cli.Flag{
						flags["username"],
						flags["email-optional"],
						flags["directory-optional"],
						flags["admin"],
						flags["noadmin"],
						flags["active"],
						flags["inactive"],
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
					Action:  addPassword,
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
					Action:  deletePassword,
					Flags: []cli.Flag{
						flags["username"],
						flags["password"],
					},
				},
			},
		}}
}
