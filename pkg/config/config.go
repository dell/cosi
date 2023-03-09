//Copyright © 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
