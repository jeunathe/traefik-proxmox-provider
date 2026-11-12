//go:build !windows

// Package main provides the traefik-proxmox-provider plugin for Traefik v3.
package main

import (
	"context"
	"encoding/json"
)

// Provider implements the Traefik provider interface
type Provider struct {
	apiEndpoint    string
	apiTokenId     string
	apiToken       string
	apiValidateSSL bool
	apiLogging     string
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

// New creates a new Provider plugin.
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
func (p *Provider) Provide(cfgChan chan<- json.Marshaler) error {
	// Create a simple configuration
	config := &Config{
		PollInterval:   "30s",
		ApiValidateSSL: "true",
		ApiLogging:     "info",
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
}
