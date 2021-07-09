/*
 * Filename: configuration.go
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

func addPassword(c *cli.Context) error {
	var (
		config      = configuration.NewConfiguration()
		username    = c.String("username")
		password    = c.String("password")
		description = c.String("description")
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

	// Check for existing password
	if _, found := config.Accounts[username].Passwords[password]; found == true {
		return errors.New("Account password already exists")
	}

	// Add new password
	config.Accounts[username].Passwords[password] = models.PasswordConfiguration{
		Password:    password,
		Description: description,
	}

	// Write Configuration
	if err := configuration.WriteConfiguration(config); err != nil {
		return err
	}

	fmt.Println("Account password added")

	return nil
}

func deletePassword(c *cli.Context) error {
	var (
		config   = configuration.NewConfiguration()
		username = c.String("username")
		password = c.String("password")
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

	// Check for existing password
	if _, found := config.Accounts[username].Passwords[password]; found == false {
		return errors.New("Account password does not exist")
	}

	delete(config.Accounts[username].Passwords, password)

	// Write Configuration
	if err := configuration.WriteConfiguration(config); err != nil {
		return err
	}

	fmt.Println("Account password deleted")

	if len(config.Accounts[username].Passwords) == 0 {
		fmt.Println("No passwords set for account " + username)
	}

	return nil
}
