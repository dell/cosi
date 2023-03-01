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
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	spec "sigs.k8s.io/container-object-storage-interface-spec"

	"google.golang.org/grpc"

	"github.com/dell/cosi-driver/pkg/identity"
	"github.com/dell/cosi-driver/pkg/provisioner"
)

// Run starts the gRPC server for the identity and provisioner servers
func Run(ctx context.Context, name string, port int) error {
	// Setup identity server and provisioner server
	identityServer := identity.New(name)
	provisionerServer := provisioner.New()
	// Some options for gRPC server may be needed
	options := []grpc.ServerOption{}
	// Crate new gRPC server
	server := grpc.NewServer(options...)
	// Register identity and provisioner servers, so they will handle gRPC requests to the server
	spec.RegisterIdentityServer(server, identityServer)
	spec.RegisterProvisionerServer(server, provisionerServer)

	// Create shared listener for gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

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
