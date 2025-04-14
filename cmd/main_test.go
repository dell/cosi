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
	"fmt"
	"os"
	"syscall"
	"testing"

	"github.com/dell/cosi/pkg/config"
	"github.com/dell/cosi/pkg/driver"
)

var (
	testRegion                = "us-east-1"
	testConfigWithConnections = &config.ConfigSchemaJson{
		Connections: []config.Configuration{
			{
				Objectscale: &config.Objectscale{
					Credentials: config.Credentials{
						Username: "testuser",
						Password: "testpassword",
					},
					Id:                 "test.id",
					ObjectscaleGateway: "gateway.objectscale.test",
					ObjectstoreGateway: "gateway.objectstore.test",
					ObjectscaleId:      "objectscale123",
					ObjectstoreId:      "objectstore123",
					Namespace:          "testnamespace",
					Protocols: config.Protocols{
						S3: &config.S3{
							Endpoint: "s3.objectstore.test",
						},
					},
					Region: &testRegion,
					Tls: config.Tls{
						Insecure: true,
					},
				},
			},
		},
	}
)

func TestRunMain(t *testing.T) {
	tests := []struct {
		name         string
		configFile   string
		otelEndpoint string
		wantErr      bool
	}{
		{
			name:         "valid config without tracing",
			configFile:   "config.json",
			otelEndpoint: "",
			wantErr:      false,
		},
		// {
		// 	name:         "valid config with tracing",
		// 	configFile:   "config.json",
		// 	otelEndpoint: "localhost:4317",
		// 	wantErr:      false,
		// },
		// {
		// 	name:         "invalid config file",
		// 	configFile:   "invalid config",
		// 	otelEndpoint: "",
		// 	wantErr:      true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the flags
			*configFile = tt.configFile
			*otelEndpoint = tt.otelEndpoint

			validConfigContent := `
			{
				"version": "1.0",
				"settings": {
					"key": "value"
				}
			}
			`

			file := createTempConfigFile(t, tt.configFile, validConfigContent)
			defer os.Remove(file)

			// create the socket for the listener
			err := os.MkdirAll(driver.COSISocket, 0755)
			if err != nil {
				t.Errorf("failed to create dir: %v", err)
			}
			defer os.Remove(driver.COSISocket)

			// Channel to signal completion
			done := make(chan struct{})

			// Simulate the server running in a separate goroutine
			go func() {
				<-done
				p, err := os.FindProcess(os.Getpid())
				if err != nil {
					t.Error(err)
				}
				err = p.Signal(syscall.SIGINT)
				if err != nil {
					t.Error(err)
				}
			}()

			fmt.Print("server started")

			// Run the main function
			err = runMain()
			if (err != nil) != tt.wantErr {
				t.Errorf("runMain() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Send the interrupt signal
			close(done)
		})
	}
}

func createTempConfigFile(t *testing.T, name string, content string) string {
	t.Helper()

	tmpFile, err := os.Create(name)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	return tmpFile.Name()
}

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
