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
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sigs.k8s.io/container-object-storage-interface-provisioner-sidecar/pkg/consts"

	l "github.com/dell/cosi/pkg/logger"
	otelCodes "go.opentelemetry.io/otel/codes"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"

	"github.com/dell/cosi/pkg/provisioner/policy"
	"github.com/dell/goobjectscale/pkg/client/model"
)

// All errors that can be returned by DriverGrantBucketAccess.
var (
	ErrEmptyBucketAccessName            = errors.New("empty bucket access name")
	ErrInvalidAuthenticationType        = errors.New("invalid authentication type")
	ErrUnknownAuthenticationType        = errors.New("unknown authentication type")
	ErrBucketNotFound                   = errors.New("bucket not found")
	ErrFailedToCreateUser               = errors.New("failed to create user")
	ErrFailedToDecodePolicy             = errors.New("failed to decode bucket policy")
	ErrFailedToUpdatePolicy             = errors.New("failed to update bucket policy")
	ErrFailedToCreateAccessKey          = errors.New("failed to create access key")
	ErrFailedToGeneratePolicyID         = errors.New("failed to generate PolicyID UUID")
	ErrGeneratedPolicyIDIsEmpty         = errors.New("generated PolicyID was empty")
	ErrAuthenticationTypeNotImplemented = errors.New("authentication type IAM not implemented")
)

// DriverGrantBucketAccess provides access to Bucket on specific Object Storage Platform.
func (s *Server) DriverGrantBucketAccess(
	ctx context.Context,
	req *cosi.DriverGrantBucketAccessRequest,
) (*cosi.DriverGrantBucketAccessResponse, error) {
	ctx, span := otel.Tracer(GrantBucketAccessTraceName).Start(ctx, "ObjectscaleDriverGrantBucketAccess")
	defer span.End()

	// Check if bucketID is not empty.
	if err := isBucketIDEmpty(req); err != nil {
		return nil, logAndTraceError(span, ErrInvalidBucketID.Error(), err, codes.InvalidArgument)
	}

	// Check if bucket access name is not empty.
	if err := isBucketAccessNameEmpty(req); err != nil {
		return nil, logAndTraceError(span, ErrEmptyBucketAccessName.Error(), err, codes.InvalidArgument)
	}

	// Check if authentication type is not unknown.
	if err := isAuthenticationTypeNotEmpty(req); err != nil {
		return nil, logAndTraceError(span, ErrInvalidAuthenticationType.Error(), err, codes.InvalidArgument)
	}

	// Now split the flow based on the type of authentication.
	if req.AuthenticationType == cosi.AuthenticationType_IAM {
		return handleIAMAuthentication(ctx, s, req)
	}

	if req.AuthenticationType == cosi.AuthenticationType_Key {
		return handleKeyAuthentication(ctx, s, req)
	}

	return nil, logAndTraceError(span, ErrUnknownAuthenticationType.Error(), ErrUnknownAuthenticationType, codes.Internal)
}

