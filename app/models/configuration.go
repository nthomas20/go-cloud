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
	Revision int                             `yaml:"revision"`
	Port     string                          `yaml:"port"`
	Logging  map[string]LogConfiguration     `yaml:"logging"`
	Accounts map[string]AccountConfiguration `yaml:"accounts"`
}

// LogConfiguration : Log Configuration
type LogConfiguration struct {
	Filename string `yaml:"filename"`
	MaxBytes int64  `yaml:"max_bytes"`
	MaxFiles int    `yaml:"max_files"`
}

// PasswordConfiguration : Account Password Configuration
type PasswordConfiguration struct {
	Password    string    `yaml:"password"`
	Description string    `yaml:"description"`
	LastUsed    time.Time `yaml:"last_used"`
}

// AccountConfiguration : Account Profile Configuration
type AccountConfiguration struct {
	Username      string                  `yaml:"username"`
	Email         string                  `yaml:"email"`
	Passwords     []PasswordConfiguration `yaml:"passwords"`
	IsActive      bool                    `yaml:"active"`
	IsAdmin       bool                    `yaml:"admin"`
	RootDirectory string                  `yaml:"root_dir"`
}
