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
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/dell/cosi/pkg/provisioner/policy"
	"github.com/dell/csmlog"
	"github.com/dell/goobjectscale/pkg/client/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	smithy "github.com/aws/smithy-go"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/status"
	cosi "sigs.k8s.io/container-object-storage-interface/proto"

	otelCodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
)

const (
	PolicySid = "cosi"
)

func GetBucketNameFromID(bucketID string) (string, error) {
	if !strings.Contains(bucketID, "-") {
		return "", fmt.Errorf("invalid bucket id %s", bucketID)
	}
	return strings.SplitN(bucketID, "-", 2)[1], nil
}

func (s *Server) DriverRevokeBucketAccess(ctx context.Context,
	req *cosi.DriverRevokeBucketAccessRequest,
) (*cosi.DriverRevokeBucketAccessResponse, error) {
	ctx, span := otel.Tracer(CreateBucketTraceName).Start(ctx, "DriverRevokeBucketAccess")
	defer span.End()

	bucketName, err := GetBucketNameFromID(req.GetBucketId())
	if err != nil {
		return nil, logAndTraceError(span, "invalid bucket name", err, codes.InvalidArgument)
	}

	log.Infof("Revoking access to bucket %s for user %s", bucketName, req.GetAccountId())
	iamClient, err := s.iamClient(ctx)
	if err != nil {
		return nil, logAndTraceError(span, "failed to create IAM client", err, codes.Internal)
	}

	parameters := map[string]string{}
	parameters["namespace"] = s.namespace

	// Check if bucket for revoking access exists.
	bucketExists, err := checkBucketExistence(ctx, s, bucketName, parameters)
	if err != nil {
		return nil, logAndTraceError(span, "failed checking if bucket exists", err, codes.Internal, "bucket", bucketName)
	}

	// Check user existence.
	userExists, err := checkUserExistence(ctx, iamClient, req.GetAccountId())
	if err != nil {
		return nil, logAndTraceError(span, "failed checking if user exists", err, codes.Internal, "user", req.GetAccountId())
	}

	principalUsername := BuildPrincipalString(req.AccountId, s.namespace)

	if bucketExists {
		err := removeBucketPolicy(ctx, s, bucketName, principalUsername, parameters)
		if err != nil {
			return nil, logAndTraceError(span, "failed removing bucket policy", err, codes.Internal, "bucket", bucketName)
		}
	}

	if userExists {
		if err := deleteUser(ctx, iamClient, req.AccountId); err != nil {
			return nil, logAndTraceError(span, "failed deleting user", err, codes.Internal, "user", req.GetAccountId())
		}
	}

	log.Infof("Revoked access to bucket %s for user %s", bucketName, req.GetAccountId())
	return &cosi.DriverRevokeBucketAccessResponse{}, nil
}

