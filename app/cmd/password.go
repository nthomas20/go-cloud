/*
 * Filename: password.go
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
	if _, found := config.Accounts[username]; !found {
		return errors.New("Account " + username + " does not exist")
	}

	// Check for existing password
	found := false
	for _, p := range config.Accounts[username].Passwords {
		if p.Password == password {
			found = true
			break
		}
	}
	if found {
		return errors.New("account password already exists")
	}

	// Add new password
	account := config.Accounts[username]
	account.Passwords = append(account.Passwords, models.PasswordConfiguration{
		Password:    password,
		Description: description,
	})
	config.Accounts[username] = account

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
		password = c.Int("password")
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

	// Check for password index out-of-bounds
	if password < 0 || password > len(config.Accounts[username].Passwords) {
		return errors.New("invalid password index")
	}

	// Delete indexed item
	account := config.Accounts[username]
	account.Passwords = append(account.Passwords[:password], account.Passwords[password+1:]...)
	config.Accounts[username] = account

	// Write Configuration
	if err := configuration.WriteConfiguration(config); err != nil {
		return err
	}

	fmt.Println("Account password deleted")

	if len(config.Accounts[username].Passwords) == 0 {
		fmt.Println("no passwords set for account " + username)
	}

	return nil
}
