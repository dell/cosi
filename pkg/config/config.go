// Copyright © 2023-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

// Package config creates configuration structures based on provided json or yaml file.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/dell/csmlog"

	"gopkg.in/yaml.v3"
)

var log = csmlog.GetLogger()

//go:generate go run github.com/atombender/go-jsonschema/cmd/gojsonschema@v0.13.1 --package=config --output=config.gen.go config.schema.json --extra-imports

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

	log.Debugf("Config file path %s", filename)

	// limit reader is used, so we will read only 20MB of the file.
	maxFileSize := 20000000
	lr := io.LimitReader(f, int64(maxFileSize))

	return io.ReadAll(lr)
}
