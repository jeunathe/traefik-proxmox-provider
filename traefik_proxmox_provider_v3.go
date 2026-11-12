// Package traefik_proxmox_provider is a plugin to use a proxmox cluster as a provider for Traefik.
// This is a compatibility version for Traefik v3 that implements a minimal standalone provider.
package traefik_proxmox_provider

import (
	"context"
	"encoding/json"
)

// Config the plugin configuration.
type Config struct {
	PollInterval   string `json:"pollInterval"`
	ApiEndpoint    string `json:"apiEndpoint"`
	ApiTokenId     string `json:"apiTokenId"`
	ApiToken       string `json:"apiToken"`
	ApiLogging     string `json:"apiLogging"`
	ApiValidateSSL string `json:"apiValidateSSL"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		PollInterval:   "30s",
		ApiValidateSSL: "true",
		ApiLogging:     "info",
	}
}

// Provider a plugin.
type Provider struct {
	apiEndpoint    string
	apiTokenId     string
	apiToken       string
	apiValidateSSL bool
	apiLogging     string
}

// New creates a new Provider plugin.
// This function is compatible with Traefik v3 plugin loading mechanism.
func New(ctx context.Context, config *Config, name string) (*Provider, error) {
	return &Provider{
		apiEndpoint:    config.ApiEndpoint,
		apiTokenId:     config.ApiTokenId,
		apiToken:       config.ApiToken,
		apiValidateSSL: config.ApiValidateSSL == "true",
		apiLogging:     config.ApiLogging,
	}, nil
}

// Init initializes the provider.
func (p *Provider) Init() error {
	return nil
}

// Provide creates and sends dynamic configuration.
// This method signature is compatible with Traefik v3.
func (p *Provider) Provide(cfgChan chan<- json.Marshaler) error {
	// Create a minimal empty configuration
	type Configuration struct {
		Backends map[string]*Backend `json:"backends,omitempty"`
	}
	
	type Backend struct {
		LoadBalancer *LoadBalancer `json:"loadBalancer,omitempty"`
	}
	
	type LoadBalancer struct {
		Servers []Server `json:"servers,omitempty"`
	}
	
	type Server struct {
		Address string `json:"address,omitempty"`
	}
	
	config := &Configuration{
		Backends: make(map[string]*Backend),
	}
	
	// Send the configuration
	cfgChan <- config
	
	return nil
}

// Stop the provider.
func (p *Provider) Stop() error {
	return nil
}

// Load is required by Traefik v3 plugin interface
func (p *Provider) Load() {
	// This method is called by Traefik v3 to load the provider
	// Implementation not required for basic functionality
}
