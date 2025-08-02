// Package core provides core functionality for GoAPI
package core

// APIConfig defines the basic configuration for the API
type APIConfig struct {
	Title       string
	Description string
	Version     string
	BasePath    string
	Host        string
	Schemes     []string
	Debug       bool
}
