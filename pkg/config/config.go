// Copyright © 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
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
	"path"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

//go:generate go run github.com/atombender/go-jsonschema/cmd/gojsonschema@main --package=config --output=config.gen.go config.schema.json

// New takes filename and returns populated configuration struct.
func New(filename string) (*ConfigSchemaJson, error) {
	ext := path.Ext(filename)
	switch ext {
	case ".json":
		b, err := readFile(filename)
		if err != nil {
			return nil, fmt.Errorf("unable to read config file: %w", err)
		}

		return NewJSON(b)
	case ".yaml", ".yml":
		b, err := readFile(filename)
		if err != nil {
			return nil, fmt.Errorf("unable to read config file: %w", err)
		}

		return NewYAML(b)
	default:
		return nil, errors.New("invalid file extension, should be .json, .yaml or .yml")
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

	log.Debug("JSON document unmarshalled")

	return cfg, nil
}

// NewYAML takes array of bytes and unmarshals it, to return populated configuration struct.
// Array of bytes is expected to be in YAML format.
func NewYAML(bytes []byte) (*ConfigSchemaJson, error) {
	// we need to unmarshall it into simple map[string]interface{}
	// config structure does not have custom UnmarshallYAML fields, so the validation is not performed
	var body map[string]interface{}

	err := yaml.Unmarshal(bytes, &body)
	if err != nil {
		return nil, err
	}

	log.Debug("YAML document unmarshalled")
	// we ignore the error, as the config was previously successfully Unmarshaled from YAML.
	// and there is no case, when the Marshaling will fail.
	b, _ := json.Marshal(body)

	// after we marshalled it to JSON, we need to call NewJSON func
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

	log.WithFields(log.Fields{
		"config_file_path": filename,
	}).Debug("config file opened")

	// limit reader is used, so the we will read only 20MB of the file.
	maxFileSize := 20000000
	lr := io.LimitReader(f, int64(maxFileSize))

	return io.ReadAll(lr)
}
