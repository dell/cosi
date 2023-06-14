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

package objectscale

import (
	"context"
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"
	otelCodes "go.opentelemetry.io/otel/codes"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"

	"github.com/dell/goobjectscale/pkg/client/model"
	"go.opentelemetry.io/otel"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DriverCreateBucket creates Bucket on specific Object Storage Platform.
func (s *Server) DriverCreateBucket(
	ctx context.Context,
	req *cosi.DriverCreateBucketRequest,
) (*cosi.DriverCreateBucketResponse, error) {
	_, span := otel.Tracer("CreateBucketRequest").Start(ctx, "ObjectscaleDriverCreateBucket")
	defer span.End()

	log.WithFields(log.Fields{
		"bucket": req.GetName(),
	}).Info("bucket is being created")

	span.AddEvent("bucket is being created")

	// Create bucket model.
	bucket := &model.Bucket{}
	bucket.Name = req.GetName()
	bucket.Namespace = s.namespace

	// Check if bucket name is not empty.
	if bucket.Name == "" {
		err := errors.New("empty bucket name")
		log.Error(err.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// Display all request parameters.
	parameters := ""
	parametersCopy := make(map[string]string)

	for key, value := range req.GetParameters() {
		parameters += key + ":" + value + ";"
		parametersCopy[key] = value
	}

	// TODO: is this good way of doing this?
	parametersCopy["namespace"] = s.namespace

	log.WithFields(log.Fields{
		"parameters": parameters,
	}).Info("parameters of the bucket")

	// Remove backendID, as this is not valid parameter for bucket creation in ObjectScale.
	delete(parametersCopy, "backendID")

	// Get bucket.
	existingBucket, err := s.getBucket(ctx, bucket.Name, parametersCopy)
	if err != nil && !errors.Is(err, model.Error{Code: model.CodeParameterNotFound}) {
		return nil, status.Error(codes.Internal, err.Error())
	} else if err == nil && existingBucket != nil {
		return &cosi.DriverCreateBucketResponse{
			BucketId: strings.Join([]string{s.backendID, bucket.Name}, "-"),
		}, nil
	}
	// Create bucket.
	err = s.createBucket(ctx, bucket)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Return response.
	return &cosi.DriverCreateBucketResponse{
		BucketId: strings.Join([]string{s.backendID, bucket.Name}, "-"),
	}, nil
}

// getBucket
func (s *Server) getBucket(ctx context.Context, bucketName string, parameters map[string]string) (*model.Bucket, error) {
	// Check if bucket with specific name and parameters already exists.
	_, span := otel.Tracer("CreateBucketRequest").Start(ctx, "ObjectscaleGetBucket")
	defer span.End()

	retievedBucket, err := s.mgmtClient.Buckets().Get(bucketName, parameters)
	if err != nil && !errors.Is(err, model.Error{Code: model.CodeParameterNotFound}) {
		errMsg := errors.New("failed to check bucket existence")
		log.WithFields(log.Fields{
			"bucket": bucketName,
			"error":  err,
		}).Error(errMsg.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, errMsg.Error())

		return nil, err
	} else if err == nil {
		log.WithFields(log.Fields{
			"bucket": bucketName,
		}).Warn("bucket already exists")

		span.AddEvent("bucket already exists")
		return retievedBucket, nil
	} else {
		return nil, nil
	}
}

// createBucket
func (s *Server) createBucket(ctx context.Context, bucket *model.Bucket) error {

	_, span := otel.Tracer("CreateBucketRequest").Start(ctx, "ObjectscaleCreateBucket")
	defer span.End()

	_, err := s.mgmtClient.Buckets().Create(*bucket)
	if err != nil {
		errMsg := errors.New("failed to create bucket")
		log.WithFields(log.Fields{
			"bucket": bucket.Name,
			"error":  err,
		}).Error(errMsg.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, errMsg.Error())

		return errMsg
	}

	log.WithFields(log.Fields{
		"bucket": bucket.Name,
	}).Info("bucket successfully created")

	span.AddEvent("bucket successfully created")

	return nil
}
