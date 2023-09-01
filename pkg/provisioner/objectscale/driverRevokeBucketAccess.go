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

	"github.com/aws/aws-sdk-go/service/iam"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"

	l "github.com/dell/cosi/pkg/logger"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"

	"github.com/dell/cosi/pkg/provisioner/policy"
	"github.com/dell/goobjectscale/pkg/client/model"
)

// All errors that can be returned by DriverRevokeBucketAccess.
var (
	ErrEmptyAccountID             = errors.New("empty accountID")
	ErrExistingPolicyIsEmpty      = errors.New("existing policy is empty")
	ErrFailedToUpdateBucketPolicy = errors.New("failed to update bucket policy")
	ErrFailedToListAccessKeys     = errors.New("failed to list access keys")
	ErrFailedToDeleteAccessKey    = errors.New("failed to delete access key")
	ErrFailedToDeleteUser         = errors.New("failed to delete user")
)

// All warnings that can be returned by DriverRevokeBucketAccess.
var (
	WarnBucketNotFound = "Bucket not found."
	WarnUserNotFound   = "User not found."
)

// DriverRevokeBucketAccess revokes access from Bucket on specific Object Storage Platform.
func (s *Server) DriverRevokeBucketAccess(ctx context.Context,
	req *cosi.DriverRevokeBucketAccessRequest,
) (*cosi.DriverRevokeBucketAccessResponse, error) {
	ctx, span := otel.Tracer(RevokeBucketAccessTraceName).Start(ctx, "ObjectscaleDriverRevokeBucketAccess")
	defer span.End()

	// Check if bucketID is not empty.
	if err := isBucketIDEmpty(req); err != nil {
		return nil, logAndTraceError(span, ErrInvalidBucketID.Error(), err, codes.InvalidArgument)
	}

	// Check if accountID is not empty.
	if err := isAccountIDEmpty(req); err != nil {
		return nil, logAndTraceError(span, ErrEmptyAccountID.Error(), err, codes.InvalidArgument)
	}

	// Extract bucket name from bucketID.
	bucketName, err := GetBucketName(req.GetBucketId())
	if err != nil {
		return nil, logAndTraceError(span, ErrInvalidBucketID.Error(), err, codes.InvalidArgument)
	}

	l.Log().V(4).Info("Bucket access for bucket is being revoked.")

	parameters := make(map[string]string)
	parameters["namespace"] = s.namespace

	l.Log().V(4).Info("Parameters of the bucket.", "parameters", parameters)

	// Check if bucket for revoking access exists.
	bucketExists, err := checkBucketExistence(ctx, s, bucketName, parameters)
	if err != nil {
		return nil, logAndTraceError(span, ErrFailedToCheckBucketExists.Error(), err, codes.Internal, "bucket", bucketName)
	}

	// Check user existence.
	userExists, err := checkUserExistence(ctx, s, req.AccountId)
	if err != nil {
		return nil, logAndTraceError(span, err.Error(), err, codes.Internal, "bucket", bucketName, "user", req.AccountId)
	}

	if bucketExists {
		err := removeBucketPolicy(ctx, s, bucketName, parameters)
		if err != nil {
			return nil, logAndTraceError(span, err.Error(), err, codes.Internal, "bucket", bucketName)
		}
	}

	if userExists {
		if err := deleteUser(s, req.AccountId); err != nil {
			return nil, logAndTraceError(span, err.Error(), err, codes.Internal, "bucket", bucketName, "user", req.AccountId)
		}
	}

	l.Log().V(4).Info("Bucket access revoked.", "userName", req.AccountId, "bucket", bucketName)

	return &cosi.DriverRevokeBucketAccessResponse{}, nil
}

// checkUserExistence checks if particual user exists on ObjectScale;
// function returns boolean value: true if user exists, false if user does not exist.
func checkUserExistence(ctx context.Context, s *Server, accountID string) (bool, error) {
	_, span := otel.Tracer(RevokeBucketAccessTraceName).Start(ctx, "ObjectscaleCheckUserExistence")
	defer span.End()

	_, err := s.iamClient.GetUser(&iam.GetUserInput{UserName: &accountID})

	// User is not found - return false. It's a valid scenario.
	if err != nil && err.Error() == iam.ErrCodeNoSuchEntityException {
		l.Log().V(0).Info(WarnUserNotFound, "user", accountID)
		span.AddEvent(WarnUserNotFound)

		return false, nil
	}

	// Connection error probably, failed to check if user exists - return error.
	if err != nil {
		return false, ErrFailedToCheckUserExists
	}

	// No errors - user exists.
	return true, nil
}

