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

	"github.com/dell/goobjectscale/pkg/client/model"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"

	log "github.com/sirupsen/logrus"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

var (
	CreateBucketTraceName = "CreateBucketRequest"

	ErrEmptyBucketName      = errors.New("empty bucket name")
	ErrFailedToCreateBucket = errors.New("failed to create bucket")
)

// DriverCreateBucket creates Bucket on specific Object Storage Platform.
func (s *Server) DriverCreateBucket(
	ctx context.Context,
	req *cosi.DriverCreateBucketRequest,
) (*cosi.DriverCreateBucketResponse, error) {
	ctx, span := otel.Tracer(CreateBucketTraceName).Start(ctx, "ObjectscaleDriverCreateBucket")
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
		return nil, logAndTraceError(log.WithFields(log.Fields{}), span, ErrEmptyBucketName.Error(), ErrEmptyBucketName, codes.InvalidArgument)
	}

	parameters := make(map[string]string)
	parameters["namespace"] = s.namespace

	log.WithFields(log.Fields{
		"parameters": parameters,
	}).Info("parameters of the bucket")

	// Get bucket.
	existingBucket, err := s.getBucket(ctx, bucket.Name, parameters)
	if err != nil && !errors.Is(err, ErrParameterNotFound) {
		fields := log.Fields{
			"bucket": bucket.Name,
		}

		return nil, logAndTraceError(log.WithFields(fields), span, ErrFailedToCheckBucketExists.Error(), err, codes.Internal)
	} else if err == nil && existingBucket != nil {
		return &cosi.DriverCreateBucketResponse{
			BucketId: strings.Join([]string{s.backendID, bucket.Name}, "-"),
		}, nil
	}

	// Create bucket.
	err = s.createBucket(ctx, bucket)
	if err != nil {
		fields := log.Fields{
			"bucket": bucket.Name,
		}

		return nil, logAndTraceError(log.WithFields(fields), span, ErrFailedToCreateBucket.Error(), err, codes.Internal)
	}

	// Return response.
	return &cosi.DriverCreateBucketResponse{
		BucketId: strings.Join([]string{s.backendID, bucket.Name}, "-"),
	}, nil
}

// getBucket is used to obtain bucket info from the Provisioner.
func (s *Server) getBucket(ctx context.Context, bucketName string, parameters map[string]string) (*model.Bucket, error) {
	ctx, span := otel.Tracer(CreateBucketTraceName).Start(ctx, "ObjectscaleGetBucket")
	defer span.End()

	// Check if bucket with specific name and parameters already exists.
	retrievedBucket, err := s.mgmtClient.Buckets().Get(ctx, bucketName, parameters)

	switch {
	// First, we don't find the bucket on the Provider.
	case errors.Is(err, ErrParameterNotFound):
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
		return nil, fmt.Errorf("%w: %w", ErrFailedToCheckBucketExists, err)
	}
}

// createBucket is used to create bucket on the Provisioner.
func (s *Server) createBucket(ctx context.Context, bucket *model.Bucket) error {
	ctx, span := otel.Tracer(CreateBucketTraceName).Start(ctx, "ObjectscaleCreateBucket")
	defer span.End()

	_, err := s.mgmtClient.Buckets().Create(ctx, *bucket)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"bucket": bucket.Name,
	}).Info("bucket successfully created")

	span.AddEvent("bucket successfully created")

	return nil
}
