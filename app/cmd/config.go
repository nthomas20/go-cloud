/*
 * Filename: config.go
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

func getConfig(c *cli.Context) error {
	var (
		config = configuration.NewConfiguration()
		key    = strings.ToLower(c.String("key"))
	)

	// Read Configuration
	bootstrap.SetupConfiguration()
	if err := configuration.ReadConfiguration(config); err != nil {
		return err
	}

	getables := map[string]interface{}{
		"port":                 config.Port,
		"logging.app.filename": config.Logging["app"].Filename,
	}

	if _, found := getables[key]; found {
		fmt.Println(key+":", getables[key])
	} else {
		return errors.New("invalid or unavailable configuration key: " + key)
	}

	return nil
}

func setConfig(c *cli.Context) error {
	var (
		config = configuration.NewConfiguration()
		key    = strings.ToLower(c.String("key"))
		value  = c.String("value")
	)

	// Read Configuration
	bootstrap.SetupConfiguration()
	if err := configuration.ReadConfiguration(config); err != nil {
		return err
	}

	switch key {
	case "port":
		{
			config.Port = value
			configuration.WriteConfiguration(config)
		}
	default:
		return errors.New("invalid or unavailable configuration key: " + key)
	}

	return nil
}
