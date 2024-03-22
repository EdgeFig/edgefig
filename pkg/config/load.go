package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads config from the specified filename
func LoadConfig(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found at %s", configPath)
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config := &Config{}

	err = yaml.Unmarshal(configBytes, config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling config yaml: %w", err)
	}

	//err = config.expandEnv()
	//if err != nil {
	//	return nil, err
	//}

	return config, nil
}
