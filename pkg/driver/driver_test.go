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

package driver

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/dell/cosi-driver/pkg/config"
	"github.com/stretchr/testify/assert"
)

var (
	testSocketPath string
	testDir        = "test"
	testRegion     = "us-east-1"

	testConfig = &config.ConfigSchemaJson{
		CosiEndpoint: "unix:///tmp/cosi.sock",
		LogLevel:     config.ConfigSchemaJsonLogLevelTrace,
	}

	testConfigNoNetwork = &config.ConfigSchemaJson{
		CosiEndpoint: "/tmp/cosi.sock",
		LogLevel:     config.ConfigSchemaJsonLogLevelTrace,
	}

	testConfigInvalidNetwork = &config.ConfigSchemaJson{
		CosiEndpoint: "tcp:///tmp/cosi.sock",
		LogLevel:     config.ConfigSchemaJsonLogLevelTrace,
	}

	testConfigWithConnections = &config.ConfigSchemaJson{
		CosiEndpoint: "unix:///tmp/cosi.sock",
		LogLevel:     config.ConfigSchemaJsonLogLevelTrace,
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

	testConfigDuplicateID = &config.ConfigSchemaJson{
		CosiEndpoint: "unix:///tmp/cosi.sock",
		LogLevel:     config.ConfigSchemaJsonLogLevelTrace,
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
			{
				Objectscale: &config.Objectscale{
					Credentials: config.Credentials{
						Username: "testuser",
						Password: "testpassword",
					},
					Id:                 "test.id",
					ObjectscaleGateway: "gateway.objectscale.test",
					ObjectstoreGateway: "gateway.objectstore.test",
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

	testConfigMissingObjectscale = &config.ConfigSchemaJson{
		CosiEndpoint: "unix:///tmp/cosi.sock",
		LogLevel:     config.ConfigSchemaJsonLogLevelTrace,
		Connections: []config.Configuration{
			{
				Objectscale: nil,
			},
		},
	}

	testConfigWithoutEndpoint = &config.ConfigSchemaJson{
		LogLevel: config.ConfigSchemaJsonLogLevelTrace,
	}
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name           string
		config         *config.ConfigSchemaJson
		auxiliaryFuncs []func(ctx context.Context) error
		expectedError  bool
	}{
		{
			name:          "success",
			config:        testConfig,
			expectedError: false,
		},
		{
			name:          "success without network in COSI Endpoint",
			config:        testConfigNoNetwork,
			expectedError: false,
		},
		{
			name:          "failure invalid network",
			config:        testConfigInvalidNetwork,
			expectedError: true,
		},
		{
			name:          "success with connections",
			config:        testConfigWithConnections,
			expectedError: false,
		},
		{
			name:          "failure no endpoint",
			config:        testConfigWithoutEndpoint,
			expectedError: true,
		},
		{
			name:          "failure duplicate ID",
			config:        testConfigDuplicateID,
			expectedError: true,
		},
		{
			name:          "failure missing connection config",
			config:        testConfigMissingObjectscale,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test server starts successfully and stops gracefully
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			errCh := make(chan error, 1)
			go func() {
				log.Printf("go Run")
				errCh <- Run(ctx, tc.config, "test")
			}()

			// Wait for server to start
			time.Sleep(500 * time.Millisecond)

			// Cancel context to stop server gracefully
			cancel()

			if tc.expectedError {
				err := <-errCh
				assert.Error(t, err)
			} else {
				err := <-errCh
				assert.NoError(t, err)
			}
		})
	}
}
