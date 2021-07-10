/*
 * Filename: account.go
 * Author: Nathaniel Thomas
 */

package cmd

import (
	"errors"
	"fmt"

	"gitea.nthomas20.net/nathaniel/go-cloud/app/bootstrap"
	"gitea.nthomas20.net/nathaniel/go-cloud/app/configuration"
	"gitea.nthomas20.net/nathaniel/go-cloud/app/models"
	"github.com/urfave/cli/v2"
)

func addAccount(c *cli.Context) error {
	var (
		config    = configuration.NewConfiguration()
		username  = c.String("username")
		email     = c.String("email")
		directory = c.String("directory")
	)

	// Read Configuration
	bootstrap.SetupConfiguration()
	if err := configuration.ReadConfiguration(config); err != nil {
		return err
	}

	// Check for existing account
	if _, found := config.Accounts[username]; found == true {
		return errors.New("Account " + username + " already exists")
	}

	// Add new account
	config.Accounts[username] = models.AccountConfiguration{
		Username:      username,
		Email:         email,
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
		username = c.String("username")
	)

	// Read Configuration
	bootstrap.SetupConfiguration()
	if err := configuration.ReadConfiguration(config); err != nil {
		return err
	}

	// Check for existing account
	if _, found := config.Accounts[username]; found == false {
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
