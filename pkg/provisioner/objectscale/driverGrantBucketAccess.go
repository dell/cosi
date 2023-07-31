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
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/dell/cosi/pkg/provisioner/policy"
	"github.com/dell/goobjectscale/pkg/client/model"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sigs.k8s.io/container-object-storage-interface-provisioner-sidecar/pkg/consts"

	log "github.com/sirupsen/logrus"
	otelCodes "go.opentelemetry.io/otel/codes"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

// All errors that can be returned by DriverGrantBucketAccess.
var (
	ErrInvalidBucketID           = errors.New("invalid bucketID")
	ErrEmptyBucketAccessName     = errors.New("empty bucket access name")
	ErrInvalidAuthenticationType = errors.New("invalid authentication type")
	ErrUnknownAuthenticationType = errors.New("unknown authentication type")
	ErrBucketNotFound            = errors.New("bucket not found")
	ErrFailedToCheckBucketExist  = errors.New("failed to check bucket existence")
	ErrFailedToCheckUserExist    = errors.New("failed to check for user existence")
	ErrFailedToCreateUser        = errors.New("failed to create user")
	ErrFailedToCheckPolicyExist  = errors.New("failed to check bucket policy existence")
	ErrFailedToDecodePolicy      = errors.New("failed to decode bucket policy")
	ErrFailedToUpdatePolicy      = errors.New("failed to update bucket policy")
	ErrFailedToCreateAccessKey   = errors.New("failed to create access key")
	ErrFailedToMarshalPolicy     = errors.New("failed to marshal updateBucketPolicyRequest into JSON")
	ErrFailedToGeneratePolicyID  = errors.New("failed to generate PolicyID UUID")
	ErrGeneratedPolicyIDIsEmpty  = errors.New("generated PolicyID was empty")
)

// Check if bucketID is not empty.
func isBucketIDEmpty(req *cosi.DriverGrantBucketAccessRequest) error {
	if req.GetBucketId() == "" {
		return ErrInvalidBucketID
	}

	return nil
}

// Check if bucket access name is not empty.
func isBucketAccessNameEmpty(req *cosi.DriverGrantBucketAccessRequest) error {
	if req.GetName() == "" {
		return ErrEmptyBucketAccessName
	}

	return nil
}

// Put error message into span and logs.
func putErrorIntoSpanAndLogs(span trace.Span, err error) {
	log.Error(err.Error())
	span.RecordError(err)
	span.SetStatus(otelCodes.Error, err.Error())
}

func putErrorIntoSpanAndLogsWithFields(span trace.Span, err error) {
	log.Error(err.Error())
	span.RecordError(err)
	span.SetStatus(otelCodes.Error, err.Error())
}

// Contruct common parameters for bucket requests.
func constructParameters(req *cosi.DriverGrantBucketAccessRequest, s *Server) map[string]string {
	parameters := ""
	parametersCopy := make(map[string]string)

	for key, value := range req.GetParameters() {
		parameters += key + ":" + value + ";"
		parametersCopy[key] = value
	}

	parametersCopy["namespace"] = s.namespace

	return parametersCopy
}

// Check if authentication type is not unknown.
func isAuthenticationTypeNotEmpty(req *cosi.DriverGrantBucketAccessRequest) error {
	if req.GetAuthenticationType() == cosi.AuthenticationType_UnknownAuthenticationType {
		return ErrInvalidAuthenticationType
	}

	return nil
}

func handleIAMAuthentication(_ context.Context, _ *Server, _ *cosi.DriverGrantBucketAccessRequest) (*cosi.DriverGrantBucketAccessResponse, error) {
	return nil, status.Error(codes.Unimplemented, "authentication type IAM not implemented")
}

