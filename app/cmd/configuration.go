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
			Usage:   "Specify directory path (e.g. /var/data/webdav) [optional]",
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
			Usage:   "Specify email [optional]",
		},
		"key": &cli.StringFlag{
			Name:     "key",
			Aliases:  []string{"k"},
			Usage:    "Specify key",
			Required: true,
		},
		"password": &cli.StringFlag{
			Name:     "password",
			Aliases:  []string{"p"},
			Usage:    "Specify password",
			Required: true,
		},
		"password-optional": &cli.StringFlag{
			Name:    "password",
			Aliases: []string{"p"},
			Usage:   "Specify password [optional]",
		},
		"username": &cli.StringFlag{
			Name:     "username",
			Aliases:  []string{"u"},
			Usage:    "Specify username",
			Required: true,
		},
		"username-optional": &cli.StringFlag{
			Name:    "username",
			Aliases: []string{"u"},
			Usage:   "Specify username [optional]",
		},
		"value": &cli.StringFlag{
			Name:     "value",
			Aliases:  []string{"v"},
			Usage:    "Specify value",
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
						flags["password-optional"],
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
		},
		// Configuration Management
		{
			Name:    "config",
			Usage:   "Configuration management command",
			Aliases: []string{"c"},
			Subcommands: []*cli.Command{
				{
					Name:    "get",
					Usage:   "Get configuration value",
					Aliases: []string{"g"},
					Action:  getConfig,
					Flags: []cli.Flag{
						flags["key"],
					},
				},
				{
					Name:    "set",
					Usage:   "Set configuration value",
					Aliases: []string{"s"},
					Action:  setConfig,
					Flags: []cli.Flag{
						flags["key"],
						flags["value"],
					},
				},
			},
		},
		// Lists!
		{
			Name:    "list",
			Usage:   "List out information",
			Aliases: []string{"l"},
			Subcommands: []*cli.Command{
				{
					Name:    "account",
					Usage:   "List accounts",
					Aliases: []string{"a"},
					Action:  listUsername,
					Flags: []cli.Flag{
						flags["username-optional"],
					},
				},
				{
					Name:    "password",
					Usage:   "List passwords",
					Aliases: []string{"p"},
					Action:  listPassword,
					Flags: []cli.Flag{
						flags["username"],
						flags["password-optional"],
					},
				},
			},
		},
	}
}