// checkBucketExistence checks if particual bucket exists on ObjectScale;
// function returns boolean value: true if bucket exists, false if bucket does not exist.
func checkBucketExistence(ctx context.Context, s *Server, bucketName string, parameters map[string]string) (bool, error) {
	ctx, span := otel.Tracer(RevokeBucketAccessTraceName).Start(ctx, "ObjectscaleCheckBucketExistence")
	defer span.End()

	_, err := s.mgmtClient.Buckets().Get(ctx, bucketName, parameters)

	// Bucket is not found - return false. It's a valid scenario.
	if errors.Is(err, ErrParameterNotFound) {
		span.AddEvent(WarnBucketNotFound)
		l.Log().V(0).Info(WarnBucketNotFound, "bucket", bucketName)

		return false, nil
	}

	// Connection error probably, failed to check if bucket exists - return error.
	if err != nil {
		return false, ErrFailedToCheckBucketExists
	}

	// No errors - bucket exists.
	return true, nil
}

// removeBucketPolicy is a function used when revoking a bucket access;
// it's responsible for updating bucket policy and removing particular right from it.
func removeBucketPolicy(
	ctx context.Context,
	s *Server,
	bucketName string,
	parameters map[string]string,
) error {
	// Get existing policy.
	existingPolicy, err := s.mgmtClient.Buckets().GetPolicy(ctx, bucketName, parameters)
	if err != nil && !errors.Is(err, model.Error{Code: model.CodeResourceNotFound}) {
		return ErrFailedToCheckPolicyExists
	} else if err == nil && existingPolicy == "" {
		return ErrExistingPolicyIsEmpty
	}

	// Amazon Resource Name, format: arn:aws:s3:<objectScaleID>:<objectStoreID>:<bucketName>/*.
	// To see more: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference-arns.html.
	awsBucketResourceARN := BuildResourceString(s.objectScaleID, s.objectStoreID, bucketName)
	// Unique ID, format: urn:osc:iam::<namespace>:user/<userName>.
	awsPrincipalString := BuildPrincipalString(s.namespace, bucketName)

	jsonPolicy := policy.Document{}

	err = json.Unmarshal([]byte(existingPolicy), &jsonPolicy)
	if err != nil {
		return ErrFailedToMarshalPolicy
	}

	for k, statement := range jsonPolicy.Statement {
		l.Log().V(6).Info("Processing next statement.", "k", k, "statement", statement)

		statement.Principal.AWS = remove(statement.Principal.AWS, awsPrincipalString)
		statement.Resource = remove(statement.Resource, awsBucketResourceARN)

		jsonPolicy.Statement[k] = statement
	}

	updatedPolicy, err := json.Marshal(jsonPolicy)
	if err != nil {
		return ErrFailedToMarshalPolicy
	}

	l.Log().V(6).Info("Updating policy.", "policy", jsonPolicy, "rawPolicy", string(updatedPolicy))

	// Update policy.
	err = s.mgmtClient.Buckets().UpdatePolicy(ctx, bucketName, string(updatedPolicy), parameters)
	if err != nil {
		return ErrFailedToUpdateBucketPolicy
	}

	return nil
}

// removeUser is a function used when revoking a bucket access;
// it's responsible for removing users tied to a specific Bucket Access through an accountID.
func deleteUser(s *Server, accountID string) error {
	// Get access keys list.
	accessKeyList, err := s.iamClient.ListAccessKeys(&iam.ListAccessKeysInput{UserName: &accountID})
	if err != nil {
		return ErrFailedToListAccessKeys
	}

	// Delete all access keys for particular user.
	for _, accessKey := range accessKeyList.AccessKeyMetadata {
		_, err = s.iamClient.DeleteAccessKey(&iam.DeleteAccessKeyInput{
			AccessKeyId: accessKey.AccessKeyId, UserName: &accountID,
		})
		if err != nil {
			return ErrFailedToDeleteAccessKey
		}
	}

	// Delete user.
	_, err = s.iamClient.DeleteUser(&iam.DeleteUserInput{UserName: &accountID})
	if err != nil {
		return ErrFailedToDeleteUser
	}

	return nil
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
