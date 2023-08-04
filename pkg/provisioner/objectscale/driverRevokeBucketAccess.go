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

	"github.com/aws/aws-sdk-go/service/iam"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"

	log "github.com/sirupsen/logrus"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"

	"github.com/dell/cosi/pkg/provisioner/policy"
	"github.com/dell/goobjectscale/pkg/client/model"
)

// All errors that can be returned by DriverRevokeBucketAccess.
var (
	ErrEmpyAccountID              = errors.New("empty accountID")
	ErrExistingPolicyIsEmpty      = errors.New("existing policy is empty")
	ErrFailedToUpdateBucketPolicy = errors.New("failed to update bucket policy")
	ErrFailedToListAccessKeys     = errors.New("failed to list access keys")
	ErrFailedToDeleteAccessKey    = errors.New("failed to delete access key")
	ErrFailedToDeleteUser         = errors.New("failed to delete user")
)

// DriverRevokeBucketAccess revokes access from Bucket on specific Object Storage Platform.
// TODO: this probably has to be refactored in order to meet the gocognit requirements (complexity < 30).
func (s *Server) DriverRevokeBucketAccess(ctx context.Context, //nolint:gocognit
	req *cosi.DriverRevokeBucketAccessRequest,
) (*cosi.DriverRevokeBucketAccessResponse, error) {
	ctx, span := otel.Tracer("RevokeBucketAccessRequest").Start(ctx, "ObjectscaleDriverRevokeBucketAccess")
	defer span.End()

	// TODO: modify errors reporting to use new system
	// Check if bucketID is not empty.
	if err := isBucketIDEmpty(req); err != nil {
		return nil, logAndTraceError(log.WithFields(log.Fields{}), span, ErrInvalidBucketID.Error(), err, codes.InvalidArgument)
	}

	// Check if accountID is not empty.
	if err := isAccountIDEmpty(req); err != nil {
		return nil, logAndTraceError(log.WithFields(log.Fields{}), span, ErrEmpyAccountID.Error(), err, codes.InvalidArgument)
	}

	// Extract bucket name from bucketID.
	bucketName, err := GetBucketName(req.GetBucketId())
	if err != nil {
		return nil, logAndTraceError(log.WithFields(log.Fields{}), span, ErrInvalidBucketID.Error(), err, codes.InvalidArgument)
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
	bucketExists := true

	_, err = s.mgmtClient.Buckets().Get(ctx, bucketName, parameters)

	if err != nil && !errors.Is(err, ErrParameterNotFound) {
		fields := log.Fields{
			"bucket": bucketName,
		}

		return nil, logAndTraceError(log.WithFields(fields), span, ErrFailedToCheckBucketExists.Error(), err, codes.Internal)
	} else if err != nil {
		warnMsg := "bucket not found"
		log.WithFields(log.Fields{
			"bucket": bucketName,
			"error":  err,
		}).Warn(warnMsg)
		span.AddEvent("bucket not found")
		bucketExists = false
	}

	// Check user existence.
	userExists := true

	_, err = s.iamClient.GetUser(&iam.GetUserInput{UserName: &req.AccountId})
	if err != nil && err.Error() != iam.ErrCodeNoSuchEntityException {
		fields := log.Fields{
			"user": req.AccountId,
		}

		return nil, logAndTraceError(
			log.WithFields(fields), span, ErrFailedToCheckUserExists.Error(), err, codes.Internal,
		)
	} else if err != nil {
		warnMsg := "user does not exist"
		log.WithFields(log.Fields{
			"user":  req.AccountId,
			"error": err,
		}).Warn(warnMsg)
		span.AddEvent(warnMsg)
		userExists = false
	}

	if bucketExists {

		// handleBucketExistsCase()

		// Get existing policy.
		existingPolicy, err := s.mgmtClient.Buckets().GetPolicy(ctx, bucketName, parameters)
		if err != nil && !errors.Is(err, model.Error{Code: model.CodeResourceNotFound}) {
			fields := log.Fields{
				"bucket": bucketName,
			}

			return nil, logAndTraceError(
				log.WithFields(fields), span, ErrFailedToCheckPolicyExists.Error(), err, codes.Internal,
			)
		} else if err == nil && existingPolicy == "" {
			fields := log.Fields{
				"bucket": bucketName,
			}

			return nil, logAndTraceError(
				log.WithFields(fields), span, ErrExistingPolicyIsEmpty.Error(), err, codes.Internal,
			)
		}

		// Amazon Resource Name, format: arn:aws:s3:<objectScaleID>:<objectStoreID>:<bucketName>/*.
		// To see more: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference-arns.html.
		awsBucketResourceARN := fmt.Sprintf("arn:aws:s3:%s:%s:%s/*", s.objectScaleID, s.objectStoreID, bucketName)
		// Unique ID, format: urn:osc:iam::<namespace>:user/<userName>.
		awsPrincipalString := fmt.Sprintf("urn:osc:iam::%s:user/%s", s.namespace, req.AccountId)

		jsonPolicy := policy.Document{}

		err = json.Unmarshal([]byte(existingPolicy), &jsonPolicy)
		if err != nil {
			fields := log.Fields{
				"bucket":   bucketName,
				"PolicyID": jsonPolicy.ID,
			}

			return nil, logAndTraceError(
				log.WithFields(fields), span, ErrFailedToMarshalPolicy.Error(), err, codes.Internal,
			)
		}

		for k, statement := range jsonPolicy.Statement {
			log.WithFields(log.Fields{
				"k":         k,
				"statement": statement,
			}).Debug("processing next statement")

			statement.Principal.AWS = remove(statement.Principal.AWS, awsPrincipalString)
			statement.Resource = remove(statement.Resource, awsBucketResourceARN)

			jsonPolicy.Statement[k] = statement
		}

		updatedPolicy, err := json.Marshal(jsonPolicy)
		if err != nil {
			fields := log.Fields{
				"bucket":   bucketName,
				"PolicyID": jsonPolicy.ID,
			}

			return nil, logAndTraceError(
				log.WithFields(fields), span, ErrFailedToMarshalPolicy.Error(), err, codes.Internal,
			)
		}

		log.WithFields(log.Fields{
			"policy":    jsonPolicy,
			"rawPolicy": string(updatedPolicy),
		}).Debug("updating policy")

		// Update policy.
		err = s.mgmtClient.Buckets().UpdatePolicy(ctx, bucketName, string(updatedPolicy), parameters)
		if err != nil {
			fields := log.Fields{
				"bucket": bucketName,
				"policy": updatedPolicy,
			}

			return nil, logAndTraceError(
				log.WithFields(fields), span, ErrFailedToUpdateBucketPolicy.Error(), err, codes.Internal,
			)
		}
	}

	if userExists {
		//TODO: move to a separate function

		// Get access keys list.
		accessKeyList, err := s.iamClient.ListAccessKeys(&iam.ListAccessKeysInput{UserName: &req.AccountId})
		if err != nil {
			fields := log.Fields{
				"userName": req.AccountId,
			}

			return nil, logAndTraceError(
				log.WithFields(fields), span, ErrFailedToListAccessKeys.Error(), err, codes.Internal,
			)
		}
		// TODO: THINK THIS THROUGH: move this to a separate function
		// Delete all access keys for particular user.
		for _, accessKey := range accessKeyList.AccessKeyMetadata {
			_, err = s.iamClient.DeleteAccessKey(&iam.DeleteAccessKeyInput{AccessKeyId: accessKey.AccessKeyId, UserName: &req.AccountId})
			if err != nil {
				fields := log.Fields{
					"userName":  req.AccountId,
					"accessKey": accessKey.AccessKeyId,
				}

				return nil, logAndTraceError(
					log.WithFields(fields), span, ErrFailedToDeleteAccessKey.Error(), err, codes.Internal,
				)
			}
		}

		// Delete user.
		_, err = s.iamClient.DeleteUser(&iam.DeleteUserInput{UserName: &req.AccountId})
		if err != nil {
			fields := log.Fields{
				"userName": req.AccountId,
			}

			return nil, logAndTraceError(
				log.WithFields(fields), span, ErrFailedToDeleteUser.Error(), err, codes.Internal,
			)
		}
	}

	log.WithFields(log.Fields{
		"userName": req.AccountId,
		"bucket":   bucketName,
	}).Info("bucket access for bucket is revoked")

	return &cosi.DriverRevokeBucketAccessResponse{}, nil
}

// remove is a generic function that removes all occurrences of an item.
func remove[T comparable](from []T, item T) []T {
	output := make([]T, 0, len(from)) // should be little bit faster if we preallocate capacity

	for _, element := range from {
		if element != item {
			output = append(output, element)
		}
	}

	return output
}

// isAccountIDEmpty checks if Account ID is not empty.
func isAccountIDEmpty(req *cosi.DriverRevokeBucketAccessRequest) error {
	if req.GetAccountId() == "" {
		return ErrEmptyBucketAccessName
	}

	return nil
}
