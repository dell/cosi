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

// Package identity implements server for handling identity requests to a driver instance.
package identity

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cosi "sigs.k8s.io/container-object-storage-interface-spec"

	l "github.com/dell/cosi/pkg/logger"
)

var ErrEmptyDriverName = errors.New("driver name is empty")

// Server is an implementation of COSI identity server.
type Server struct {
	name string
}

var _ cosi.IdentityServer = (*Server)(nil)

// New returns new server.
func New(provisioner string) *Server {
	return &Server{
		name: provisioner,
	}
}

// DriverGetInfo returns name of server.
func (srv *Server) DriverGetInfo(_ context.Context,
	_ *cosi.DriverGetInfoRequest,
) (*cosi.DriverGetInfoResponse, error) {
	if srv.name == "" {
		l.Log().Error(ErrEmptyDriverName, "driver name is empty")

		return nil, status.Error(codes.InvalidArgument, "DriverName is empty")
	}

	return &cosi.DriverGetInfoResponse{
		Name: srv.name,
	}, nil
}
