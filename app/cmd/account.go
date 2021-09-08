/*
 * Filename: account.go
 * Author: Nathaniel Thomas
 */

package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gitea.nthomas20.net/nathaniel/go-cloud/app/bootstrap"
	"gitea.nthomas20.net/nathaniel/go-cloud/app/configuration"
	"gitea.nthomas20.net/nathaniel/go-cloud/app/models"
	"github.com/urfave/cli/v2"
)

func addAccount(c *cli.Context) error {
	var (
		config    = configuration.NewConfiguration()
		username  = strings.ToLower(c.String("username"))
		email     = c.String("email")
		directory = c.String("directory")
		admin     = c.Bool("admin")
	)

	// Read Configuration
	bootstrap.SetupConfiguration()
	if err := configuration.ReadConfiguration(config); err != nil {
		return err
	}

	// Check for existing account
	if _, found := config.Accounts[username]; found {
		return errors.New("Account " + username + " already exists")
	}

	// Check directory existence
	_, err := os.Stat(directory)
	if os.IsNotExist(err) {
		return errors.New("directory does not exist")
	}

	// Add new account
	config.Accounts[username] = models.AccountConfiguration{
		Username:      username,
		Email:         email,
		Passwords:     make([]models.PasswordConfiguration, 0),
		IsActive:      true,
		IsAdmin:       admin,
		RootDirectory: directory,
	}

	// Write Configuration
	if err := configuration.WriteConfiguration(config); err != nil {
		return err
	}

	fmt.Println("Account " + username + " added")

	return nil
}

func deleteAccount(c *cli.Context) error {
	var (
		config   = configuration.NewConfiguration()
		username = strings.ToLower(c.String("username"))
	)

	// Read Configuration
	bootstrap.SetupConfiguration()
	if err := configuration.ReadConfiguration(config); err != nil {
		return err
	}

	// Check for existing account
	if _, found := config.Accounts[username]; !found {
		return errors.New("Account " + username + " does not exist")
	}

	delete(config.Accounts, username)

	// Write Configuration
	if err := configuration.WriteConfiguration(config); err != nil {
		return err
	}

	fmt.Println("Account " + username + " deleted")

	return nil
}

func updateAccount(c *cli.Context) error {
	var (
		config    = configuration.NewConfiguration()
		username  = strings.ToLower(c.String("username"))
		email     = c.String("email")
		directory = c.String("directory")
		admin     = c.Bool("admin")
		noadmin   = c.Bool("noadmin")
		active    = c.Bool("active")
		inactive  = c.Bool("inactive")
	)

	// Read Configuration
	bootstrap.SetupConfiguration()
	if err := configuration.ReadConfiguration(config); err != nil {
		return err
	}

	// Check for existing account
	if _, found := config.Accounts[username]; !found {
		return errors.New("Account " + username + " does not exist")
	}

	// Grab the account for manipulation
	account := config.Accounts[username]

	if email != "" {
		account.Email = email
	}

	if directory != "" {
		// Check directory existence
		_, err := os.Stat(directory)
		if os.IsNotExist(err) {
			return errors.New("directory does not exist")
		}

		account.RootDirectory = directory
	}

	if noadmin {
		account.IsAdmin = false
	} else if admin {
		account.IsAdmin = true
	}

	if inactive {
		account.IsActive = false
	} else if active {
		account.IsActive = true
	}

	// Assign the account after changes
	config.Accounts[username] = account

	// Write Configuration
	if err := configuration.WriteConfiguration(config); err != nil {
		return err
	}

	fmt.Println("Account " + username + " updated")

	return nil
}