// handleKeyAuthentication is a function providing the bucket access granting functionality,
// which uses the key type authentication method.
func handleKeyAuthentication(ctx context.Context, s *Server, req *cosi.DriverGrantBucketAccessRequest) (*cosi.DriverGrantBucketAccessResponse, error) {
	ctx, span := otel.Tracer(GrantBucketAccessTraceName).Start(ctx, "ObjectscaleHandleKeyAuthentication")
	defer span.End()

	// Get bucket name from bucketID.
	bucketName, err := GetBucketName(req.GetBucketId())
	if err != nil {
		return nil, logAndTraceError(span, ErrInvalidBucketID.Error(), err, codes.InvalidArgument)
	}

	l.Log().V(4).Info("Bucket access for bucket is being created.", "bucket", bucketName, "bucket_access", req.Name)

	// Construct common parameters for bucket requests.
	parameters := make(map[string]string)
	parameters["namespace"] = s.namespace

	l.Log().V(4).Info("Parameters of the bucket.", "parameters", parameters)

	// Check if bucket for granting access exists.
	_, err = s.mgmtClient.Buckets().Get(ctx, bucketName, parameters)
	if err != nil {
		if errors.Is(err, ErrParameterNotFound) {
			return nil, logAndTraceError(span, ErrBucketNotFound.Error(), err, codes.NotFound, "bucket", bucketName)
		}

		return nil, logAndTraceError(span, ErrFailedToCheckBucketExists.Error(), err, codes.Internal, "bucket", bucketName)
	}

	// This flow below will check for user existence; if user does not exist, it will create one. It will only fail
	// in case of an unknown error, e.g. network issues, to adhere to idempotency requirement.
	userName := BuildUsername(s.namespace, bucketName)

	// Retrieve the user.
	userGet, err := s.iamClient.GetUserWithContext(ctx, &iam.GetUserInput{UserName: &userName})
	if err != nil {
		var myAwsErr awserr.Error

		if errors.As(err, &myAwsErr) {
			// If we got a known error, but it's not "user does not exist" error, we fail.
			if myAwsErr.Code() != iam.ErrCodeNoSuchEntityException {
				span.RecordError(myAwsErr)
				return nil, logAndTraceError(span, ErrFailedToCheckUserExists.Error(), err, codes.Internal, "user", userName)
			}
		} else {
			// If we got an unknown error, we fail.
			return nil, logAndTraceError(span, ErrFailedToCheckUserExists.Error(), err, codes.Internal, "user", userName)
		}
	}

	// Check if IAM user exists.
	if userGet.User != nil {
		// Case when user exists.
		l.Log().V(0).Info("User already exists.", "user", userName)
	} else {
		// Case when user does not exist - create one.
		user, err := s.iamClient.CreateUserWithContext(ctx, &iam.CreateUserInput{
			UserName: &userName,
		})
		if err != nil {
			return nil, logAndTraceError(span, ErrFailedToCreateUser.Error(), err, codes.Internal, "user", userName)
		}

		l.Log().V(4).Info("ObjectScale IAM user was created.", "user", userName, "userID", user.User.UserId)
	}

	// Check if policy for a specific bucket exists.
	existingPolicy, err := s.mgmtClient.Buckets().GetPolicy(ctx, bucketName, parameters)
	if err != nil && !errors.Is(err, model.Error{Code: model.CodeResourceNotFound}) {
		return nil, logAndTraceError(span, ErrFailedToCheckPolicyExists.Error(), err, codes.Internal, "bucket", bucketName)
	}

	policyRequest := policy.Document{}
	if existingPolicy != "" {
		err = json.NewDecoder(strings.NewReader(existingPolicy)).Decode(&policyRequest)
		if err != nil {
			return nil, logAndTraceError(span, ErrFailedToDecodePolicy.Error(), err, codes.Internal, "bucket", bucketName, "policy", existingPolicy)
		}
	}

	// Update policy.
	awsBucketResourceARN := BuildResourceString(s.objectScaleID, s.objectStoreID, bucketName)
	awsPrincipalString := BuildPrincipalString(s.namespace, bucketName)
	policyRequest.Statement = parsePolicyStatement(
		ctx, policyRequest.Statement, awsBucketResourceARN, awsPrincipalString,
	)

	l.Log().V(4).Info("Policy request statement was parsed.", "awsBucketResourceARN", awsBucketResourceARN, "awsPrincipalString", awsPrincipalString, "statement", policyRequest.Statement)

	if policyRequest.ID == "" {
		policyID, err := generatePolicyID(ctx)
		if err != nil {
			return nil, logAndTraceError(span, ErrFailedToGeneratePolicyID.Error(), err, codes.Internal, "bucket", bucketName, "PolicyID", policyID)
		}

		l.Log().V(4).Info("policyID was generated.", "policy", policyRequest)

		span.AddEvent("policyID was generated")
	}

	if policyRequest.Version == "" {
		policyRequest.Version = bucketVersion
	}

	// Marshal the struct to JSON to confirm JSON validity.
	updateBucketPolicyJSON, err := json.Marshal(policyRequest)
	if err != nil {
		return nil, logAndTraceError(span, ErrFailedToMarshalPolicy.Error(), err, codes.Internal, "bucket", bucketName, "PolicyID", policyRequest.ID)
	}

	err = s.mgmtClient.Buckets().UpdatePolicy(ctx, bucketName, string(updateBucketPolicyJSON), parameters)
	if err != nil {
		return nil, logAndTraceError(span, ErrFailedToUpdatePolicy.Error(), err, codes.Internal, "bucket", bucketName, "policy", updateBucketPolicyJSON)
	}

	accessKey, err := s.iamClient.CreateAccessKey(&iam.CreateAccessKeyInput{UserName: &userName})
	if err != nil {
		return nil, logAndTraceError(span, ErrFailedToCreateAccessKey.Error(), err, codes.Internal, "user", userName)
	}

	// TODO: can credentials have empty values? if no, should we check any specific fields for non-emptiness?
	credentials := assembleCredentials(ctx, accessKey, s.s3Endpoint, userName, bucketName)

	return &cosi.DriverGrantBucketAccessResponse{AccountId: userName, Credentials: credentials}, nil
}

// TODO: this function will be implemented if we decide to add the IAM authentication.
// handleIAMAuthentication is a function providing the bucket access granting functionality,
// which uses the IAM type authentication method.
func handleIAMAuthentication(_ context.Context, _ *Server, _ *cosi.DriverGrantBucketAccessRequest) (*cosi.DriverGrantBucketAccessResponse, error) {
	return nil, status.Error(codes.Unimplemented, ErrAuthenticationTypeNotImplemented.Error())
}

// isBucketAccessNameEmpty checks if bucket access name is not empty.
func isBucketAccessNameEmpty(req *cosi.DriverGrantBucketAccessRequest) error {
	if req.GetName() == "" {
		return ErrEmptyBucketAccessName
	}

	return nil
}

