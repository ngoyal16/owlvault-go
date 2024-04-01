package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Config represents the configuration for the OwlVault service.
type Config struct {
	Storage struct {
		Type  string `yaml:"type"`
		MySQL struct {
			ConnectionString string `yaml:"connection_string"`
		} `yaml:"mysql"`
		DDB struct {
			Region      string `yaml:"region"`
			TablePrefix string `yaml:"table_prefix"`
		} `yaml:"ddb"`
		// Add other storage types here
	} `yaml:"storage"`
}

// ReadConfig reads configuration from the specified YAML file path provided by the environment variable.
func ReadConfig() (*Config, error) {
	// Get the config file path from the environment variable
	configPath := os.Getenv("OWLVAULT_CONFIG_PATH")
	if configPath == "" {
		return nil, fmt.Errorf("environment variable OWLVAULT_CONFIG_PATH is not set")
	}

	// Get the absolute path of the configuration file
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	// Check if the file exists
	_, err = os.Stat(absPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist at path: %s", absPath)
	}

	// Read YAML configuration file
	configFile, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	// Unmarshal YAML into Config struct
	var config Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		return nil, err
	}

	return &config, nil
}