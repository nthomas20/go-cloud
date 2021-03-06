/*
 * Filename: password.go
 * Author: Nathaniel Thomas
 */

package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

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

	// If not password is provided, ask for it
	if password == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Password: ")
		password, _ = reader.ReadString('\n')
		password = strings.Trim(password, " \n")
	}

	if password == "" {
		return errors.New("invalid password")
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

	fmt.Println("account password added")

	return nil
}

func deletePassword(c *cli.Context) error {
	var (
		config        = configuration.NewConfiguration()
		username      = c.String("username")
		passwordIndex = c.Int("index")
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
	if len(config.Accounts[username].Passwords) > 0 {
		if passwordIndex < 0 || passwordIndex > len(config.Accounts[username].Passwords) {
			return errors.New("invalid password index")
		}

		// Delete indexed item
		account := config.Accounts[username]
		account.Passwords = append(account.Passwords[:passwordIndex], account.Passwords[passwordIndex+1:]...)
		config.Accounts[username] = account

		// Write Configuration
		if err := configuration.WriteConfiguration(config); err != nil {
			return err
		}

		fmt.Println("account password deleted")
	}

	if len(config.Accounts[username].Passwords) == 0 {
		fmt.Println("no passwords set for account " + username)
	}

	return nil
}
