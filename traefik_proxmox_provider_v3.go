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
	client       *internal.ProxmoxClient
	cancel       func()
}

// New creates a new Provider plugin.
// This function is compatible with Traefik v3 plugin loading mechanism.
func New(ctx context.Context, config *Config, name string) (*Provider, error) {
	providerConfig := &provider.Config{
		PollInterval:   config.PollInterval,
		ApiEndpoint:    config.ApiEndpoint,
		ApiTokenId:     config.ApiTokenId,
		ApiToken:       config.ApiToken,
		ApiLogging:     config.ApiLogging,
		ApiValidateSSL: config.ApiValidateSSL,
	}

	innerProvider, err := provider.New(ctx, providerConfig, name)
	if err != nil {
		return nil, err
	}

	return &Provider{
		name:         name,
		pollInterval: innerProvider.pollInterval,
		client:       innerProvider.client,
		cancel:       innerProvider.cancel,
	}, nil
}

// Init initializes the provider.
func (p *Provider) Init() error {
	return nil
}

// Provide creates and sends dynamic configuration.
// This method signature is compatible with Traefik v3.
// Implements the provider logic directly to avoid delegation issues with yaegi
func (p *Provider) Provide(cfgChan chan<- json.Marshaler) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initial configuration
	if err := p.updateConfiguration(ctx, cfgChan); err != nil {
		time.Sleep(p.pollInterval)
	}

	// Start polling
	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := p.updateConfiguration(ctx, cfgChan); err != nil {
				// Continue on error
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (p *Provider) updateConfiguration(ctx context.Context, cfgChan chan<- json.Marshaler) error {
	if p.client == nil {
		return errors.New("client not initialized")
	}

	servicesMap, err := getServiceMap(p.client, ctx)
	if err != nil {
		return err
	}

	configuration := generateConfiguration(servicesMap)
	cfgChan <- &dynamic.JSONPayload{Configuration: configuration}
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

// getServiceMap retrieves all services from Proxmox API
func getServiceMap(client *internal.ProxmoxClient, ctx context.Context) (map[string][]internal.Service, error) {
	servicesMap := make(map[string][]internal.Service)

	nodes, err := client.GetNodes(ctx)
	if err != nil {
		return nil, err
	}

	for _, node := range nodes {
		// Get VMs
		vms, err := client.GetVirtualMachines(ctx, node.Name)
		if err != nil {
			continue
		}
		for _, vm := range vms {
			servicesMap[vm.Name] = append(servicesMap[vm.Name], internal.Service{
				Name:     vm.Name,
				Type:     internal.ServiceTypeVM,
				Node:     node.Name,
				Hostname: vm.Hostname,
			})
		}

		// Get Containers
		containers, err := client.GetContainers(ctx, node.Name)
		if err != nil {
			continue
		}
		for _, container := range containers {
			servicesMap[container.Name] = append(servicesMap[container.Name], internal.Service{
				Name:     container.Name,
				Type:     internal.ServiceTypeContainer,
				Node:     node.Name,
				Hostname: container.Hostname,
			})
		}
	}

	return servicesMap, nil
}

// generateConfiguration generates the Traefik configuration
func generateConfiguration(servicesMap map[string][]internal.Service) *dynamic.Configuration {
	config := &dynamic.Configuration{
		Backends: make(map[string]*dynamic.Backend),
	}

	for _, service := range servicesMap {
		for _, svc := range service {
			backendName := svc.Name
			config.Backends[backendName] = &dynamic.Backend{
				LoadBalancer: &dynamic.LoadBalancer{
					Servers: []dynamic.Server{},
				},
			}
		}
	}

	return config
}
