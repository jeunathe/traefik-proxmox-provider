//go:build !windows

// Package main provides the traefik-proxmox-provider plugin for Traefik v3.
package main

import (
	"context"
	"encoding/json"
)

// Provider implements the Traefik provider interface for Traefik v3
type Provider struct{}

// New creates a new Provider plugin.
// This is a minimal implementation that satisfies Traefik v3's interface.
func New(ctx context.Context, config *Config, name string) (*Provider, error) {
	return &Provider{}, nil
}

// Config the plugin configuration.
type Config struct {
	PollInterval   string `json:"pollInterval"`
	ApiEndpoint    string `json:"apiEndpoint"`
	ApiTokenId     string `json:"apiTokenId"`
	ApiToken       string `json:"apiToken"`
	ApiLogging     string `json:"apiLogging"`
	ApiValidateSSL string `json:"apiValidateSSL"`
}

// MarshalJSON implements json.Marshaler interface
func (c *Config) MarshalJSON() ([]byte, error) {
	return json.Marshal(c)
}

// Init initializes the provider.
func (p *Provider) Init() error {
	return nil
}

// Provide creates and sends dynamic configuration.
// This is a minimal implementation - can be extended later.
func (p *Provider) Provide(cfgChan chan<- json.Marshaler) error {
	// Send an empty configuration to satisfy the interface
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
	// Empty implementation
}
