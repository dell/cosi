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
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	smithy "github.com/aws/smithy-go"

	cosi "sigs.k8s.io/container-object-storage-interface/proto"
)

func (s *Server) DriverGrantBucketAccess(ctx context.Context,
	req *cosi.DriverGrantBucketAccessRequest,
) (*cosi.DriverGrantBucketAccessResponse, error) {
	ctx, span := otel.Tracer(CreateBucketTraceName).Start(ctx, "DriverGrantBucketAccess")
	defer span.End()

	// Get bucket name from bucketID.
	bucketName, err := GetBucketNameFromID(req.GetBucketId())
	if err != nil {
		return nil, logAndTraceError(span, "invalid bucket name", err, codes.InvalidArgument)
	}

	log.Infof("Creating Bucket Access %s for bucket %s", req.Name, bucketName)
	iamClient, err := s.iamClient(ctx)
	if err != nil {
		return nil, logAndTraceError(span, "failed getting IAM client", err, codes.Internal, "bucket", bucketName)
	}
	// Construct common parameters for bucket requests.
	parameters := make(map[string]string)
	parameters["namespace"] = s.namespace

	log.Debugf("Bucket parameters: %v", parameters)
	// Check if bucket for granting access exists.
	bucketExists, err := checkBucketExistence(ctx, s, bucketName, parameters)
	if err != nil {
		return nil, logAndTraceError(span, "failed checking if bucket exists", err, codes.Internal, "bucket", bucketName)
	}

	if !bucketExists {
		return nil, logAndTraceError(span, "bucket not found", err, codes.NotFound, "bucket", bucketName)
	}

	// This flow below will check for user existence; if user does not exist, it will create one. It will only fail
	// in case of an unknown error, e.g. network issues, to adhere to idempotency requirement.
	userName := BuildUsername(s.namespace, req.Name)
	var user *types.User
	result, err := iamClient.GetUser(ctx, &iam.GetUserInput{
		UserName: aws.String(userName),
	})
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NoSuchEntityException:
				err = nil
			default:
				return nil, logAndTraceError(span, "failed getting user", err, codes.Internal, "user", userName)
			}
		}
	} else {
		user = result.User
	}

	// Check if IAM user exists.
	if user != nil {
		// Case when user exists.
		log.Warnf("User %s already exists", userName)
	} else {
		// Case when user does not exist - create one.
		user, err := iamClient.CreateUser(ctx, &iam.CreateUserInput{
			UserName: &userName,
		})
		if err != nil {
			return nil, logAndTraceError(span, "failed creating user", err, codes.Internal, "user", userName)
		}
		log.Infof("Created ObjectScale IAM user %s with ID %v", userName, user.User.UserId)
	}

	// Check if policy for a specific bucket exists.
	existingPolicy, err := s.mgmtClient.Buckets().GetPolicy(ctx, bucketName, parameters)
	if err != nil {
		return nil, logAndTraceError(span, "failed getting bucket policy", err, codes.Internal, "bucket", bucketName)
	}

	policyRequest := policy.Document{}
	if existingPolicy != "" {
		err = json.NewDecoder(strings.NewReader(existingPolicy)).Decode(&policyRequest)
		if err != nil {
			return nil, logAndTraceError(span, "error parsing policy", err, codes.Internal, "bucket", bucketName)
		}
	}

	// Update policy.
	awsBucketResourceARNs := BuildResourceStrings(bucketName)
	awsPrincipalString := BuildPrincipalString(userName, s.namespace)

	policyRequest.Statement = parsePolicyStatement(
		ctx, policyRequest.Statement, awsBucketResourceARNs, awsPrincipalString,
	)

	log.Debugf("Policy request details: awsBucketResourceARNs: %v, awsPrincipalString: %v, statement: %v", awsBucketResourceARNs, awsPrincipalString, policyRequest.Statement)
	if policyRequest.Version == "" {
		policyRequest.Version = bucketVersion
	}

	if policyRequest.ID == "" {
		policyRequest.ID = "bucket-policy"
	}

	// Marshal the struct to JSON to confirm JSON validity.
	updateBucketPolicyJSON, err := json.Marshal(policyRequest)
	if err != nil {
		return nil, logAndTraceError(span, "error marshalling policy", err, codes.Internal, "bucket", bucketName)
	}

	err = s.mgmtClient.Buckets().UpdatePolicy(ctx, bucketName, string(updateBucketPolicyJSON), parameters)
	if err != nil {
		return nil, logAndTraceError(span, "error updating policy", err, codes.Internal, "bucket", bucketName)
	}

	accessKey, err := iamClient.CreateAccessKey(ctx, &iam.CreateAccessKeyInput{UserName: &userName})
	if err != nil {
		return nil, logAndTraceError(span, "failed creating access key", err, codes.Internal, "user", userName)
	}

	credentials := assembleCredentials(ctx, accessKey, s.s3Endpoint, userName, bucketName)
	return &cosi.DriverGrantBucketAccessResponse{AccountId: userName, Credentials: credentials}, nil
}

func BuildResourceStrings(bucketName string) []string {
	return []string{
		fmt.Sprintf("arn:aws:s3:::%s/*", bucketName),
		fmt.Sprintf("arn:aws:s3:::%s", bucketName),
	}
}

func BuildPrincipalString(userName, namespace string) string {
	return fmt.Sprintf("urn:ecs:iam::%s:user/%s", namespace, userName)
}
