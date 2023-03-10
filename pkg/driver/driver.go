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
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"

	"google.golang.org/grpc"

	objectscaleRest "github.com/emcecs/objectscale-management-go-sdk/pkg/client/rest"
	objectscaleClient "github.com/emcecs/objectscale-management-go-sdk/pkg/client/rest/client"
	log "github.com/sirupsen/logrus"
	spec "sigs.k8s.io/container-object-storage-interface-spec"

	"github.com/dell/cosi-driver/pkg/identity"
	"github.com/dell/cosi-driver/pkg/provisioner/objectscale"
)

var (
	objectscaleGateway  = flag.String("objectscale-gateway", "https://localhost:9443", "ObjectScale Gateway")
	objectscaleUser     = flag.String("objectscale-user", "admin", "ObjectScale User")
	objectscalePassword = flag.String("objectscale-password", "admin", "ObjectScale Password")
	objectstoreGateway  = flag.String("objectstore-gateway", "https://localhost:9443", "ObjectStore Gateway")
	unsafeClient        = flag.Bool("unsafe-client", false, "Use unsafe client")
)

// Run starts the gRPC server for the identity and provisioner servers
func Run(ctx context.Context, name, backendID, namespace string, port int) error {
	// Setup identity server and provisioner server
	identityServer := identity.New(name)

	/* #nosec */
	// ObjectScale clientset
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// FIXME: not validating if client should be secure
	unsafeClient := &http.Client{Transport: transport}

	objectscaleAuthUser := objectscaleClient.AuthUser{
		Gateway:  *objectscaleGateway,
		Username: *objectscaleUser,
		Password: *objectscalePassword,
	}
	mngtClient := objectscaleRest.NewClientSet(
		&objectscaleClient.Simple{
			Endpoint:       *objectstoreGateway,
			Authenticator:  &objectscaleAuthUser,
			HTTPClient:     unsafeClient,
			OverrideHeader: false,
		},
	)

	provisionerServer := provisioner.New(mngtClient, backendID, namespace)
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
