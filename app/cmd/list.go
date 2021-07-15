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
		password = c.String("password")
	)

	// Read Configuration
	bootstrap.SetupConfiguration()
	if err := configuration.ReadConfiguration(config); err != nil {
		return err
	}

	if password == "" {
		// List all passwords
		for p := range config.Accounts[username].Passwords {
			fmt.Println(p)
		}
	} else {
		// Check for existing account
		if _, found := config.Accounts[username]; found == false {
			return errors.New("Account " + username + " does not exist")
		}

		// Check for existing password information
		if _, found := config.Accounts[username].Passwords[password]; found == false {
			return errors.New("Account password does not exist")
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
		if _, found := config.Accounts[username]; found == false {
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
