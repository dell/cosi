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

// TODO: write documentation comment for driver package
package driver

import (
	"context"
	"errors"
	"fmt"
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
	// COSISocket is a default location of COSI API UNIX socket.
	COSISocket = "/var/lib/cosi/cosi.sock"
)

// Driver structure for storing server and listener instances.
type Driver struct {
	// gRPC server
	server *grpc.Server
	// socket listener
	lis net.Listener
}

// New creates a new driver for COSI API with identity and provisioner servers.
func New(ctx context.Context, config *config.ConfigSchemaJson, socket, name string) (*Driver, error) {
	// Setup identity server and provisioner server
	identityServer := identity.New(name)

	driverset := &provisioner.Driverset{}

	for _, cfg := range config.Connections {
		driver, err := provisioner.NewVirtualDriver(ctx, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to validate provided object storage platform connection: %w", err)
		}

		log.WithFields(log.Fields{
			"driver": driver.ID(),
		}).Debug("configuration for specified object storage platform validated")

		err = driverset.Add(driver)
		if err != nil {
			return nil, fmt.Errorf("failed to add object storage platform configuration: %w", err)
		}

		log.WithFields(log.Fields{
			"driver": driver.ID(),
		}).Debug("new configuration for specified object storage platform added")
	}

	provisionerServer := provisioner.New(driverset)
	// Some options for gRPC server may be needed.
	options := []grpc.ServerOption{}
	// Create new gRPC server.
	server := grpc.NewServer(options...)
	// Register identity and provisioner servers, so they will handle gRPC requests to the driver.
	spec.RegisterIdentityServer(server, identityServer)
	spec.RegisterProvisionerServer(server, provisionerServer)

	// Remove socket file if it already exists
	// so we can start a new driver after crash or pod restart
	if _, err := os.Stat(socket); !errors.Is(err, fs.ErrNotExist) {
		if err := os.RemoveAll(socket); err != nil {
			log.Fatalf("failed to remove socket: %v", err)
		}
	}

	// Create shared listener for gRPC server
	listener, err := net.Listen("unix", socket)
	if err != nil {
		return nil, fmt.Errorf("failed to announce on the local network address: %w", err)
	}

	log.WithFields(log.Fields{
		"socket": socket,
	}).Debug("shared listener created")

	return &Driver{server, listener}, nil
}

// starts the gRPC server and returns a channel that will be closed when it is ready.
func (s *Driver) start(ctx context.Context) <-chan struct{} {
	ready := make(chan struct{})
	go func() {
		close(ready)

		if err := s.server.Serve(s.lis); err != nil {
			log.Fatalf("failed to serve gRPC server: %v", err)
		}
	}()

	return ready
}

// Run starts the gRPC server for the identity and provisioner servers.
// This function will not block and instead will provide channel for checking when the driver is ready.
// Await for context if you want the thread you are running this in to block.
func Run(ctx context.Context, config *config.ConfigSchemaJson, socket, name string) (<-chan struct{}, error) {
	// Create new driver
	driver, err := New(ctx, config, socket, name)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to start gRPC server")

		return nil, err
	}

	log.Info("gRPC server started")

	return driver.start(ctx), nil
}

// RunBlocking is a blocking version of Run.
func RunBlocking(ctx context.Context, config *config.ConfigSchemaJson, socket, name string) error {
	// Create new driver
	driver, err := New(ctx, config, socket, name)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to start gRPC server")

		return err
	}

	log.Info("gRPC server started")
	// Block until driver is ready
	<-driver.start(ctx)

	// Block until context is done
	<-ctx.Done()

	// Gracefully stop the driver
	driver.server.GracefulStop()

	return nil
}
