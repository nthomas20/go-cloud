/*
 * Filename: list.go
 * Author: Nathaniel Thomas
 */

package cmd

import (
	"errors"
	"fmt"
	"strings"

	"gitea.nthomas20.net/nathaniel/go-cloud/app/bootstrap"
	"gitea.nthomas20.net/nathaniel/go-cloud/app/configuration"
	"github.com/urfave/cli/v2"
)

func listPassword(c *cli.Context) error {
	var (
		config   = configuration.NewConfiguration()
		username = strings.ToLower(c.String("username"))
		password = c.Int("password")
	)

	// Read Configuration
	bootstrap.SetupConfiguration()
	if err := configuration.ReadConfiguration(config); err != nil {
		return err
	}

	if password == 0 {
		// List all passwords
		for i, p := range config.Accounts[username].Passwords {
			fmt.Println(i, p.Password)
		}
	} else {
		// Check for existing account
		if _, found := config.Accounts[username]; !found {
			return errors.New("Account " + username + " does not exist")
		}

		// Check for password index out-of-bounds
		if password < 0 || password > len(config.Accounts[username].Passwords) {
			return errors.New("invalid password index")
		}

		// Grab the account passwords
		account := config.Accounts[username].Passwords[password]

		fmt.Println("Password:   ", account.Password)
		fmt.Println("Description:", account.Description)
	}

	return nil
}

func listUsername(c *cli.Context) error {
	var (
		config   = configuration.NewConfiguration()
		username = strings.ToLower(c.String("username"))
	)

	// Read Configuration
	bootstrap.SetupConfiguration()
	if err := configuration.ReadConfiguration(config); err != nil {
		return err
	}

	if username == "" {
		// List all usernames
		for u := range config.Accounts {
			fmt.Println(u)
		}
	} else {
		// Check for existing account
		if _, found := config.Accounts[username]; !found {
			return errors.New("Account " + username + " does not exist")
		}

		// Grab the account information
		account := config.Accounts[username]

		fmt.Println("Username: ", account.Username)
		fmt.Println("Active:   ", account.IsActive)
		fmt.Println("Admin:    ", account.IsAdmin)
		fmt.Println("Email:    ", account.Email)
		fmt.Println("Directory:", account.RootDirectory)
		fmt.Println("Passwords:", len(account.Passwords))
	}

	return nil
}
