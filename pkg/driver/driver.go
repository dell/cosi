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

// Package driver implements gRPC server for handling requests to the COSI driver
// as specified by COSI specification.
package driver

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"os"

	"google.golang.org/grpc"

	spec "sigs.k8s.io/container-object-storage-interface-spec"

	"github.com/dell/cosi/pkg/config"
	"github.com/dell/cosi/pkg/identity"
	l "github.com/dell/cosi/pkg/logger"
	"github.com/dell/cosi/pkg/provisioner"
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
func New(config *config.ConfigSchemaJson, socket, name string) (*Driver, error) {
	// Setup identity server and provisioner server
	identityServer := identity.New(name)

	driverset := &provisioner.Driverset{}

	for _, cfg := range config.Connections {
		driver, err := provisioner.NewVirtualDriver(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to validate provided object storage platform connection: %w", err)
		}

		l.Log().V(6).Info("Configuration for specified object storage platform validated.", "driver", driver.ID())

		err = driverset.Add(driver)
		if err != nil {
			return nil, fmt.Errorf("failed to add object storage platform configuration: %w", err)
		}

		l.Log().V(6).Info("New configuration for specified object storage platform added.", "driver", driver.ID())
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
			l.Log().Error(err, "failed to remove socket")
			os.Exit(1)
		}
	}

	// Create shared listener for gRPC server
	fmt.Println("socket = ", socket)
	listener, err := net.Listen("unix", socket)
	if err != nil {
		return nil, fmt.Errorf("failed to announce on the local network address: %w", err)
	}

	l.Log().V(6).Info("Shared listener created.", "socket", socket)

	fmt.Println("HERE IN NEW")
	return &Driver{server, listener}, nil
}

// starts the gRPC server and returns a channel that will be closed when it is ready.
func (s *Driver) start(_ context.Context) <-chan struct{} {
	ready := make(chan struct{})
	go func() {
		close(ready)

		if err := s.server.Serve(s.lis); err != nil {
			l.Log().Error(err, "failed to serve gRPC server")
			os.Exit(1)
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
		l.Log().Error(err, "failed to start gRPC server")

		return nil, err
	}

	l.Log().V(4).Info("gRPC server started.")

	return driver.start(ctx), nil
}

// RunBlocking is a blocking version of Run.
func RunBlocking(ctx context.Context, config *config.ConfigSchemaJson, socket, name string) error {
	// Create new driver
	fmt.Println("socket in RunBlocking: ", socket)
	driver, err := New(config, socket, name)
	if err != nil {
		l.Log().Error(err, "failed to start gRPC server")
		return err
	}

	l.Log().V(4).Info("gRPC server started.")
	// Block until driver is ready
	<-driver.start(ctx)
	fmt.Print("in RunBlocking \n")

	// Block until context is done
	<-ctx.Done()

	// Gracefully stop the driver
	driver.server.GracefulStop()

	return nil
}
