package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

func ReadConfig(configReader ...io.Reader) (*Config, error) {
	result := &Config{
		General: General{
			StartBrightness: 100,
		},
	}

	for i, reader := range configReader {
		err := yaml.NewDecoder(reader).Decode(result)
		if err != nil {
			if f, isFile := reader.(*os.File); isFile {
				return nil, fmt.Errorf("could not decode config [%s]: %w", f.Name(), err)
			}

			return nil, fmt.Errorf("could not decode config [%d]: %w", i, err)
		}
	}

	err := validate(result)
	if err != nil {
		return nil, fmt.Errorf("config is invalid: %w", err)
	}

	return result, nil
}
