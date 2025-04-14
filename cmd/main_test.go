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

package main

import (
	"context"
	"testing"
)

func TestTracerProvider(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid URL",
			url:     "localhost:4317",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			_, err := tracerProvider(ctx, tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("tracerProvider() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetOtelLogger(t *testing.T) {
	tests := []struct {
		name         string
		configFile   string
		logFormat    string
		logLevel     int
		otelEndpoint string
	}{
		{
			name:         "set OpenTelemetry logger",
			configFile:   "/cosi/config.yaml",
			logFormat:    "text",
			logLevel:     4,
			otelEndpoint: "10.0.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setOtelLogger()
		})
	}
}
