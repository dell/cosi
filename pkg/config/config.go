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
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:generate go run github.com/atombender/go-jsonschema/cmd/gojsonschema@main --package=config --output=config.gen.go config.schema.json

// New takes filename and returns populated configuration struct.
func New(filename string) (*ConfigSchemaJson, error) {
	if strings.HasSuffix(filename, ".json") {
		b, err := readFile(filename)
		if err != nil {
			return nil, fmt.Errorf("unable to read config file: %w", err)
		}

		return NewJSON(b)
	} else if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
		b, err := readFile(filename)
		if err != nil {
			return nil, fmt.Errorf("unable to read config file: %w", err)
		}

		return NewYAML(b)
	} else {
		return nil, errors.New("file extension unknown")
	}
}

// NewJSON takes array of bytes and unmarshals it, to return populated configuration struct.
// Array of bytes is expected to be in JSON format.
func NewJSON(bytes []byte) (*ConfigSchemaJson, error) {
	cfg := &ConfigSchemaJson{}

	err := json.Unmarshal(bytes, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// NewYAML takes array of bytes and unmarshals it, to return populated configuration struct.
// Array of bytes is expected to be in YAML format.
func NewYAML(bytes []byte) (*ConfigSchemaJson, error) {
	var body map[string]interface{}
	err := yaml.Unmarshal(bytes, &body)
	if err != nil {
		return nil, err
	}

	b, _ := json.Marshal(body)

	return NewJSON(b)
}

func readFile(filename string) ([]byte, error) {
	// ignore G304 error, reasons:
	// - the filetype is validated to be one of (json, yaml, yml)
	// - file is read using LimitReader (max 20MB)
	/* #nosec G304 */
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// limit reader is used, so the we will read only 20MB of the file.
	lr := io.LimitReader(f, 20*1000000)

	return io.ReadAll(lr)
}
