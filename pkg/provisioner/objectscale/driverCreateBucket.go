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
	"fmt"
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

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

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

	parameters := make(map[string]string)
	parameters["namespace"] = s.namespace

	log.WithFields(log.Fields{
		"parameters": parameters,
	}).Info("parameters of the bucket")

	// Get bucket.
	existingBucket, err := s.getBucket(ctx, bucket.Name, parameters)
	if err != nil && !errors.Is(err, model.Error{Code: model.CodeParameterNotFound}) {
		msg := "failed to check if bucket exists"

		log.WithFields(log.Fields{
			"error": err,
		}).Error(msg)

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, msg)

		return nil, status.Error(codes.Internal, msg)

	} else if err == nil && existingBucket != nil {
		return &cosi.DriverCreateBucketResponse{
			BucketId: strings.Join([]string{s.backendID, bucket.Name}, "-"),
		}, nil
	}

	// Create bucket.
	err = s.createBucket(ctx, bucket)
	if err != nil {
		msg := "failed to create bucket"

		log.WithFields(log.Fields{
			"error": err,
		}).Error(msg)

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, msg)

		return nil, status.Error(codes.Internal, msg)
	}

	// Return response.
	return &cosi.DriverCreateBucketResponse{
		BucketId: strings.Join([]string{s.backendID, bucket.Name}, "-"),
	}, nil
}

// getBucket is used to obtain bucket info from the Provisioner.
func (s *Server) getBucket(ctx context.Context, bucketName string, parameters map[string]string) (*model.Bucket, error) {
	_, span := otel.Tracer("CreateBucketRequest").Start(ctx, "ObjectscaleGetBucket")
	defer span.End()

	// Check if bucket with specific name and parameters already exists.
	retrievedBucket, err := s.mgmtClient.Buckets().Get(ctx, bucketName, parameters)

	switch {
	// First, we don't find the bucket on the Provider.
	case errors.Is(err, model.Error{Code: model.CodeParameterNotFound}):
		return nil, nil

	// Second case is the error is nil, which means we actually found a bucket.
	case err == nil:
		log.WithFields(log.Fields{
			"bucket": bucketName,
		}).Info("bucket already exists")

		span.AddEvent("bucket already exists")

		return retrievedBucket, nil

	// Final case, when we receive an unknown error.
	default:
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "failed to check bucket existence")

		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}
}

// createBucket is used to create bucket on the Provisioner.
func (s *Server) createBucket(ctx context.Context, bucket *model.Bucket) error {
	_, span := otel.Tracer("CreateBucketRequest").Start(ctx, "ObjectscaleCreateBucket")
	defer span.End()

	_, err := s.mgmtClient.Buckets().Create(ctx, *bucket)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "failed to create bucket")

		return fmt.Errorf("failed to create bucket: %w", err)
	}

	log.WithFields(log.Fields{
		"bucket": bucket.Name,
	}).Info("bucket successfully created")

	span.AddEvent("bucket successfully created")

	return nil
}
