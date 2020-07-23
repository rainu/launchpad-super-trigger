package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
)

func ReadConfig(configReader io.Reader) (*Config, error) {
	result := &Config{}

	err := yaml.NewDecoder(configReader).Decode(result)
	if err != nil {
		return nil, fmt.Errorf("could not decode config: %w", err)
	}

	err = validate(result)
	if err != nil {
		return nil, fmt.Errorf("config is invalid: %w", err)
	}

	return result, nil
}