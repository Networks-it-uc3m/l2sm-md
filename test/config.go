// Copyright 2024 Universidad Carlos III de Madrid
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
