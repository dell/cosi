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
	testRegion = "us-east-1"

	testConfigEmpty = &config.ConfigSchemaJson{}

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
	t.Parallel()
	for name, test := range map[string]func(t *testing.T){
		"empty configuration":                    TestRunEmptyConfiguration,
		"configuration with connections":         TestRunConfigurationWithConnections,
		"configuration with duplicate ID":        TestRunConfigurationWithDuplicateID,
		"configuration with missing objectscale": TestRunConfigurationWithMissingObjectscale,
		"with preexisting socket file":           TestRunWithPreexistingSocketFile,
		"fail on non-existing socket directory":  TestRunFailOnNonExistingDirectory,
	} {
		t.Run(name, test)
	}
}

func TestRunEmptyConfiguration(t *testing.T) {
	t.Parallel()
	// TODO: server should start successfully with an empty configuration?
	err := runWithParameters(t, testConfigEmpty, t.TempDir())
	assert.NoError(t, err)
}

func TestRunConfigurationWithConnections(t *testing.T) {
	t.Parallel()
	// server should start successfully with the provided configuration
	err := runWithParameters(t, testConfigWithConnections, t.TempDir())
	assert.NoError(t, err)
}

func TestRunConfigurationWithDuplicateID(t *testing.T) {
	t.Parallel()
	// server should fail if the configuration contains duplicate IDs
	err := runWithParameters(t, testConfigDuplicateID, t.TempDir())
	assert.Error(t, err)
}

func TestRunConfigurationWithMissingObjectscale(t *testing.T) {
	t.Parallel()
	// server should fail if the configuration is missing the objectscale connection
	err := runWithParameters(t, testConfigMissingObjectscale, t.TempDir())
	assert.Error(t, err)
}

func TestRunWithPreexistingSocketFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	//create a socket file
	socketPath := path.Join(dir, "cosi.sock")
	_, err := os.Create(socketPath)
	assert.NoError(t, err)
	// server should delete the socket file and start with a new one
	err = runWithParameters(t, testConfigWithConnections, dir)
	assert.NoError(t, err)
}

func TestRunFailOnNonExistingDirectory(t *testing.T) {
	t.Parallel()
	// server should fail if the directory does not exist
	err := runWithParameters(t, testConfigWithConnections, "/nonexistent")
	assert.Error(t, err)
}

func runWithParameters(t *testing.T, configuration *config.ConfigSchemaJson, socketDirectoryPath string) error {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := make(chan error, 1)
	go func() {
		testSocketPath := path.Join(socketDirectoryPath, "cosi.sock")
		err <- Run(ctx, configuration, testSocketPath, "test")
	}()
	// Wait for server to start
	time.Sleep(500 * time.Millisecond)
	// Cancel context to stop server gracefully
	cancel()
	return <-err
}
