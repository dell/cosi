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

package pkg

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

type IdentityServer struct {
	provisioner string
}

func (id *IdentityServer) DriverGetInfo(ctx context.Context,
	req *cosi.DriverGetInfoRequest) (*cosi.DriverGetInfoResponse, error) {

	if id.provisioner == "" {
		log.Printf("Invalid argument: %v", errors.New("provisioner name cannot be empty"))
		return nil, status.Error(codes.InvalidArgument, "ProvisionerName is empty")
	}

	return &cosi.DriverGetInfoResponse{
		Name: id.provisioner,
	}, nil
}
