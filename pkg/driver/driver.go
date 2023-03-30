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
	"errors"
	"io/fs"
	"net"
	"os"

	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
	spec "sigs.k8s.io/container-object-storage-interface-spec"

	"github.com/dell/cosi-driver/pkg/config"
	"github.com/dell/cosi-driver/pkg/identity"
	"github.com/dell/cosi-driver/pkg/provisioner"
)

const (
	// COSISocket is a default location of COSI API UNIX socket
	COSISocket = "/var/lib/cosi/cosi.sock"
)

// Representation of Driver for COSI API
type Driver struct {
	// gRPC Driver
	server *grpc.Server
	// socket listener
	lis net.Listener
}

// NewDriver creates a new driver for COSI API with identity and provisioner servers
func New(config *config.ConfigSchemaJson, socket, name string) (*Driver, error) {
	// Setup identity Driver and provisioner Driver
	identityServer := identity.New(name)

	driverset := &provisioner.Driverset{}
	for _, cfg := range config.Connections {
		driver, err := provisioner.NewVirtualDriver(cfg)
		if err != nil {
			return nil, err
		}

		err = driverset.Add(driver)
		if err != nil {
			return nil, err
		}
	}

	provisionerServer := provisioner.New(driverset)
	// Some options for gRPC Driver may be needed
	options := []grpc.ServerOption{}
	// Crate new gRPC Driver
	server := grpc.NewServer(options...)
	// Register identity and provisioner Drivers, so they will handle gRPC requests to the Driver
	spec.RegisterIdentityServer(server, identityServer)
	spec.RegisterProvisionerServer(server, provisionerServer)

	// Remove socket file if it already exists
	// so we can start a new Driver after crash or pod restart
	if _, err := os.Stat(socket); !errors.Is(err, fs.ErrNotExist) {
		if err := os.RemoveAll(socket); err != nil {
			log.Fatal(err)
		}
	}

	// Create shared listener for gRPC Driver
	listener, err := net.Listen("unix", socket)
	if err != nil {
		return nil, err
	}

	return &Driver{server, listener}, nil
}

// Start starts the gRPC server and returns a channel that will be closed when it is ready
func (s *Driver) Start(ctx context.Context) <-chan struct{} {
	ready := make(chan struct{})
	go func() {
		close(ready)
		if err := s.server.Serve(s.lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()
	return ready
}

// Run starts the gRPC server for the identity and provisioner servers.
// This function will not block and instead will provide channel for checking when the driver is ready.
// Await for context if you want the thread you are running this in to block.
func Run(ctx context.Context, config *config.ConfigSchemaJson, socket, name string) (<-chan struct{}, error) {
	// Create new driver
	driver, err := New(config, socket, name)
	if err != nil {
		return nil, err
	}

	log.Infoln("gRPC server started")

	return driver.Start(ctx), nil
}

// Blocking version of Run
func RunBlocking(ctx context.Context, config *config.ConfigSchemaJson, socket, name string) error {
	// Create new driver
	driver, err := New(config, socket, name)
	if err != nil {
		return err
	}

	log.Infoln("gRPC server started")
	// Block until driver is ready
	<-driver.Start(ctx)

	// Block until context is done
	<-ctx.Done()

	// Gracefully stop the driver
	driver.server.GracefulStop()

	return nil
}