// isAuthenticationTypeNotEmpty checks if authentication type is not unknown.
func isAuthenticationTypeNotEmpty(req *cosi.DriverGrantBucketAccessRequest) error {
	if req.GetAuthenticationType() == cosi.AuthenticationType_UnknownAuthenticationType {
		return ErrInvalidAuthenticationType
	}

	return nil
}

// parsePolicyStatement generates new bucket policy statements array with updated resource and principal.
func parsePolicyStatement(
	ctx context.Context,
	inputStatements []policy.StatementEntry,
	awsBucketResourceARN,
	awsPrincipalString string,
) []policy.StatementEntry {
	_, span := otel.Tracer(GrantBucketAccessTraceName).Start(ctx, "ObjectscaleParsePolicyStatement")
	defer span.End()

	outputStatements := []policy.StatementEntry{}

	// Omitting a nil check, as the len() is defined as at lest zero.
	if len(inputStatements) > 0 {
		outputStatements = inputStatements
	} else {
		outputStatements = append(outputStatements, policy.StatementEntry{})
	}

	for k := range outputStatements {
		statement := &outputStatements[k]
		if statement.Resource == nil {
			statement.Resource = []string{}
		}

		if !awsBucketResourceArnExists(statement, awsBucketResourceARN) {
			statement.Resource = append(statement.Resource, awsBucketResourceARN)
		}

		span.AddEvent("update resource in policy statement")

		if statement.Effect == "" {
			statement.Effect = allowEffect
		}

		if !principalExists(statement, awsPrincipalString) {
			statement.Principal.AWS = append(statement.Principal.AWS, awsPrincipalString)
		}

		span.AddEvent("update principal AWS in policy statement")

		// TODO: shouldn't action be validated with params? Maybe we only want to grant read access by default?
		// if yes, then this should be done later, when we have more info about the params (MVP is to grant all permissions)
		if !actionExists(statement) {
			statement.Action = append(statement.Action, "*")
		}

		span.AddEvent("update principal action in policy statement")

		// TODO: I don't think this is necessary after the changes to addressing "statement" variable
		outputStatements[k] = *statement
	}

	return outputStatements
}

// actionExists is a function used when parsing statements, which adds the Action field if none are found.
func actionExists(statement *policy.StatementEntry) bool {
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

// principalExists is a function used when parsing statements, which adds the Principal field if none are found.
func principalExists(statement *policy.StatementEntry, principalString string) bool {
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

// awsBucketResourceArnExists is a function used when parsing statements,
// which adds the awsBucketResourceARN field if none are found.
func awsBucketResourceArnExists(statement *policy.StatementEntry, awsBucketResourceARN string) bool {
	for _, r := range statement.Resource {
		if r == awsBucketResourceARN {
			return true
		}
	}

	return false
}

// generatePolicyID creates new policy for the bucket.
func generatePolicyID(ctx context.Context) (*uuid.UUID, error) {
	_, span := otel.Tracer(GrantBucketAccessTraceName).Start(ctx, "ObjectscaleGeneratePolicyID")
	defer span.End()

	policyID, err := uuid.NewUUID()
	if err != nil {
		return nil, ErrFailedToGeneratePolicyID
	}

	if policyID.String() == "" {
		return nil, ErrGeneratedPolicyIDIsEmpty
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
	_, span := otel.Tracer(GrantBucketAccessTraceName).Start(ctx, "ObjectscaleAssembeCredentials")
	defer span.End()

	secretsMap := make(map[string]string)
	secretsMap[consts.S3SecretAccessKeyID] = *accessKey.AccessKey.AccessKeyId
	secretsMap[consts.S3SecretAccessSecretKey] = *accessKey.AccessKey.SecretAccessKey
	secretsMap[consts.S3Endpoint] = s3Endpoint
	secretsMap["bucketName"] = bucketName

	l.Log().V(4).Info("Secret access key for user with endpoint was created.", "user", userName, "secretKeyId", *accessKey.AccessKey.AccessKeyId, "endpoint", s3Endpoint)

	span.AddEvent("secret access key for user with endpoint was created")

	credentialDetails := cosi.CredentialDetails{Secrets: secretsMap}
	credentials := make(map[string]*cosi.CredentialDetails)
	credentials[consts.S3Key] = &credentialDetails

	l.Log().V(4).Info("Access to the bucket for user successfully granted.", "bucket", bucketName, "user", userName)

	span.AddEvent("access to the bucket for user successfully granted")

	return credentials
}

// logAndTraceError is a helper function that logs an error with specified fields and records it in a span.
func logAndTraceError(span trace.Span, errMsg string, err error, code codes.Code, keysAndValues ...interface{}) error {
	l.Log().Error(err, errMsg, keysAndValues)

	span.RecordError(err)
	span.SetStatus(otelCodes.Error, errMsg)

	return status.Error(code, errMsg)
}