func removeBucketPolicy(
	ctx context.Context,
	s *Server,
	bucketName string,
	principalUsername string,
	parameters map[string]string,
) error {
	// Get existing policy.
	existingPolicy, err := s.mgmtClient.Buckets().GetPolicy(ctx, bucketName, parameters)
	if errors.Is(err, model.ErrParameterNotFound) {
		return nil
	}
	if err != nil {
		return err
	}

	// no policy exists, return
	if existingPolicy == "" {
		return nil
	}

	jsonPolicy := policy.Document{}

	err = json.Unmarshal([]byte(existingPolicy), &jsonPolicy)
	if err != nil {
		return err
	}

	updatedPolicyDoc := policy.Document{}
	updatedPolicyDoc.ID = jsonPolicy.ID
	updatedPolicyDoc.Version = jsonPolicy.Version
	updatedPolicyDoc.Statement = []policy.StatementEntry{}

	for _, statement := range jsonPolicy.Statement {
		isMatch := statement.Sid == PolicySid && statement.Principal["AWS"] == principalUsername
		if !isMatch {
			updatedPolicyDoc.Statement = append(updatedPolicyDoc.Statement, statement)
		}
	}

	updatedPolicy, err := json.Marshal(updatedPolicyDoc)
	if err != nil {
		return &aws.RequestCanceledError{}
	}

	if len(updatedPolicyDoc.Statement) == 0 {
		log.Infof("Deleting policy")
		err = s.mgmtClient.Buckets().DeletePolicy(ctx, bucketName, parameters)
		if errors.Is(err, model.ErrParameterNotFound) {
			return nil
		}
		if err != nil {
			return err
		}
	} else {
		log.Infof("Updating policy %s", updatedPolicyDoc)
		log.Debugf("Raw policy %s", string(updatedPolicy))
		// Update policy.
		err = s.mgmtClient.Buckets().UpdatePolicy(ctx, bucketName, string(updatedPolicy), parameters)
		if errors.Is(err, model.ErrParameterNotFound) {
			return nil
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func checkUserExistence(ctx context.Context, iamClient IAM, accountID string) (bool, error) {
	_, span := otel.Tracer(RevokeBucketAccessTraceName).Start(ctx, "ObjectscaleCheckUserExistence")
	defer span.End()

	_, err := iamClient.GetUser(ctx, &iam.GetUserInput{UserName: &accountID})
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NoSuchEntityException:
				return false, nil
			default:
				return false, err
			}
		}
	}
	return true, nil
}

func getBucket(ctx context.Context, s *Server, bucketName string, parameters map[string]string) (*model.Bucket, error) {
	ctx, span := otel.Tracer(RevokeBucketAccessTraceName).Start(ctx, "ObjectscaleGetBucket")
	defer span.End()

	bucket, err := s.mgmtClient.Buckets().Get(ctx, bucketName, parameters)

	if errors.Is(err, model.ErrParameterNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return bucket, nil
}

func checkBucketExistence(ctx context.Context, s *Server, bucketName string, parameters map[string]string) (bool, error) {
	ctx, span := otel.Tracer(RevokeBucketAccessTraceName).Start(ctx, "ObjectscaleCheckBucketExistence")
	defer span.End()
	bucket, err := getBucket(ctx, s, bucketName, parameters)
	if err != nil {
		return false, err
	}
	return bucket != nil, nil
}

func deleteUser(ctx context.Context, iamClient IAM, accountID string) error {
	// Get access keys list.
	accessKeyList, err := iamClient.ListAccessKeys(ctx, &iam.ListAccessKeysInput{UserName: &accountID})
	if err != nil {
		return err
	}

	// Delete all access keys for particular user.
	for _, accessKey := range accessKeyList.AccessKeyMetadata {
		_, err = iamClient.DeleteAccessKey(ctx, &iam.DeleteAccessKeyInput{
			AccessKeyId: accessKey.AccessKeyId, UserName: &accountID,
		})
		if err != nil {
			return err
		}
	}

	// Delete user.
	_, err = iamClient.DeleteUser(ctx, &iam.DeleteUserInput{UserName: &accountID})
	if err != nil {
		return err
	}

	return nil
}

// kvToFields converts variadic key-value pairs into csmlog.Fields.
// Only string keys are kept; malformed pairs are skipped safely.
func kvToFields(keysAndValues ...any) csmlog.Fields {
	fields := csmlog.Fields{}
	for i := 0; i+1 < len(keysAndValues); i += 2 {
		k, ok := keysAndValues[i].(string)
		if !ok {
			continue
		}
		fields[k] = keysAndValues[i+1]
	}
	return fields
}

// logAndTraceError is a helper function that logs an error with specified fields and records it in a span.
func logAndTraceError(span trace.Span, errMsg string, err error, code codes.Code, keysAndValues ...any) error {
	fields := kvToFields(keysAndValues...)

	// Add the error as a structured field
	if err != nil {
		fields["error"] = err
	}

	log.WithFields(fields).Error(errMsg)
	span.RecordError(err)
	span.SetStatus(otelCodes.Error, errMsg)

	return status.Error(code, errMsg)
}
