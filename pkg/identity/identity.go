// Copyright © 2023-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

// Package identity implements server for handling identity requests to a driver instance.
package identity

import (
	"context"

	"github.com/dell/csmlog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cosi "sigs.k8s.io/container-object-storage-interface/proto"
)

// Server is an implementation of COSI identity server.
type Server struct {
	name string
	cosi.UnimplementedIdentityServer
}

var (
	_   cosi.IdentityServer = (*Server)(nil)
	log                     = csmlog.GetLogger()
)

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
		log.Error("driver name is empty")
		return nil, status.Error(codes.InvalidArgument, "DriverName is empty")
	}

	return &cosi.DriverGetInfoResponse{
		Name: srv.name,
	}, nil
}
