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
	"errors"
	"strings"

	"github.com/dell/csmlog"
	"github.com/dell/goobjectscale/pkg/client/model"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	cosi "sigs.k8s.io/container-object-storage-interface/proto"
)

var log = csmlog.GetLogger()

// DriverCreateBucket is an idempotent method for creating buckets
// It is expected to create the same bucket given a bucketName and protocol
// If the bucket already exists, then it MUST return codes.AlreadyExists
// Return values
//
//	nil -                   Bucket successfully created
//	codes.AlreadyExists -   Bucket already exists. No more retries
//	non-nil err -           Internal error                                [requeue'd with exponential backoff]
func (s *Server) DriverCreateBucket(ctx context.Context,
	req *cosi.DriverCreateBucketRequest,
) (*cosi.DriverCreateBucketResponse, error) {
	ctx, span := otel.Tracer(CreateBucketTraceName).Start(ctx, "DriverCreateBucket")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	log.Infof("Creating Bucket %s", req.GetName())
	createParams := &model.CreateBucketRequestParams{}
	err := createParams.ParseFrom(req.GetParameters())
	if err != nil {
		return nil, logAndTraceError(span, "failed parsing parameters", err, codes.Internal)
	}

	// check rg if user input replicationGroup value
	vPoolID := ""
	if len(createParams.ReplicationGroup) > 0 {
		vPools, err := s.mgmtClient.VPools().List(ctx)
		if err != nil {
			return nil, logAndTraceError(span, "failed listing replication groups", err, codes.Internal)
		}
		rgExists := false
		for _, vPool := range vPools {
			if vPool.Name == createParams.ReplicationGroup {
				vPoolID = vPool.ID
				rgExists = true
				break
			}
		}
		if !rgExists {
			return nil, logAndTraceError(span, "replication group not found", err, codes.NotFound, "replicationGroup", createParams.ReplicationGroup)
		}
	}

	existingBucket, err := getBucket(ctx, s, req.GetName(), map[string]string{"namespace": s.namespace})
	if err != nil && !errors.Is(err, model.ErrParameterNotFound) {
		return nil, logAndTraceError(span, "error finding bucket", err, codes.Internal, "namespace", s.namespace, "bucket", req.GetName())
	} else if err == nil && existingBucket != nil {
		return &cosi.DriverCreateBucketResponse{
			BucketId: strings.Join([]string{s.backendID, existingBucket.Name}, "-"),
		}, nil
	}

	toBeCreatedBucket := &model.ObjectBucketParam{}
	toBeCreatedBucket.Name = req.GetName()
	toBeCreatedBucket.Namespace = s.namespace
	toBeCreatedBucket.Vpool = vPoolID
	toBeCreatedBucket.HeadType = model.S3
	toBeCreatedBucket.EncryptionEnabled = createParams.EncryptionEnabled
	toBeCreatedBucket.FsAccessEnabled = createParams.FilesystemEnabled
	toBeCreatedBucket.IsStaleAllowed = createParams.AccessDuringOutageEnabled
	toBeCreatedBucket.Retention = createParams.DefaultRetention
	toBeCreatedBucket.BlockSize = createParams.QuotaLimit
	toBeCreatedBucket.NotificationSize = createParams.QuotaWarn

	bucket, err := s.mgmtClient.Buckets().Create(ctx, toBeCreatedBucket)
	if err != nil {
		return nil, logAndTraceError(span, "failed to create bucket", err, codes.Internal, "namespace", s.namespace, "bucket", req.GetName())
	}

	log.Infof("Successfully created bucket %s in namespace %s", req.GetName(), s.namespace)
	return &cosi.DriverCreateBucketResponse{BucketId: strings.Join([]string{s.backendID, bucket.Name}, "-")}, nil
}
