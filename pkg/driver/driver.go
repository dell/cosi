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
	"strings"

	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
	spec "sigs.k8s.io/container-object-storage-interface-spec"

	"github.com/dell/cosi-driver/pkg/config"
	"github.com/dell/cosi-driver/pkg/identity"
	"github.com/dell/cosi-driver/pkg/provisioner"
)

var (
	// ErrNoEndpoint indicates that the COSI endpoint was not provided
	ErrNoEndpoint = errors.New("COSI Endpoint not configured")

	// ErrBadEndpoint indicates that the COSI endpoint had bad format
	ErrBadEndpoint = errors.New("COSI Endpoint has bad format, should be unix://<path> or <path>")
)

// Run starts the gRPC server for the identity and provisioner servers
func Run(ctx context.Context, config *config.ConfigSchemaJson, name string) error {
	// Setup identity server and provisioner server
	identityServer := identity.New(name)

	driverset := &provisioner.Driverset{}
	for _, cfg := range config.Connections {
		driver, err := provisioner.NewVirtualDriver(cfg)
		if err != nil {
			return err
		}

		err = driverset.Add(driver)
		if err != nil {
			return err
		}
	}

	provisionerServer := provisioner.New(driverset)
	// Some options for gRPC server may be needed
	options := []grpc.ServerOption{}
	// Crate new gRPC server
	server := grpc.NewServer(options...)
	// Register identity and provisioner servers, so they will handle gRPC requests to the server
	spec.RegisterIdentityServer(server, identityServer)
	spec.RegisterProvisionerServer(server, provisionerServer)

	if config.CosiEndpoint == "" {
		return ErrNoEndpoint
	}

	var (
		network string
		address string
	)
	connection := strings.Split(config.CosiEndpoint, "://")
	if len(connection) == 2 {
		switch connection[0] {
		case "unix":
			network = connection[0]
			address = connection[1]
		default:
			return ErrBadEndpoint
		}
	} else if len(connection) == 1 {
		network = "unix"
		address = connection[0]
	}

	// Remove socket file if it already exists
	// so we can start a new server after crash or pod restart
	if _, err := os.Stat(address); !errors.Is(err, fs.ErrNotExist) {
		if err := os.RemoveAll(address); err != nil {
			log.Fatal(err)
		}
	}

	// Create shared listener for gRPC server
	lis, err := net.Listen(network, address)
	if err != nil {
		return err
	}

	log.Infoln("gRPC server started")
	// Run gRPC server in a separate goroutine
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	// Wait for context cancellation or server error
	<-ctx.Done()
	server.GracefulStop()

	return nil
}
