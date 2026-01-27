// Copyright © 2023-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

package objectscale

import (
	"context"

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	cosi "sigs.k8s.io/container-object-storage-interface/proto"
)

func (s *Server) DriverDeleteBucket(ctx context.Context,
	req *cosi.DriverDeleteBucketRequest,
) (*cosi.DriverDeleteBucketResponse, error) {
	ctx, span := otel.Tracer(CreateBucketTraceName).Start(ctx, "DriverDeleteBucket")
	defer span.End()

	bucketName, err := GetBucketNameFromID(req.GetBucketId())
	if err != nil {
		return nil, logAndTraceError(span, "invalid bucket name", err, codes.InvalidArgument)
	}
	log.Infof("Deleting Bucket %s", bucketName)
	parameters := map[string]string{}
	parameters["namespace"] = s.namespace

	bucketExists, err := checkBucketExistence(ctx, s, bucketName, parameters)
	if err != nil {
		return nil, logAndTraceError(span, "failed checking if bucket exists", err, codes.Internal, "bucket", bucketName)
	}

	if bucketExists {
		err := s.mgmtClient.Buckets().Delete(ctx, bucketName, parameters)
		if err != nil {
			return nil, logAndTraceError(span, "failed deleting bucket", err, codes.Internal, "bucket", bucketName)
		}
	} else {
		log.Warnf("Bucket %s does not exist", bucketName)
	}

	log.Infof("Deleted Bucket %s", bucketName)
	return &cosi.DriverDeleteBucketResponse{}, nil
}
