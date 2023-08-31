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

	l "github.com/dell/cosi/pkg/logger"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
)

// All errors that can be returned by DriverDeleteBucket.
var ErrFailedToDeleteBucket = errors.New("bucket was not successfully deleted")

// DriverDeleteBucket deletes Bucket on specific Object Storage Platform.
func (s *Server) DriverDeleteBucket(ctx context.Context,
	req *cosi.DriverDeleteBucketRequest,
) (*cosi.DriverDeleteBucketResponse, error) {
	ctx, span := otel.Tracer(DeleteBucketTraceName).Start(ctx, "ObjectscaleDriverDeleteBucket")
	defer span.End()

	l.Log().V(4).Info("Bucket id being deleted.", "bucketID", req.BucketId)

	span.AddEvent("bucket is being deleted")

	// Check if bucketID is not empty.
	if req.GetBucketId() == "" {
		return nil, logAndTraceError(span, ErrInvalidBucketID.Error(), ErrInvalidBucketID, codes.InvalidArgument)
	}

	// Extract bucket name from bucketID.
	bucketName, err := GetBucketName(req.BucketId)
	if err != nil {
		return nil, logAndTraceError(span, ErrInvalidBucketID.Error(), err, codes.InvalidArgument, "bucketID", req.BucketId)
	}

	// Delete bucket.
	err = s.mgmtClient.Buckets().Delete(ctx, bucketName, s.namespace, s.emptyBucket)

	if errors.Is(err, ErrParameterNotFound) {
		l.Log().V(0).Info("Bucket does not exist.", "bucket", bucketName)

		span.AddEvent("bucket does not exist")

		return &cosi.DriverDeleteBucketResponse{}, nil
	}

	if err != nil {
		return nil, logAndTraceError(span, ErrFailedToDeleteBucket.Error(), err, codes.Internal, "bucket", bucketName)
	}

	l.Log().V(4).Info("Bucket successfully deleted.", "bucket", bucketName)

	span.AddEvent("bucket successfully deleted")

	return &cosi.DriverDeleteBucketResponse{}, nil
}
