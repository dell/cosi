// Copyright © 2023-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

// Package provisioner implements main COSI driver functionality
// namely creating/deleting buckets and granting/revoking access to them.
package provisioner

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	otelCodes "go.opentelemetry.io/otel/codes"
	cosi "sigs.k8s.io/container-object-storage-interface/proto"

	"github.com/dell/cosi/pkg/provisioner/objectscale"
	"github.com/dell/csmlog"
)

// Server is an implementation of a provisioner server.
type Server struct {
	driverset *Driverset
	cosi.UnimplementedProvisionerServer
}

var (
	_   cosi.ProvisionerServer = (*Server)(nil)
	log                        = csmlog.GetLogger()
)

const (
	ErrInvalidBackendID = "invalid backend ID"
)

// New initializes Server based on the config file.
func New(driverset *Driverset) *Server {
	return &Server{
		driverset: driverset,
	}
}

// DriverCreateBucket creates Bucket on specific Object Storage Platform.
func (s *Server) DriverCreateBucket(ctx context.Context,
	req *cosi.DriverCreateBucketRequest,
) (*cosi.DriverCreateBucketResponse, error) {
	tracedCtx, span := otel.Tracer(objectscale.CreateBucketTraceName).Start(ctx, "ProvisionerCreateBucket")
	defer span.End()

	id := req.Parameters["id"]

	// get the driver from driverset
	// if there is no correct driver, log error, and return standard error message
	d, err := s.driverset.Get(id)
	if err != nil {
		log.Errorf("%s %s: %v", ErrInvalidBackendID, id, err)
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, ErrInvalidBackendID)

		return nil, status.Error(codes.InvalidArgument, ErrInvalidBackendID)
	}

	// execute DriverCreateBucket from correct driver
	return d.DriverCreateBucket(tracedCtx, req)
}

// DriverDeleteBucket deletes Bucket on specific Object Storage Platform.
func (s *Server) DriverDeleteBucket(ctx context.Context,
	req *cosi.DriverDeleteBucketRequest,
) (*cosi.DriverDeleteBucketResponse, error) {
	tracedCtx, span := otel.Tracer("DeleteBucketRequest").Start(ctx, "ProvisionerDeleteBucket")
	defer span.End()

	id := getID(req.BucketId)

	// get the driver from driverset
	// if there is no correct driver, log error, and return standard error message
	d, err := s.driverset.Get(id)
	if err != nil {
		log.Errorf("%s %s: %v", ErrInvalidBackendID, id, err)
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, ErrInvalidBackendID)

		return nil, status.Error(codes.InvalidArgument, ErrInvalidBackendID)
	}

	// execute DriverDeleteBucket from correct driver
	return d.DriverDeleteBucket(tracedCtx, req)
}

// DriverGrantBucketAccess provides access to Bucket on specific Object Storage Platform.
func (s *Server) DriverGrantBucketAccess(ctx context.Context,
	req *cosi.DriverGrantBucketAccessRequest,
) (*cosi.DriverGrantBucketAccessResponse, error) {
	tracedCtx, span := otel.Tracer("GrantBucketAccessRequest").Start(ctx, "ProvisionerGrantBucketAccess")
	defer span.End()

	id := req.Parameters["id"]

	// get the driver from driverset
	// if there is no correct driver, log error, and return standard error message
	d, err := s.driverset.Get(id)
	if err != nil {
		log.Errorf("%s %s: %v", ErrInvalidBackendID, id, err)
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, ErrInvalidBackendID)

		return nil, status.Error(codes.InvalidArgument, ErrInvalidBackendID)
	}

	// execute DriverGrantBucketAccess from correct driver
	return d.DriverGrantBucketAccess(tracedCtx, req)
}

// DriverRevokeBucketAccess revokes access from Bucket on specific Object Storage Platform.
func (s *Server) DriverRevokeBucketAccess(ctx context.Context,
	req *cosi.DriverRevokeBucketAccessRequest,
) (*cosi.DriverRevokeBucketAccessResponse, error) {
	tracedCtx, span := otel.Tracer("RevokeBucketAccessRequest").Start(ctx, "ProvisionerRevokeBucketAccess")
	defer span.End()

	id := getID(req.BucketId)

	// get the driver from driverset
	// if there is no correct driver, log error, and return standard error message
	d, err := s.driverset.Get(id)
	if err != nil {
		log.Errorf("%s %s: %v", ErrInvalidBackendID, id, err)
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, ErrInvalidBackendID)

		return nil, status.Error(codes.InvalidArgument, ErrInvalidBackendID)
	}

	// execute DriverRevokeBucketAccess from correct driver
	return d.DriverRevokeBucketAccess(tracedCtx, req)
}

// getID splits the string and returns ID from it
// correct format of string is:
// (ID)-(other identifers).
func getID(s string) string {
	id := strings.Split(s, "-")

	correctFormatLength := 2
	if len(id) < correctFormatLength {
		return ""
	}

	return id[0]
}
