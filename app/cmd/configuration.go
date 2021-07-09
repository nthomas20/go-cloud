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
				},
				{
					Name:    "delete",
					Usage:   "Delete an account",
					Aliases: []string{"d"},
					Action: func(c *cli.Context) error {
						return nil
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
				},
				{
					Name:    "delete",
					Usage:   "Delete an account password",
					Aliases: []string{"d"},
					Action: func(c *cli.Context) error {
						return nil
					},
				},
			},
		}}
}
