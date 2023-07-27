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

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	log "github.com/sirupsen/logrus"
	otelCodes "go.opentelemetry.io/otel/codes"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

// DriverDeleteBucket deletes Bucket on specific Object Storage Platform.
func (s *Server) DriverDeleteBucket(ctx context.Context,
	req *cosi.DriverDeleteBucketRequest,
) (*cosi.DriverDeleteBucketResponse, error) {
	_, span := otel.Tracer("DeleteBucketRequest").Start(ctx, "ObjectscaleDriverDeleteBucket")
	defer span.End()

	log.WithFields(log.Fields{
		"bucketID": req.BucketId,
	}).Info("bucket is being deleted")

	span.AddEvent("bucket is being deleted")

	// Check if bucketID is not empty.
	if req.GetBucketId() == "" {
		err := errors.New("empty bucketID")
		log.Error(err.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Extract bucket name from bucketID.
	bucketName, err := GetBucketName(req.BucketId)
	if err != nil {
		log.WithFields(log.Fields{
			"bucketID": req.BucketId,
			"error":    err,
		}).Error(err.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Delete bucket.
	err = s.mgmtClient.Buckets().Delete(ctx, bucketName, s.namespace, s.emptyBucket)

	if errors.Is(err, ErrParameterNotFound) {
		log.WithFields(log.Fields{
			"bucket": bucketName,
		}).Warn("bucket does not exist")

		span.AddEvent("bucket does not exist")

		return &cosi.DriverDeleteBucketResponse{}, nil
	}

	if err != nil {
		log.WithFields(log.Fields{
			"bucket": bucketName,
			"error":  err,
		}).Error("failed to delete bucket")

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "failed to delete bucket")

		return nil, status.Error(codes.Internal, "bucket was not successfully deleted")
	}

	log.WithFields(log.Fields{
		"bucket": bucketName,
	}).Info("bucket successfully deleted")

	span.AddEvent("bucket successfully deleted")

	return &cosi.DriverDeleteBucketResponse{}, nil
}
