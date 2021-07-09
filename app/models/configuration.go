/*
 * Filename: configuration.go
 * Author: Nathaniel Thomas
 */

package models

import (
	"time"
)

// Configuration : Application Configuration File Sructure
// Revision 1
type Configuration struct {
	Revision int                             `json:"revision"`
	Logging  map[string]LogConfiguration     `json:"logging"`
	Accounts map[string]AccountConfiguration `json:"accounts"`
}

// LogConfiguration : Log Configuration
type LogConfiguration struct {
	Filename string `json:"filename"`
	MaxBytes int64  `json:"max_bytes"`
	MaxFiles int    `json:"max_files"`
}

// PasswordConfiguration : Account Password Configuration
type PasswordConfiguration struct {
	Password    string    `json:"password"`
	Description string    `json:"description"`
	LastUsed    time.Time `json:"last_used"`
}

// AccountConfiguration : Account Profile Configuration
type AccountConfiguration struct {
	Username      string                  `json:"username"`
	Email         string                  `json:"email"`
	Passwords     []PasswordConfiguration `json:"passwords"`
	IsActive      bool                    `json:"active"`
	IsAdmin       bool                    `json:"admin"`
	RootDirectory string                  `json:"root_dir"`
}
