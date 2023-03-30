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

package provisioner

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	log "github.com/sirupsen/logrus"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

// Server is an implementation of a provisioner server.
type Server struct {
	driverset *Driverset
}

var _ cosi.ProvisionerServer = (*Server)(nil)

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
	id := req.Parameters["id"]

	// get the driver from driverset
	// if there is no correct driver, log error, and return standard error message
	d, err := s.driverset.Get(id)
	if err != nil {
		log.WithFields(log.Fields{
			"id":    id,
			"error": err,
		}).Error("DriverCreateBucket: Invalid backend ID")

		return nil, status.Error(codes.InvalidArgument, "DriverCreateBucket: Invalid backend ID")
	}

	// execute DriverCreateBucket from correct driver
	return d.DriverCreateBucket(ctx, req)
}

// DriverDeleteBucket deletes Bucket on specific Object Storage Platform.
func (s *Server) DriverDeleteBucket(ctx context.Context,
	req *cosi.DriverDeleteBucketRequest,
) (*cosi.DriverDeleteBucketResponse, error) {
	id := getID(req.BucketId)

	// get the driver from driverset
	// if there is no correct driver, log error, and return standard error message
	d, err := s.driverset.Get(id)
	if err != nil {
		log.WithFields(log.Fields{
			"id":    id,
			"error": err,
		}).Error("DriverDeleteBucket: Invalid backend ID")

		return nil, status.Error(codes.InvalidArgument, "DriverDeleteBucket: Invalid backend ID")
	}

	// execute DriverDeleteBucket from correct driver
	return d.DriverDeleteBucket(ctx, req)
}

// DriverGrantBucketAccess provides access to Bucket on specific Object Storage Platform.
func (s *Server) DriverGrantBucketAccess(ctx context.Context,
	req *cosi.DriverGrantBucketAccessRequest,
) (*cosi.DriverGrantBucketAccessResponse, error) {
	id := req.Parameters["id"]

	// get the driver from driverset
	// if there is no correct driver, log error, and return standard error message
	d, err := s.driverset.Get(id)
	if err != nil {
		log.WithFields(log.Fields{
			"id":    id,
			"error": err,
		}).Error("DriverGrantBucketAccess: Invalid backend ID")

		return nil, status.Error(codes.InvalidArgument, "DriverGrantBucketAccess: Invalid backend ID")
	}

	// execute DriverGrantBucketAccess from correct driver
	return d.DriverGrantBucketAccess(ctx, req)
}

// DriverRevokeBucketAccess revokes access from Bucket on specific Object Storage Platform.
func (s *Server) DriverRevokeBucketAccess(ctx context.Context,
	req *cosi.DriverRevokeBucketAccessRequest,
) (*cosi.DriverRevokeBucketAccessResponse, error) {
	id := getID(req.BucketId)

	// get the driver from driverset
	// if there is no correct driver, log error, and return standard error message
	d, err := s.driverset.Get(id)
	if err != nil {
		log.WithFields(log.Fields{
			"id":    id,
			"error": err,
		}).Error("DriverRevokeBucketAccess: Invalid backend ID")

		return nil, status.Error(codes.InvalidArgument, "DriverRevokeBucketAccess: Invalid backend ID")
	}

	// execute DriverRevokeBucketAccess from correct driver
	return d.DriverRevokeBucketAccess(ctx, req)
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
