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

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dell/cosi-driver/pkg"
	"sigs.k8s.io/container-object-storage-interface-provisioner-sidecar/pkg/provisioner"
)

const (
	driverAddress = "unix:///var/lib/cosi/cosi.sock"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT)
		<-sigc
		cancel()
	}()

	if err := run(ctx); err != nil {
		log.Fatalf("error during execution: %s", err.Error())
	}
}

func run(ctx context.Context) error {
	identityServer, bucketProvisioner, err := pkg.NewDriver(ctx, "dell-cosi-driver")
	if err != nil {
		return err
	}

	server, err := provisioner.NewDefaultCOSIProvisionerServer(
		driverAddress,
		identityServer,
		bucketProvisioner)
	if err != nil {
		return err
	}

	return server.Run(ctx)
}
