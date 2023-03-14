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

package provisioner

import (
	"context"
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"

	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/api"
	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server is an implementation of a provisioner server.
type Server struct {
	mgmtClient api.ClientSet
	backendID  string
	namespace  string
}

var _ cosi.ProvisionerServer = (*Server)(nil)

// New initializs Server based on the config file.
func New(mgmtClient api.ClientSet, backendID, namespace string) *Server {
	return &Server{
		mgmtClient: mgmtClient,
		backendID:  backendID,
		namespace:  namespace,
	}
}

// ID extends COSI interface by adding ID method.
func (s *Server) ID() string {
	return s.backendID
}

// DriverCreateBucket creates Bucket on specific Object Storage Platform.
func (s *Server) DriverCreateBucket(ctx context.Context,
	req *cosi.DriverCreateBucketRequest) (*cosi.DriverCreateBucketResponse, error) {

	log.WithFields(log.Fields{
		"bucket": req.GetName(),
	}).Info("Bucket is being created")

	// Create bucket model.
	bucket := &model.Bucket{}
	bucket.Name = req.GetName()
	bucket.Namespace = s.namespace

	// Check if bucket name is not empty.
	if bucket.Name == "" {
		log.Error("DriverCreateBucket: Empty bucket name")
		return nil, status.Error(codes.InvalidArgument, "Empty bucket name")
	}

	// Display all request parameters.
	parameters := ""
	parametersCopy := make(map[string]string)
	for key, value := range req.GetParameters() {
		parameters += key + ":" + value + ";"
		parametersCopy[key] = value
	}

	log.WithFields(log.Fields{
		"parameters": parameters,
	}).Debug("Parameters of the bucket")

	// Remove backendID, as this is not valid parameter for bucket creation in ObjectScale.
	delete(parametersCopy, "backendID")

	// Check if bucket with specific name and parameters already exists.
	_, err := s.mgmtClient.Buckets().Get(bucket.Name, parametersCopy)
	if err != nil && errors.Is(err, model.Error{Code: 1004}) == false {
		log.WithFields(log.Fields{
			"existing_bucket": bucket.Name,
		}).Error("DriverCreateBucket: Failed to check bucket existence")
		return nil, status.Error(codes.Internal, "An unexpected error occurred")
	} else if err == nil {
		log.WithFields(log.Fields{
			"existing_bucket": bucket.Name,
		}).Error("DriverCreateBucket: Bucket already exists")
		return nil, status.Error(codes.AlreadyExists, "Bucket already exists")
	}

	// Create bucket.
	bucket, err = s.mgmtClient.Buckets().Create(*bucket)
	if err != nil {
		log.WithFields(log.Fields{
			"bucket": bucket.Name,
		}).Error("DriverCreateBucket: Bucket was not sucessfully created")
		return nil, status.Error(codes.Internal, "Bucket was not sucessfully created")
	}

	log.WithFields(log.Fields{
		"bucket": bucket.Name,
	}).Info("DriverCreateBucket: Bucket has been successfully created")

	// Return response.
	return &cosi.DriverCreateBucketResponse{
		BucketId: strings.Join([]string{s.backendID, bucket.Name}, "-"),
	}, nil
}

// DriverDeleteBucket deletes Bucket on specific Object Storage Platform.
func (s *Server) DriverDeleteBucket(ctx context.Context,
	req *cosi.DriverDeleteBucketRequest) (*cosi.DriverDeleteBucketResponse, error) {

	return nil, status.Error(codes.Unimplemented, "DriverCreateBucket: not implemented")
}

// DriverGrantBucketAccess provides access to Bucket on specific Object Storage Platform.
func (s *Server) DriverGrantBucketAccess(ctx context.Context,
	req *cosi.DriverGrantBucketAccessRequest) (*cosi.DriverGrantBucketAccessResponse, error) {

	return nil, status.Error(codes.Unimplemented, "DriverCreateBucket: not implemented")
}

// DriverRevokeBucketAccess revokes access from Bucket on specific Object Storage Platform.
func (s *Server) DriverRevokeBucketAccess(ctx context.Context,
	req *cosi.DriverRevokeBucketAccessRequest) (*cosi.DriverRevokeBucketAccessResponse, error) {

	return nil, status.Error(codes.Unimplemented, "DriverCreateBucket: not implemented")
}
