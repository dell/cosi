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

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/dell/goobjectscale/pkg/client/model"
	"github.com/google/uuid"
	"go.opencensus.io/trace"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sigs.k8s.io/container-object-storage-interface-provisioner-sidecar/pkg/consts"

	log "github.com/sirupsen/logrus"
	otelCodes "go.opentelemetry.io/otel/codes"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

// All errors that can be returned by DriverGrantBucketAccess.
var (
	ErrInvalidBucketID       = errors.New("invalid bucketID")
	ErrEmptyBucketAccessName = errors.New("empty bucket access name")
)

// Check if bucketID is not empty.
func isBucketIDEmpty(ctx context.Context, req *cosi.DriverGrantBucketAccessRequest) error {
	if req.GetBucketId() == "" {
		return ErrInvalidBucketID
	}

	return nil
}

// Check if bucket access name is not empty.
func isBucketAccessNameEmpty(ctx context.Context, req *cosi.DriverGrantBucketAccessRequest) error {
	if req.GetName() == "" {
		return ErrEmptyBucketAccessName
	}

	return nil
}

// Put error message into span and logs.
func putErrorIntoSpanAndLogs(ctx context.Context, span trace.Span, err error) {
	log.Error(err.Error())
	span.RecordError(err)
	span.SetStatus(otelCodes.Error, err.Error())
}

