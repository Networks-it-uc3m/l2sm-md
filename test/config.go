package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config holds all your configuration parameters.
type Config struct {
	ServerAddress string          `yaml:"serverAddress"`
	NetworkName   string          `yaml:"networkName"`
	Provider      ProviderConfig  `yaml:"provider"`
	Clusters      []ClusterConfig `yaml:"clusters"`
	NetworkType   string          `yaml:"networkType"`
	Namespace     string          `yaml:"namespace"`
}

type ProviderConfig struct {
	Name   string `yaml:"name"`
	Domain string `yaml:"domain"`
}

type ClusterConfig struct {
	Name        string     `yaml:"name"`
	ApiKey      string     `yaml:"apiKey"`
	BearerToken string     `yaml:"bearerToken"`
	Nodes       []string   `yaml:"nodes"`
	GatewayNode NodeConfig `yaml:"gatewayNode"`
}

type NodeConfig struct {
	Name      string `yaml:"name"`
	IPAddress string `yaml:"ipAddress"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
