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
	"os"
	"path"
	"testing"
	"time"

	"github.com/dell/cosi-driver/pkg/config"
	"github.com/stretchr/testify/assert"
)

var (
	testDir    = "test"
	testRegion = "us-east-1"

	testConfig = &config.ConfigSchemaJson{}

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
		Connections: []config.Configuration{
			{
				Objectscale: nil,
			},
		},
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
			name:          "success with connections",
			config:        testConfigWithConnections,
			expectedError: false,
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

			dir, err := os.MkdirTemp("", testDir)
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)

			errCh := make(chan error, 1)
			go func() {
				testSocketPath := path.Join(dir, "cosi.sock")

				errCh <- Run(ctx, tc.config, testSocketPath, "test")
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

func TestRunWithPreexistingSocketFile(t *testing.T) {
	// Test server starts successfully and stops gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dir, err := os.MkdirTemp("", testDir)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	errCh := make(chan error, 1)
	go func() {
		testSocketPath := path.Join(dir, "cosi.sock")

		// Create preexisting socket file
		_, err := os.Create(testSocketPath)
		if err != nil {
			errCh <- err
		}

		errCh <- Run(ctx, testConfig, testSocketPath, "test")
	}()

	// Wait for server to start
	time.Sleep(500 * time.Millisecond)

	// Cancel context to stop server gracefully
	cancel()

	err = <-errCh
	assert.NoError(t, err)
}

func TestRunFailWithPathErrorOnRemoveAll(t *testing.T) {
	// Test server starts successfully and stops gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dir, err := os.MkdirTemp("", testDir)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	errCh := make(chan error, 1)
	go func() {
		testSocketPath := path.Join(dir, "cosi.sock")

		// Create preexisting socket file
		_, err := os.Create(testSocketPath)
		if err != nil {
			errCh <- err
		}

		// Remove directory to cause error on RemoveAll
		err = os.RemoveAll(dir)
		if err != nil {
			errCh <- err
		}

		errCh <- Run(ctx, testConfig, testSocketPath, "test")
	}()

	// Wait for server to start
	time.Sleep(500 * time.Millisecond)

	// Cancel context to stop server gracefully
	cancel()

	err = <-errCh
	assert.Error(t, err)
}