// DriverGrantBucketAccess provides access to Bucket on specific Object Storage Platform.
func (s *Server) DriverGrantBucketAccess(
	ctx context.Context,
	req *cosi.DriverGrantBucketAccessRequest,
) (*cosi.DriverGrantBucketAccessResponse, error) {
	ctx, span := otel.Tracer("GrantBucketAccessRequest").Start(ctx, "ObjectscaleDriverGrantBucketAccess")
	defer span.End()

	// Check if bucketID is not empty.
	if err := isBucketIDEmpty(ctx, req); err != nil {
		putErrorIntoSpanAndLogs(ctx, span, err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Check if bucket access name is not empty.

	// Check if bucket access name is not empty.
	if req.GetName() == "" {
		err := errors.New("empty bucket access name")
		log.Error(err.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// TODO: after adding IAM the flow here could be if auth type key -> run key, if auth type iam -> run IAM, else error
	// Check authentication type.
	if req.GetAuthenticationType() == cosi.AuthenticationType_UnknownAuthenticationType {
		err := errors.New("invalid authentication type")
		log.Error(err.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if req.AuthenticationType == cosi.AuthenticationType_IAM {
		return nil, status.Error(codes.Unimplemented, "authentication type IAM not implemented")
	}

	// TODO: this should probably be moved to a separate function
	// Extract bucket name from bucketID.
	bucketName := strings.SplitN(req.BucketId, "-", splitNumber)[1]

	log.WithFields(log.Fields{
		"bucket":        bucketName,
		"bucket_access": req.Name,
	}).Info("bucket access for bucket is being created")

	parameters := ""
	parametersCopy := make(map[string]string)

	for key, value := range req.GetParameters() {
		parameters += key + ":" + value + ";"
		parametersCopy[key] = value
	}

	parametersCopy["namespace"] = s.namespace

	log.WithFields(log.Fields{
		"parameters": parameters,
	}).Info("parameters of the bucket")

	// Check if bucket for granting access exists.
	_, err := s.mgmtClient.Buckets().Get(ctx, bucketName, parametersCopy)
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

	userName := BuildUsername(s.namespace, bucketName)

	userGet, err := s.iamClient.GetUserWithContext(ctx, &iam.GetUserInput{UserName: &userName})
	if err != nil {
		var myAwsErr awserr.Error
		if errors.As(err, &myAwsErr) {
			if myAwsErr.Code() != iam.ErrCodeNoSuchEntityException {
				errMsg := errors.New("failed to check for user existence")
				log.WithFields(log.Fields{
					"user":  userName,
					"error": myAwsErr,
				}).Error(errMsg.Error())

				span.RecordError(myAwsErr)
				span.SetStatus(otelCodes.Error, errMsg.Error())

				return nil, status.Error(codes.Internal, errMsg.Error())
			}
		} else {
			errMsg := errors.New("failed to check for user existence")
			log.WithFields(log.Fields{
				"user":  userName,
				"error": err,
			}).Error(errMsg.Error())

			span.RecordError(err)
			span.SetStatus(otelCodes.Error, errMsg.Error())

			return nil, status.Error(codes.Internal, errMsg.Error())
		}
	}

	// The userGet is being set to &iam.GetUserOutput{} in
	// GetUser function so I don't think wrapping this in additional if is necessary.
	// Check if IAM user exists.
	if userGet.User != nil {
		log.WithFields(log.Fields{
			"user": userName,
		}).Warn("user already exists")
	} else {
		user, err := s.iamClient.CreateUserWithContext(ctx, &iam.CreateUserInput{
			UserName: &userName,
		})
		// add idempotency case (user exists)
		if err != nil {
			log.WithFields(log.Fields{
				"user":  userName,
				"error": err,
			}).Error("cannot create user")

			span.RecordError(err)
			span.SetStatus(otelCodes.Error, "cannot create user")

			return nil, status.Error(codes.Internal, fmt.Sprintf("cannot create user %s", userName))
		}

		log.WithFields(log.Fields{
			"user":   userName,
			"userId": user.User.UserId,
		}).Info("ObjectScale IAM user was created")
	}

	// Check if policy for specific bucket exists.
	policy, err := s.mgmtClient.Buckets().GetPolicy(ctx, bucketName, parametersCopy)
	if err != nil && !errors.Is(err, model.Error{Code: model.CodeResourceNotFound}) {
		errMsg := errors.New("failed to check bucket policy existence")
		log.WithFields(log.Fields{
			"bucket": bucketName,
			"error":  err,
		}).Error(errMsg.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, errMsg.Error())

		return nil, status.Error(codes.Internal, errMsg.Error())
	} else if err == nil {
		// TODO: this block is not necessary
		// Even if we get no error, the policy can be empty, e.g. we get a 200 OK response and empty policy
		log.WithFields(log.Fields{
			"bucket": bucketName,
		}).Info("bucket policy already exists")
	}

	policyRequest := UpdateBucketPolicyRequest{}
	if policy != "" {
		err = json.NewDecoder(strings.NewReader(policy)).Decode(&policyRequest)
		if err != nil {
			errMsg := errors.New("failed to decode existing bucket policy")
			log.WithFields(log.Fields{
				"bucket": bucketName,
				"policy": policy,
				"error":  err,
			}).Error(errMsg.Error())

			span.RecordError(err)
			span.SetStatus(otelCodes.Error, errMsg.Error())

			return nil, status.Error(codes.Internal, errMsg.Error())
		}
	}

	// Update policy.
	awsBucketResourceARN := BuildResourceString(s.objectScaleID, s.objectStoreID, bucketName)
	awsPrincipalString := BuildPrincipalString(s.namespace, bucketName)
	policyRequest.Statement = parsePolicyStatement(
		ctx, policyRequest.Statement, awsBucketResourceARN, awsPrincipalString,
	)

	log.WithFields(log.Fields{
		"awsBucketResourceARN": awsBucketResourceARN,
		"awsPrincipalString":   awsPrincipalString,
		"statement":            policyRequest.Statement,
	}).Info("policy request statement was parsed")

	if policyRequest.PolicyID == "" {
		policyID, err := generatePolicyID(ctx, bucketName)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		log.WithFields(log.Fields{
			"policy": policyRequest,
		}).Infof("policyID %v was generated", policyID)
		span.AddEvent("policyID was generated")
	}

	if policyRequest.Version == "" {
		policyRequest.Version = bucketVersion
	}

	// Marshal the struct to JSON to confirm JSON validity
	updateBucketPolicyJSON, err := json.Marshal(policyRequest)
	if err != nil {
		errMsg := errors.New("failed to marshal updateBucketPolicyRequest into JSON")
		log.WithFields(log.Fields{
			"bucket":   bucketName,
			"PolicyID": policyRequest.PolicyID,
			"error":    err,
		}).Error(errMsg.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, errMsg.Error())

		return nil, status.Error(codes.Internal, errMsg.Error())
	}

	err = s.mgmtClient.Buckets().UpdatePolicy(ctx, bucketName, string(updateBucketPolicyJSON), parametersCopy)
	if err != nil {
		errMsg := errors.New("failed to update bucket policy")
		log.WithFields(log.Fields{
			"bucket": bucketName,
			"policy": updateBucketPolicyJSON,
			"error":  err,
		}).Error(errMsg.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, errMsg.Error())

		return nil, status.Error(codes.Internal, errMsg.Error())
	}

	accessKey, err := s.iamClient.CreateAccessKey(&iam.CreateAccessKeyInput{UserName: &userName})
	if err != nil {
		errMsg := errors.New("failed to create access key")
		log.WithFields(log.Fields{
			"user":  userName,
			"error": err,
		}).Error(errMsg.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, errMsg.Error())

		return nil, status.Error(codes.Internal, errMsg.Error())
	}

	// TODO: can credentials have empty values? if no, should we check any specific fields for non-emptiness?
	credentials := assembleCredentials(ctx, accessKey, s.s3Endpoint, userName, bucketName)

	return &cosi.DriverGrantBucketAccessResponse{AccountId: userName, Credentials: credentials}, nil
}

// parsePolicyStatement generates new bucket policy statements array with updated resource and principal.
// TODO: this probably has to be refactored in order to meet the gocognit requirements (complexity < 30).
func parsePolicyStatement( // nolint:gocognit
	ctx context.Context,
	inputStatements []UpdateBucketPolicyStatement,
	awsBucketResourceARN,
	awsPrincipalString string,
) []UpdateBucketPolicyStatement {
	_, span := otel.Tracer("GrantBucketAccessRequest").Start(ctx, "ObjectscaleParsePolicyStatement")
	defer span.End()

	outputStatements := []UpdateBucketPolicyStatement{}

	// Omitting a nil check, as the len() is defined as at lest zero.
	if len(inputStatements) > 0 {
		outputStatements = inputStatements
	} else {
		outputStatements = append(outputStatements, UpdateBucketPolicyStatement{})
	}

	for k, statement := range outputStatements {
		foundResource := false

		if statement.Resource == nil {
			statement.Resource = []string{}
		}

		for _, r := range statement.Resource {
			if r == awsBucketResourceARN {
				foundResource = true
			}
		}

		if !foundResource {
			statement.Resource = append(statement.Resource, awsBucketResourceARN)
		}

		span.AddEvent("update resource in policy statement")

		if statement.Effect == "" {
			statement.Effect = allowEffect
		}

		foundPrincipal := false

		if statement.Principal.AWS == nil {
			statement.Principal.AWS = []string{}
		}

		for _, p := range statement.Principal.AWS {
			if p == awsPrincipalString {
				foundPrincipal = true
			}
		}

		if !foundPrincipal {
			statement.Principal.AWS = append(statement.Principal.AWS, awsPrincipalString)
		}

		span.AddEvent("update principal AWS in policy statement")

		// TODO: shouldn't action be validated with params? Maybe we only want to grant read access by default?
		// if yes, then this should be done later, when we have more info about the params (MVP is to grant all permissions)
		foundAction := false

		if statement.Action == nil {
			statement.Action = []string{}
		}

		for _, a := range statement.Action {
			if a == "*" {
				foundAction = true
			}
		}

		if !foundAction {
			statement.Action = append(statement.Action, "*")
		}

		span.AddEvent("update principal action in policy statement")

		outputStatements[k] = statement
	}

	return outputStatements
}

// generatePolicyID creates new policy for the bucket.
func generatePolicyID(ctx context.Context, bucketName string) (*uuid.UUID, error) {
	_, span := otel.Tracer("GrantBucketAccessRequest").Start(ctx, "ObjectscaleGeneratePolicyID")
	defer span.End()

	policyID, err := uuid.NewUUID()
	if err != nil {
		errMsg := errors.New("failed to generate PolicyID UUID")
		log.WithFields(log.Fields{
			"bucket": bucketName,
			"error":  err,
		}).Error(errMsg.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, errMsg.Error())

		return nil, errMsg
	}

	if policyID.String() == "" {
		errMsg := errors.New("generated PolicyID was empty")
		log.WithFields(log.Fields{
			"bucket":   bucketName,
			"PolicyID": policyID,
		}).Error(errMsg.Error())

		span.RecordError(errMsg)
		span.SetStatus(otelCodes.Error, errMsg.Error())

		return nil, errMsg
	}

	return &policyID, nil
}

// assembleCredentials assembles credentials details and adds them to the credentialRepo.
func assembleCredentials(
	ctx context.Context,
	accessKey *iam.CreateAccessKeyOutput,
	s3Endpoint,
	userName,
	bucketName string,
) map[string]*cosi.CredentialDetails {
	_, span := otel.Tracer("GrantBucketAccessRequest").Start(ctx, "ObjectscaleAssembeCredentials")
	defer span.End()

	secretsMap := make(map[string]string)
	secretsMap[consts.S3SecretAccessKeyID] = *accessKey.AccessKey.AccessKeyId
	secretsMap[consts.S3SecretAccessSecretKey] = *accessKey.AccessKey.SecretAccessKey
	secretsMap[consts.S3Endpoint] = s3Endpoint
	secretsMap["bucketName"] = bucketName

	log.WithFields(log.Fields{
		"user":        userName,
		"secretKeyId": *accessKey.AccessKey.AccessKeyId,
		"endpoint":    s3Endpoint,
	}).Info("secret access key for user with endpoint was created")

	span.AddEvent("secret access key for user with endpoint was created")

	credentialDetails := cosi.CredentialDetails{Secrets: secretsMap}
	credentials := make(map[string]*cosi.CredentialDetails)
	credentials[consts.S3Key] = &credentialDetails

	log.WithFields(log.Fields{
		"bucket": bucketName,
		"user":   userName,
	}).Info("access to the bucket for user successfully granted")

	span.AddEvent("access to the bucket for user successfully granted")

	return credentials
}
