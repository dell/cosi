package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:generate go run github.com/atombender/go-jsonschema/cmd/gojsonschema@main --package=config --output=config.gen.go config.schema.json

// Any additional functionality or configuration related utils should be placed in this file

func New(filename string) (*ConfigSchemaJson, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read config file: %w", err)
	}

	cfg := &ConfigSchemaJson{}
	if strings.HasSuffix(filename, ".json") {

		err = json.Unmarshal(b, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal config from JSON: %w", err)
		}
	} else if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
		err = yaml.Unmarshal(b, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal config from YAML: %w", err)
		}
	} else {
		return nil, errors.New("file extension unknown")
	}

	return cfg, nil
}
