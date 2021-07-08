/*
 * Filename: configuration.go
 * Author: Nathaniel Thomas
 */

package models

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

// AccountConfiguration : Account Profile Configuration
type AccountConfiguration struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"admin"`
}
