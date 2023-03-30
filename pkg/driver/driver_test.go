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

func TestDriver(t *testing.T) {
	t.Parallel()
	for name, test := range map[string]func(t *testing.T){
		"empty configuration":                             TestDriverEmptyConfiguration,
		"configuration with connections":                  TestDriverConfigurationWithConnections,
		"configuration with duplicate ID":                 TestDriverConfigurationWithDuplicateID,
		"configuration with missing objectscale":          TestDriverConfigurationWithMissingObjectscale,
		"with preexisting socket file":                    TestDriverWithPreexistingSocketFile,
		"fail on non-existing socket directory":           TestDriverFailOnNonExistingDirectory,
		"run blocking server":                             TestDriverRunBlockingServer,
		"blocking server configuration with duplicate ID": TestDriverBlockingServerConfigurationWithDuplicateID,
	} {
		t.Run(name, test)
	}
}

func TestDriverEmptyConfiguration(t *testing.T) {
	t.Parallel()
	// TODO: server should start successfully with an empty configuration?
	err := runWithParameters(t, testConfigEmpty, t.TempDir())
	assert.NoError(t, err)
}

func TestDriverConfigurationWithConnections(t *testing.T) {
	t.Parallel()
	// server should start successfully with the provided configuration
	err := runWithParameters(t, testConfigWithConnections, t.TempDir())
	assert.NoError(t, err)
}

func TestDriverConfigurationWithDuplicateID(t *testing.T) {
	t.Parallel()
	// server should fail if the configuration contains duplicate IDs
	err := runWithParameters(t, testConfigDuplicateID, t.TempDir())
	assert.Error(t, err)
}

func TestDriverConfigurationWithMissingObjectscale(t *testing.T) {
	t.Parallel()
	// server should fail if the configuration is missing the objectscale connection
	err := runWithParameters(t, testConfigMissingObjectscale, t.TempDir())
	assert.Error(t, err)
}

func TestDriverWithPreexistingSocketFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	// create a socket file
	socketPath := path.Join(dir, "cosi.sock")
	_, err := os.Create(socketPath)
	assert.NoError(t, err)
	// server should delete the socket file and start with a new one
	err = runWithParameters(t, testConfigWithConnections, dir)
	assert.NoError(t, err)
}

func TestDriverFailOnNonExistingDirectory(t *testing.T) {
	t.Parallel()
	// server should fail if the directory does not exist
	err := runWithParameters(t, testConfigWithConnections, "/nonexistent")
	assert.Error(t, err)
}

func TestDriverRunBlockingServer(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	err := RunBlocking(ctx, testConfigWithConnections, t.TempDir(), "test")

	assert.NoError(t, err)
}

func TestDriverBlockingServerConfigurationWithDuplicateID(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	err := RunBlocking(ctx, testConfigDuplicateID, t.TempDir(), "test")

	assert.Error(t, err)
}

func runWithParameters(t *testing.T, configuration *config.ConfigSchemaJson, socketDirectoryPath string) error {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	testSocketPath := path.Join(socketDirectoryPath, "cosi.sock")
	ready, err := Run(ctx, configuration, testSocketPath, "test")
	if err != nil {
		return err
	}

	// Block until server is ready
	<-ready
	// Cancel context to stop server gracefully
	cancel()

	return nil
}
