// Package traefik_proxmox_provider is a plugin to use a proxmox cluster as a provider for Traefik.
// This is a compatibility version for Traefik v3.
package traefik_proxmox_provider

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jeunathe/traefik-proxmox-provider/internal"
	"github.com/jeunathe/traefik-proxmox-provider/provider"
	"github.com/traefik/genconf/dynamic"
)

// Config the plugin configuration.
type Config struct {
	PollInterval   string `json:"pollInterval" yaml:"pollInterval" toml:"pollInterval"`
	ApiEndpoint    string `json:"apiEndpoint" yaml:"apiEndpoint" toml:"apiEndpoint"`
	ApiTokenId     string `json:"apiTokenId" yaml:"apiTokenId" toml:"apiTokenId"`
	ApiToken       string `json:"apiToken" yaml:"apiToken" toml:"apiToken"`
	ApiLogging     string `json:"apiLogging" yaml:"apiLogging" toml:"apiLogging"`
	ApiValidateSSL string `json:"apiValidateSSL" yaml:"apiValidateSSL" toml:"apiValidateSSL"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	cfg := provider.CreateConfig()
	return &Config{
		PollInterval:   cfg.PollInterval,
		ApiEndpoint:    cfg.ApiEndpoint,
		ApiTokenId:     cfg.ApiTokenId,
		ApiToken:       cfg.ApiToken,
		ApiLogging:     cfg.ApiLogging,
		ApiValidateSSL: cfg.ApiValidateSSL,
	}
}

// Provider a plugin.
type Provider struct {
	name         string
	pollInterval time.Duration
	cancel       func()
}

// New creates a new Provider plugin.
// This function is compatible with Traefik v3 plugin loading mechanism.
func New(ctx context.Context, config *Config, name string) (*Provider, error) {
	// Validate config
	if config.ApiEndpoint == "" || config.ApiTokenId == "" || config.ApiToken == "" {
		return nil, errors.New("missing required configuration")
	}

	return &Provider{
		name:         name,
		pollInterval: 30 * time.Second,
		cancel:       func() {},
	}, nil
}

// Init initializes the provider.
func (p *Provider) Init() error {
	return nil
}

// Provide creates and sends dynamic configuration.
// This method signature is compatible with Traefik v3.
func (p *Provider) Provide(cfgChan chan<- json.Marshaler) error {
	// Create a simple configuration for now
	config := &dynamic.Configuration{
		Backends: make(map[string]*dynamic.Backend),
	}
	
	// Send the configuration
	cfgChan <- &dynamic.JSONPayload{Configuration: config}
	
	return nil
}

// Stop the provider.
func (p *Provider) Stop() error {
	if p.cancel != nil {
		p.cancel()
	}
	return nil
}

// Load is required by Traefik v3 plugin interface
func (p *Provider) Load() {
	// This method is called by Traefik v3 to load the provider
	// Implementation not required for basic functionality
}