func handleKeyAuthentication(ctx context.Context, s *Server, req *cosi.DriverGrantBucketAccessRequest) (*cosi.DriverGrantBucketAccessResponse, error) {
	ctx, span := otel.Tracer("GrantBucketAccessRequest").Start(ctx, "ObjectscaleHandleKeyAuthentication")
	defer span.End()

	// Get bucket name from bucketID.
	bucketName, err := GetBucketName(req.GetBucketId())
	if err != nil {
		putErrorIntoSpanAndLogs(span, err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	log.WithFields(log.Fields{
		"bucket":        bucketName,
		"bucket_access": req.Name,
	}).Info("bucket access for bucket is being created")

	// Construct common parameters for bucket requests.
	parameters := constructParameters(req, s)

	log.WithFields(log.Fields{
		"parameters": parameters,
	}).Info("parameters of the bucket")

	// Check if bucket for granting access exists.
	_, err = s.mgmtClient.Buckets().Get(ctx, bucketName, parameters)
	if err != nil {
		fields := log.Fields{
			"bucket": bucketName,
			"error":  err,
		}
		if errors.Is(err, ErrParameterNotFound) {
			log.WithFields(fields).Error(ErrBucketNotFound)
			span.RecordError(err)
			span.SetStatus(otelCodes.Error, ErrBucketNotFound.Error())

			return nil, status.Error(codes.NotFound, ErrBucketNotFound.Error())
		}

		log.WithFields(fields).Error(ErrFailedToCheckBucketExist)
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, ErrFailedToCheckBucketExist.Error())

		return nil, status.Error(codes.Internal, ErrFailedToCheckBucketExist.Error())
	}

	userName := BuildUsername(s.namespace, bucketName)

	userGet, err := s.iamClient.GetUserWithContext(ctx, &iam.GetUserInput{UserName: &userName})
	if err != nil {
		fields := log.Fields{
			"user":  userName,
			"error": err,
		}

		var myAwsErr awserr.Error

		if errors.As(err, &myAwsErr) {
			if myAwsErr.Code() != iam.ErrCodeNoSuchEntityException {
				log.WithFields(fields).Error(ErrFailedToCheckUserExist)
				span.RecordError(myAwsErr)
				span.RecordError(ErrFailedToCheckUserExist)
				span.SetStatus(otelCodes.Error, ErrFailedToCheckUserExist.Error())

				return nil, status.Error(codes.Internal, ErrFailedToCheckUserExist.Error())
			}
		} else {
			log.WithFields(fields).Error(ErrFailedToCheckUserExist)
			span.RecordError(err)
			span.SetStatus(otelCodes.Error, ErrFailedToCheckUserExist.Error())

			return nil, status.Error(codes.Internal, ErrFailedToCheckUserExist.Error())
		}
	}

	// Check if IAM user exists.
	if userGet.User != nil {
		// Case when user exists.
		log.WithFields(log.Fields{
			"user": userName,
		}).Warn("user already exists")
	} else {
		// Case when user does not exist.
		// TODO: tutaj skończyłem !!
		user, err := s.iamClient.CreateUserWithContext(ctx, &iam.CreateUserInput{
			UserName: &userName,
		})
		// add idempotency case (user exists)
		if err != nil {
			log.WithFields(log.Fields{
				"user":  userName,
				"error": err,
			}).Error(ErrFailedToCreateUser)

			span.RecordError(err)
			span.SetStatus(otelCodes.Error, ErrFailedToCreateUser.Error())

			return nil, status.Error(codes.Internal, ErrFailedToCreateUser.Error())
		}

		log.WithFields(log.Fields{
			"user":   userName,
			"userId": user.User.UserId,
		}).Info("ObjectScale IAM user was created")
	}

	// Check if policy for specific bucket exists.

	existingPolicy, err := s.mgmtClient.Buckets().GetPolicy(ctx, bucketName, parameters)
	if err != nil && !errors.Is(err, model.Error{Code: model.CodeResourceNotFound}) {
		log.WithFields(log.Fields{
			"bucket": bucketName,
			"error":  err,
		}).Error(ErrFailedToCheckPolicyExist)

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, ErrFailedToCheckPolicyExist.Error())

		return nil, status.Error(codes.Internal, ErrFailedToCheckPolicyExist.Error())
	}

	policyRequest := policy.Document{}
	if existingPolicy != "" {
		err = json.NewDecoder(strings.NewReader(existingPolicy)).Decode(&policyRequest)
		if err != nil {
			log.WithFields(log.Fields{
				"bucket": bucketName,
				"policy": existingPolicy,
				"error":  err,
			}).Error(ErrFailedToDecodePolicy)

			span.RecordError(err)
			span.SetStatus(otelCodes.Error, ErrFailedToDecodePolicy.Error())

			return nil, status.Error(codes.Internal, ErrFailedToDecodePolicy.Error())
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

	if policyRequest.ID == "" {
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
			"PolicyID": policyRequest.ID,
			"error":    err,
		}).Error(errMsg.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, errMsg.Error())

		return nil, status.Error(codes.Internal, errMsg.Error())
	}

	err = s.mgmtClient.Buckets().UpdatePolicy(ctx, bucketName, string(updateBucketPolicyJSON), parameters)
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

// DriverGrantBucketAccess provides access to Bucket on specific Object Storage Platform.
func (s *Server) DriverGrantBucketAccess(
	ctx context.Context,
	req *cosi.DriverGrantBucketAccessRequest,
) (*cosi.DriverGrantBucketAccessResponse, error) {
	ctx, span := otel.Tracer("GrantBucketAccessRequest").Start(ctx, "ObjectscaleDriverGrantBucketAccess")
	defer span.End()

	// Check if bucketID is not empty.
	if err := isBucketIDEmpty(req); err != nil {
		putErrorIntoSpanAndLogs(span, err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Check if bucket access name is not empty.
	if err := isBucketAccessNameEmpty(req); err != nil {
		putErrorIntoSpanAndLogs(span, err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Check if authentication type is not unknown.
	if err := isAuthenticationTypeNotEmpty(req); err != nil {
		putErrorIntoSpanAndLogs(span, err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if req.AuthenticationType == cosi.AuthenticationType_IAM {
		return handleIAMAuthentication(ctx, s, req)
	}

	if req.AuthenticationType == cosi.AuthenticationType_Key {
		return handleKeyAuthentication(ctx, s, req)
	}

	putErrorIntoSpanAndLogs(span, ErrUnknownAuthenticationType)

	return nil, status.Error(codes.Internal, ErrUnknownAuthenticationType.Error())
}

// parsePolicyStatement generates new bucket policy statements array with updated resource and principal.
// TODO: this probably has to be refactored in order to meet the gocognit requirements (complexity < 30).
func parsePolicyStatement(
	ctx context.Context,
	inputStatements []policy.StatementEntry,
	awsBucketResourceARN,
	awsPrincipalString string,
) []policy.StatementEntry {
	_, span := otel.Tracer("GrantBucketAccessRequest").Start(ctx, "ObjectscaleParsePolicyStatement")
	defer span.End()

	outputStatements := []policy.StatementEntry{}

	// Omitting a nil check, as the len() is defined as at lest zero.
	if len(inputStatements) > 0 {
		outputStatements = inputStatements
	} else {
		outputStatements = append(outputStatements, policy.StatementEntry{})
	}

	for k, statement := range outputStatements {

		if statement.Resource == nil {
			statement.Resource = []string{}
		}

		if !awsBucketResourceArnExists(&statement, awsBucketResourceARN) {
			statement.Resource = append(statement.Resource, awsBucketResourceARN)
		}

		span.AddEvent("update resource in policy statement")

		if statement.Effect == "" {
			statement.Effect = allowEffect
		}

		if !principalExists(&statement, awsPrincipalString) {
			statement.Principal.AWS = append(statement.Principal.AWS, awsPrincipalString)
		}

		span.AddEvent("update principal AWS in policy statement")

		// TODO: shouldn't action be validated with params? Maybe we only want to grant read access by default?
		// if yes, then this should be done later, when we have more info about the params (MVP is to grant all permissions)
		if !actionExists(&statement) {
			statement.Action = append(statement.Action, "*")
		}

		span.AddEvent("update principal action in policy statement")

		outputStatements[k] = statement
	}

	return outputStatements
}

func actionExists(statement *UpdateBucketPolicyStatement) bool {
	if statement.Action == nil {
		statement.Action = []string{}
	}

	for _, a := range statement.Action {
		if a == "*" {
			return true
		}
	}

	return false
}

func principalExists(statement *UpdateBucketPolicyStatement, principalString string) bool {
	if statement.Principal.AWS == nil {
		statement.Principal.AWS = []string{}
	}

	for _, p := range statement.Principal.AWS {
		if p == principalString {
			return true
		}
	}

	return false
}

func awsBucketResourceArnExists(statement *UpdateBucketPolicyStatement, awsBucketResourceARN string) bool {
	for _, r := range statement.Resource {
		if r == awsBucketResourceARN {
			return true
		}
	}

	return false
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
