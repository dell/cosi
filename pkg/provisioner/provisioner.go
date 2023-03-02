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
	"os"
	"strings"

	_ "github.com/emcecs/objectscale-management-go-sdk/pkg/client/fake"
	log "github.com/sirupsen/logrus"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"

	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/api"
	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

type Server struct {
	mgmtClient api.ClientSet
	backendID  string
	namespace  string
}

var _ cosi.ProvisionerServer = (*Server)(nil)

// FIXME: this is boilerplate, needs proper constructor
func New(mgmtClient api.ClientSet, backendID, namespace string) *Server {
	// TODO: fill all required fields
	return &Server{
		mgmtClient: mgmtClient,
		backendID:  backendID,
		namespace:  namespace,
	}
}

func (s *Server) ID() string {
	return s.backendID
}

func (s *Server) DriverCreateBucket(ctx context.Context,
	req *cosi.DriverCreateBucketRequest) (*cosi.DriverCreateBucketResponse, error) {

	log.WithFields(log.Fields{
		"bucket": req.GetName(),
	}).Info("Bucket is being created")

	bucket := &model.Bucket{}
	bucket.Name = req.GetName()
	bucket.Namespace = s.namespace

	if bucket.Name == "" {
		log.Error("DriverCreateBucket: Empty bucket name")
		return nil, status.Error(codes.InvalidArgument, "Empty bucket name")
	}

	// display all parameters
	parameters := ""
	parametersCopy := make(map[string]string)
	for key, value := range req.GetParameters() {
		parameters += key + ":" + value + ";"
		parametersCopy[key] = value
	}

	log.WithFields(log.Fields{
		"parameters": parameters,
	}).Debug("Parameters of the bucket")

	// remove backendID, as this is not valid parameter for bucket creation in ObjectScale
	delete(parametersCopy, "backendID")

	_, err := s.mgmtClient.Buckets().Get(bucket.Name, parametersCopy)
	if err != nil && errors.Is(err, &model.Error{Code: 404}) == false {
		log.WithFields(log.Fields{
			"existing_bucket": bucket.Name,
		}).Error("DriverCreateBucket: Failed to check bucket existence")
		return nil, status.Error(codes.Internal, err.Error())
	} else if err == nil {
		log.WithFields(log.Fields{
			"existing_bucket": bucket.Name,
		}).Error("DriverCreateBucket: Bucket already exists")
		return nil, status.Error(codes.AlreadyExists, "Bucket already exists")
	}

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

	return &cosi.DriverCreateBucketResponse{
		BucketId: strings.Join([]string{s.backendID, bucket.Name}, "-"),
	}, nil
}

func (s *Server) DriverDeleteBucket(ctx context.Context,
	req *cosi.DriverDeleteBucketRequest) (*cosi.DriverDeleteBucketResponse, error) {

	return nil, status.Error(codes.Unimplemented, "DriverCreateBucket: not implemented")
}

func (s *Server) DriverGrantBucketAccess(ctx context.Context,
	req *cosi.DriverGrantBucketAccessRequest) (*cosi.DriverGrantBucketAccessResponse, error) {

	return nil, status.Error(codes.Unimplemented, "DriverCreateBucket: not implemented")
}

func (s *Server) DriverRevokeBucketAccess(ctx context.Context,
	req *cosi.DriverRevokeBucketAccessRequest) (*cosi.DriverRevokeBucketAccessResponse, error) {

	return nil, status.Error(codes.Unimplemented, "DriverCreateBucket: not implemented")
}
