// Copyright © 2023-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

package driver

import (
	"context"
	"errors"
	"net"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dell/cosi/pkg/config"
	"github.com/dell/cosi/pkg/provisioner"
	"github.com/dell/cosi/pkg/provisioner/objectscale"
	"github.com/dell/cosi/pkg/provisioner/virtualdriver"
)

var (
	testRegion    = "us-east-1"
	testNamespace = "testnamespace"

	testConfigEmpty = &config.ConfigSchemaJson{}

	testConfigWithConnections = &config.ConfigSchemaJson{
		Connections: []config.Configuration{
			{
				Objectscale: &config.Objectscale{
					Credentials: config.Credentials{
						Username: "testuser",
						Password: "testpassword",
					},
					Id:           "test.id",
					MgmtEndpoint: "https://example.com/api/s3",
					Namespace:    &testNamespace,
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
					Id:           "test.id",
					MgmtEndpoint: "https://example.com/api/s3",
					Namespace:    &testNamespace,
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
					Id:           "test.id",
					MgmtEndpoint: "https://example.com/api/s3",
					Namespace:    &testNamespace,
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

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestDriver(t *testing.T) {
	os.Setenv("KUBERNETES_SERVICE_HOST", "127.0.0.1")
	os.Setenv("KUBERNETES_SERVICE_PORT", "6443")

	for name, test := range map[string]func(t *testing.T){
		"empty configuration":                             testDriverEmptyConfiguration,
		"configuration with missing objectscale":          testDriverConfigurationWithMissingObjectscale,
		"configuration with connections":                  testDriverConfigurationWithConnections,
		"configuration with duplicate ID":                 testDriverConfigurationWithDuplicateID,
		"with preexisting socket file":                    testDriverWithPreexistingSocketFile,
		"fail on non-existing socket directory":           testDriverFailOnNonExistingDirectory,
		"run blocking server":                             testDriverRunBlockingServer,
		"blocking server configuration with duplicate ID": testDriverBlockingServerConfigurationWithDuplicateID,
		"fail on Listen error":                            testDriverFailOnListen,
	} {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			test(t)
		})
	}
}

func testDriverEmptyConfiguration(t *testing.T) {
	// TODO: server should start successfully with an empty configuration?
	err := runWithParameters(t, testConfigEmpty, t.TempDir())
	assert.NoError(t, err)
}

func testDriverConfigurationWithConnections(t *testing.T) {
	// server should start successfully with the provided configuration

	defer func() {
		ProvisionerNewVirtualDriverFunc = provisioner.NewVirtualDriver
	}()

	ProvisionerNewVirtualDriverFuncMock := func(_ config.Configuration) (virtualdriver.Driver, error) {
		return &objectscale.Server{}, nil
	}

	ProvisionerNewVirtualDriverFunc = ProvisionerNewVirtualDriverFuncMock

	err := runWithParameters(t, testConfigWithConnections, t.TempDir())
	assert.NoError(t, err)
}

func testDriverConfigurationWithDuplicateID(t *testing.T) {
	// server should fail if the configuration contains duplicate IDs

	defer func() {
		ProvisionerNewVirtualDriverFunc = provisioner.NewVirtualDriver
	}()

	ProvisionerNewVirtualDriverFuncMock := func(_ config.Configuration) (virtualdriver.Driver, error) {
		return &objectscale.Server{}, nil
	}

	ProvisionerNewVirtualDriverFunc = ProvisionerNewVirtualDriverFuncMock

	err := runWithParameters(t, testConfigDuplicateID, t.TempDir())
	assert.Error(t, err)
}

func testDriverConfigurationWithMissingObjectscale(t *testing.T) {
	// server should fail if the configuration is missing the objectscale connection
	err := runWithParameters(t, testConfigMissingObjectscale, t.TempDir())
	assert.Error(t, err)
}

func testDriverWithPreexistingSocketFile(t *testing.T) {
	dir := t.TempDir()
	// create a socket file
	socketPath := path.Join(dir, "cosi.sock")
	_, err := os.Create(socketPath)
	assert.NoError(t, err)

	defer func() {
		ProvisionerNewVirtualDriverFunc = provisioner.NewVirtualDriver
	}()

	ProvisionerNewVirtualDriverFuncMock := func(_ config.Configuration) (virtualdriver.Driver, error) {
		return &objectscale.Server{}, nil
	}

	ProvisionerNewVirtualDriverFunc = ProvisionerNewVirtualDriverFuncMock

	// server should delete the socket file and start with a new one
	err = runWithParameters(t, testConfigWithConnections, dir)
	assert.NoError(t, err)
}

func testDriverFailOnNonExistingDirectory(t *testing.T) {
	// server should fail if the directory does not exist
	err := runWithParameters(t, testConfigWithConnections, "/nonexistent")
	assert.Error(t, err)
}

func testDriverRunBlockingServer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	defer func() {
		ProvisionerNewVirtualDriverFunc = provisioner.NewVirtualDriver
	}()

	ProvisionerNewVirtualDriverFuncMock := func(_ config.Configuration) (virtualdriver.Driver, error) {
		return &objectscale.Server{}, nil
	}

	ProvisionerNewVirtualDriverFunc = ProvisionerNewVirtualDriverFuncMock

	err := RunBlocking(ctx, testConfigWithConnections, t.TempDir(), "test")
	assert.NoError(t, err)
}

func testDriverBlockingServerConfigurationWithDuplicateID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := RunBlocking(ctx, testConfigDuplicateID, t.TempDir(), "test")
	assert.Error(t, err)
}

func testDriverFailOnListen(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	defer func() {
		ProvisionerNewVirtualDriverFunc = provisioner.NewVirtualDriver
		NetListenFunc = net.Listen
	}()

	ProvisionerNewVirtualDriverFuncMock := func(_ config.Configuration) (virtualdriver.Driver, error) {
		return &objectscale.Server{}, nil
	}

	NetListenFuncMock := func(_, _ string) (net.Listener, error) {
		return nil, errors.New("injected listen error for uT")
	}

	ProvisionerNewVirtualDriverFunc = ProvisionerNewVirtualDriverFuncMock
	NetListenFunc = NetListenFuncMock

	err := RunBlocking(ctx, testConfigWithConnections, t.TempDir(), "test")
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
