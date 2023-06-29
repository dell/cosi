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

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aws/aws-sdk-go/service/iam"
	log "github.com/sirupsen/logrus"
	otelCodes "go.opentelemetry.io/otel/codes"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

// DriverRevokeBucketAccess revokes access from Bucket on specific Object Storage Platform.
func (s *Server) DriverRevokeBucketAccess(ctx context.Context,
	req *cosi.DriverRevokeBucketAccessRequest,
) (*cosi.DriverRevokeBucketAccessResponse, error) {
	ctx, span := otel.Tracer("RevokeBucketAccessRequest").Start(ctx, "ObjectscaleDriverRevokeBucketAccess")
	defer span.End()

	// Check if bucketID is not empty.
	if req.GetBucketId() == "" {
		err := errors.New("empty bucketID")
		log.Error(err.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Check if accountID is not empty.
	if req.GetAccountId() == "" {
		err := errors.New("empty accountID")
		log.Error(err.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Extract bucket name from bucketID.
	bucketName := strings.SplitN(req.BucketId, "-", splitNumber)[1]

	log.WithFields(log.Fields{
		"bucket": bucketName,
	}).Info("bucket access for bucket is being revoked")

	parameters := make(map[string]string)
	parameters["namespace"] = s.namespace

	log.WithFields(log.Fields{
		"parameters": parameters,
	}).Info("parameters of the bucket")

	// Check if bucket for granting access exists.
	_, err := s.mgmtClient.Buckets().Get(ctx, bucketName, parameters)
	if err != nil && !errors.Is(err, ErrParameterNotFound) {
		log.WithFields(log.Fields{
			"bucket": bucketName,
			"error":  err,
		}).Error("failed to check bucket existence")

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "failed to check bucket existence")

		return nil, status.Error(codes.Internal, "an unexpected error occurred")
	} else if err != nil {
		log.WithFields(log.Fields{
			"bucket": bucketName,
			"error":  err,
		}).Error("bucket not found")

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "bucket not found")

		return nil, status.Error(codes.NotFound, "bucket not found")
	}

	// Check user existence.
	_, err = s.iamClient.GetUser(&iam.GetUserInput{UserName: &req.AccountId})
	if err != nil && err.Error() != iam.ErrCodeNoSuchEntityException {
		errMsg := errors.New("failed to check for user existence")
		log.WithFields(log.Fields{
			"user":  req.AccountId,
			"error": err,
		}).Error(errMsg.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, errMsg.Error())

		return nil, status.Error(codes.Internal, errMsg.Error())
	} else if err != nil {
		errMsg := errors.New("failed to get user")
		log.WithFields(log.Fields{
			"user":  req.AccountId,
			"error": err,
		}).Error(errMsg.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, errMsg.Error())

		return nil, status.Error(codes.Internal, errMsg.Error())
	}

	// Get access keys list.
	accessKeyList, err := s.iamClient.ListAccessKeys(&iam.ListAccessKeysInput{UserName: &req.AccountId})
	if err != nil {
		errMsg := errors.New("failed to get access key list")
		log.WithFields(log.Fields{
			"userName": req.AccountId,
			"error":    err,
		}).Error(errMsg.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, errMsg.Error())

		return nil, status.Error(codes.Internal, errMsg.Error())
	}

	// Delete all access keys for particular user.
	for _, accessKey := range accessKeyList.AccessKeyMetadata {
		_, err = s.iamClient.DeleteAccessKey(&iam.DeleteAccessKeyInput{AccessKeyId: accessKey.AccessKeyId, UserName: &req.AccountId})
		if err != nil {
			errMsg := errors.New("failed to delete access key")
			log.WithFields(log.Fields{
				"userName":  req.AccountId,
				"accessKey": accessKey.AccessKeyId,
				"error":     err,
			}).Error(errMsg.Error())

			span.RecordError(err)
			span.SetStatus(otelCodes.Error, errMsg.Error())

			return nil, status.Error(codes.Internal, errMsg.Error())
		}
	}

	// delete/update policy
	// awsBucketResourceARN := fmt.Sprintf("arn:aws:s3:%s:%s:%s/*", s.objectScaleID, s.objectStoreID, bucketName)
	// awsPrincipalString := fmt.Sprintf("urn:osc:iam::%s:user/%s", s.namespace, req.AccountId)

	// Delete user.
	_, err = s.iamClient.DeleteUser(&iam.DeleteUserInput{UserName: &req.AccountId})
	if err != nil {
		errMsg := errors.New("failed to delete user")
		log.WithFields(log.Fields{
			"userName": req.AccountId,
			"error":    err,
		}).Error(errMsg.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, errMsg.Error())

		return nil, status.Error(codes.Internal, errMsg.Error())
	}

	log.WithFields(log.Fields{
		"userName": req.AccountId,
		"bucket":   bucketName,
	}).Info("bucket access for bucket is revoked")

	return &cosi.DriverRevokeBucketAccessResponse{}, nil
}
