/*
 * Filename: configuration.go
 * Author: Nathaniel Thomas
 */

package jobs

import (
	"time"

	"gitea.nthomas20.net/nathaniel/go-cloud/app/configuration"
	"gitea.nthomas20.net/nathaniel/go-cloud/app/models"
)

// RefreshConfiguration : Start configuration refresh
func RefreshConfiguration(config *models.Configuration, frequency time.Duration) {
	go func() {
		for {
			time.Sleep(frequency)

			configuration.ReadConfiguration(config)
		}
	}()
}
