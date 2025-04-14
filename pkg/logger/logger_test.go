// Copyright © 2025 Dell Inc. or its subsidiaries. All Rights Reserved.
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

package logger

import (
	"testing"

	"github.com/bombsimon/logrusr/v4"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
)

func TestNewAWSLogger(t *testing.T) {
	tests := []struct {
		name   string
		logger logr.Logger
	}{
		{
			name:   "valid logger",
			logger: logrusr.New(logrus.New()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAWSLogger(tt.logger)
			if got.impl != tt.logger {
				t.Errorf("NewAWSLogger() = %v, want %v", got.impl, tt.logger)
			}
		})
	}
}

func TestAWSLogger_Log(t *testing.T) {
	tests := []struct {
		name           string
		logger         AWSLogger
		keysAndValues  []any
		expectedOutput string
	}{
		{
			name:           "log with key-value pairs",
			logger:         NewAWSLogger(logrusr.New(logrus.New())),
			keysAndValues:  []any{"key1", "value1", "key2", "value2"},
			expectedOutput: "internal logger message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logger.Log(tt.keysAndValues...)
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		level     int
		formatter string
	}{
		{
			name:      "default level and formatter",
			level:     defaultLevel,
			formatter: "text",
		},
		{
			name:      "json formatter",
			level:     defaultLevel,
			formatter: "json",
		},
		{
			name:      "pretty formatter",
			level:     defaultLevel,
			formatter: "pretty",
		},
		{
			name:      "invalid formatter",
			level:     defaultLevel,
			formatter: "invalid",
		},
		{
			name:      "level within range",
			level:     5,
			formatter: "text",
		},
		{
			name:      "level out of range",
			level:     11,
			formatter: "text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			New(tt.level, tt.formatter)
			if log.GetSink() == nil {
				t.Errorf("log should not be nil")
			}
		})
	}
}

func TestLog(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "get logger",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Log(); got.GetSink() == nil {
				t.Errorf("Log() sink should not be nil")
			}
		})
	}
}
