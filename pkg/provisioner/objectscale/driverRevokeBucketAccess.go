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
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/dell/goobjectscale/pkg/client/model"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	log "github.com/sirupsen/logrus"
	otelCodes "go.opentelemetry.io/otel/codes"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

// DriverRevokeBucketAccess revokes access from Bucket on specific Object Storage Platform.
// TODO: this probably has to be refactored in order to meet the gocognit requirements (complexity < 30).
func (s *Server) DriverRevokeBucketAccess(ctx context.Context, // nolint:gocognit
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

	log.WithFields(log.Fields{
		"bucket": bucketName,
	}).Info("bucket access for bucket is being revoked")

	parameters := make(map[string]string)
	parameters["namespace"] = s.namespace

	log.WithFields(log.Fields{
		"parameters": parameters,
	}).Info("parameters of the bucket")

	// Check if bucket for revoking access exists.
	_, err = s.mgmtClient.Buckets().Get(ctx, bucketName, parameters)
	if err != nil && !errors.Is(err, ErrParameterNotFound) {
		errMsg := errors.New("failed to check bucket existence")
		log.WithFields(log.Fields{
			"bucket": bucketName,
			"error":  err,
		}).Error(errMsg.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, errMsg.Error())

		return nil, status.Error(codes.Internal, errMsg.Error())
	} else if err != nil {
		errMsg := errors.New("bucket not found")
		log.WithFields(log.Fields{
			"bucket": bucketName,
			"error":  err,
		}).Error(errMsg.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, errMsg.Error())

		return nil, status.Error(codes.NotFound, errMsg.Error())
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

	// Get existing policy.
	policy, err := s.mgmtClient.Buckets().GetPolicy(ctx, bucketName, parameters)
	if err != nil && !errors.Is(err, model.Error{Code: model.CodeResourceNotFound}) {
		errMsg := errors.New("failed to check bucket policy existence")
		log.WithFields(log.Fields{
			"bucket": bucketName,
		}).Error(errMsg.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, errMsg.Error())

		return nil, status.Error(codes.Internal, errMsg.Error())
	} else if err == nil && policy == "" {
		errMsg := errors.New("policy is empty")
		log.WithFields(log.Fields{
			"bucket": bucketName,
		}).Error(errMsg.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, errMsg.Error())

		return nil, status.Error(codes.Internal, errMsg.Error())
	}

	// Amazon Resource Name, format: arn:aws:s3:<objectScaleID>:<objectStoreID>:<bucketName>/*.
	// To see more: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference-arns.html.
	awsBucketResourceARN := fmt.Sprintf("arn:aws:s3:%s:%s:%s/*", s.objectScaleID, s.objectStoreID, bucketName)
	// Unique ID, format: urn:osc:iam::<namespace>:user/<userName>.
	awsPrincipalString := fmt.Sprintf("urn:osc:iam::%s:user/%s", s.namespace, req.AccountId)

	jsonPolicy := UpdateBucketPolicyRequest{}

	err = json.Unmarshal([]byte(policy), &jsonPolicy)
	if err != nil {
		errMsg := errors.New("failed to marshall policy")
		log.WithFields(log.Fields{
			"bucket":   bucketName,
			"PolicyID": jsonPolicy.PolicyID,
			"error":    err,
		}).Error(errMsg.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, errMsg.Error())

		return nil, status.Error(codes.Internal, errMsg.Error())
	}

	for k, statement := range jsonPolicy.Statement {
		isPrincipal := false
		isResource := false

		for _, p := range statement.Principal.AWS {
			if p == awsPrincipalString {
				isPrincipal = true
			}
		}

		for _, r := range statement.Resource {
			if r == awsBucketResourceARN {
				isResource = true
			}
		}

		if isPrincipal && isResource {
			jsonPolicy.Statement = append(jsonPolicy.Statement[:k], jsonPolicy.Statement[k+1:]...)
		}
	}

	updatedPolicy, err := json.Marshal(jsonPolicy)
	if err != nil {
		errMsg := errors.New("failed to marshal updatePolicy into JSON")
		log.WithFields(log.Fields{
			"bucket":   bucketName,
			"PolicyID": jsonPolicy.PolicyID,
			"error":    err,
		}).Error(errMsg.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, errMsg.Error())

		return nil, status.Error(codes.Internal, errMsg.Error())
	}

	// Update policy.
	err = s.mgmtClient.Buckets().UpdatePolicy(ctx, bucketName, string(updatedPolicy), parameters)
	if err != nil {
		errMsg := errors.New("failed to update bucket policy")
		log.WithFields(log.Fields{
			"bucket": bucketName,
			"policy": updatedPolicy,
			"error":  err,
		}).Error(errMsg.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, errMsg.Error())

		return nil, status.Error(codes.Internal, errMsg.Error())
	}

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

// GetBucketName splits BucketID by -, the first element is backendID, the second element is bucketName.
func GetBucketName(bucketId string) (string, error) {
	list := strings.Split(bucketId, "-")
	if len(list) != 2 {
		return "", errors.New("improper bucketId")
	}
	return list[1], nil

}
