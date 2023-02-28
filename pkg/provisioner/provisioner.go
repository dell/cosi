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
	"os"
	"strings"

	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/api"
	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/model"
	log "github.com/sirupsen/logrus"
	"k8s.io/utils/strings/slices"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

type Server struct {
	mgmtClient api.ClientSet
}

var _ cosi.ProvisionerServer = (*Server)(nil)

// FIXME: this is boilerplate, needs proper constructor
func New() *Server {
	return &Server{}
}

func (s *Server) DriverCreateBucket(ctx context.Context,
	req *cosi.DriverCreateBucketRequest) (*cosi.DriverCreateBucketResponse, error) {

	log.WithFields(log.Fields{
		"bucket": req.GetName(),
	}).Info("Bucket is being created")

	bucket := &model.Bucket{}
	bucket.Name = req.GetName()
	bucket.Namespace = req.Parameters["namespace"]

	// display all parameters
	parameters := ""
	for key, value := range req.GetParameters() {
		parameters += key + ":" + value + ";"
	}

	log.WithFields(log.Fields{
		"parameters": parameters,
	}).Debug("Parameters of the bucket")

	//supported protocols
	protocols := []string{"S3"}

	log.WithFields(log.Fields{
		"protocols": strings.Join(protocols, ","),
	}).Debug("Supported protocols")

	// create bucket only if in field protocol is one of the supported protocols
	if !slices.Contains(protocols, req.GetParameters()["protocol"]) {
		log.WithFields(log.Fields{
			"protocol": req.GetParameters()["protocol"],
		}).Error("DriverCreateBucket: Protocol not supported")
		return nil, status.Error(codes.InvalidArgument, "Protocol not supported")
	}

	bucketExisted, err := s.mgmtClient.Buckets().Get(bucket.Name, req.GetParameters())
	if err != nil {
		log.Error("DriverCreateBucket: Failed to check bucket existence")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// return error status if bucket exists, otherwise create new bucket
	if bucketExisted != nil {
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
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.WithFields(log.Fields{
		"bucket": bucket.Name,
	}).Info("Bucket has been successfully created")

	return &cosi.DriverCreateBucketResponse{
		BucketId: strings.Join([]string{bucket.Name, bucket.Namespace}, "-"),
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
