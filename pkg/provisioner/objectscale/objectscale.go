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

// Package objectscale ...
// TODO: write documentation comment for objectscale package
package objectscale

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/dell/goobjectscale/pkg/client/api"
	"github.com/dell/goobjectscale/pkg/client/model"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sigs.k8s.io/container-object-storage-interface-provisioner-sidecar/pkg/consts"

	driver "github.com/dell/cosi-driver/pkg/provisioner/virtualdriver"
	objectscaleRest "github.com/dell/goobjectscale/pkg/client/rest"
	objectscaleClient "github.com/dell/goobjectscale/pkg/client/rest/client"
	iamObjectscale "github.com/dell/goobjectscale/pkg/client/rest/iam"
	log "github.com/sirupsen/logrus"
	otelCodes "go.opentelemetry.io/otel/codes"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"

	"github.com/dell/cosi-driver/pkg/config"
	"github.com/dell/cosi-driver/pkg/transport"
)

const (
	splitNumber = 2
	// bucketVersion is used when sending the bucket policy update request.
	bucketVersion = "2012-10-17"
	// allowEffect is used when updating the bucket policy, in order to grant permissions to user.
	allowEffect = "Allow"
)

// defaultTimeout is default call length before context is getting canceled.
var defaultTimeout = time.Second * 20

// Server is implementation of driver.Driver interface for ObjectScale platform.
type Server struct {
	mgmtClient         api.ClientSet
	iamClient          iamiface.IAMAPI
	x509Client         http.Client
	objClient          objectscaleClient.AuthUser
	backendID          string
	emptyBucket        bool
	namespace          string
	username           string
	password           string
	region             string
	objectScaleGateway string
	objectStoreGateway string
	objectScaleID      string
	objectStoreID      string
	s3Endpoint         string
}

var _ driver.Driver = (*Server)(nil)

// New initializes server based on the config file.
// TODO: verify if emptiness verification can be moved to a separate function.
func New(ctx context.Context, config *config.Objectscale) (*Server, error) {
	id := config.Id
	if id == "" {
		return nil, errors.New("empty driver id")
	}

	namespace := config.Namespace
	if namespace == "" {
		return nil, errors.New("empty objectstore id")
	}

	username := config.Credentials.Username
	if username == "" {
		return nil, errors.New("empty username")
	}

	password := config.Credentials.Password
	if password == "" {
		return nil, errors.New("empty password")
	}

	region := config.Region
	if region == nil {
		return nil, errors.New("region was not specified in config")
	}

	if *region == "" {
		return nil, errors.New("empty region provided")
	}

	objectscaleGateway := config.ObjectscaleGateway
	if objectscaleGateway == "" {
		return nil, errors.New("empty objectscale gateway")
	}

	objectstoreGateway := config.ObjectstoreGateway
	if objectstoreGateway == "" {
		return nil, errors.New("empty objectstore gateway")
	}

	objectscaleID := config.ObjectscaleId
	if objectscaleID == "" {
		return nil, errors.New("empty objectscaleID")
	}

	objectstoreID := config.ObjectstoreId
	if objectstoreID == "" {
		return nil, errors.New("empty objectstoreID")
	}

	protocolS3Endpoint := config.Protocols.S3.Endpoint
	if protocolS3Endpoint == "" {
		return nil, errors.New("empty protocol S3 endpoint")
	}

	if strings.Contains(id, "-") {
		id = strings.ReplaceAll(id, "-", "_")

		log.WithFields(log.Fields{
			"id":        id,
			"config.id": config.Id,
		}).Warn("id in config contains hyphens, which will be replaced with underscores")
	}

	transport, err := transport.New(config.Tls)
	if err != nil {
		return nil, fmt.Errorf("initialization of transport failed: %w", err)
	}

	objectscaleAuthUser := objectscaleClient.AuthUser{
		Gateway:  objectscaleGateway,
		Username: username,
		Password: password,
	}
	mgmtClient := objectscaleRest.NewClientSet(
		&objectscaleClient.Simple{
			Endpoint:       objectstoreGateway,
			Authenticator:  &objectscaleAuthUser,
			HTTPClient:     &http.Client{Transport: transport},
			OverrideHeader: false,
		},
	)

	x509Client := *http.DefaultClient
	objClient := objectscaleClient.AuthUser{
		Gateway:  objectscaleGateway,
		Username: username,
		Password: password,
	}

	handler := request.NamedHandler{
		Name: iamObjectscale.SDSHandlerName,
		Fn: func(r *request.Request) {
			if !objClient.IsAuthenticated() {
				ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
				defer cancel()

				err := objClient.Login(ctx, &x509Client)
				if err != nil {
					r.Error = err // no return intentional
				}
			}

			token := objClient.Token()
			r.HTTPRequest.Header.Add(iamObjectscale.SDSHeaderName, token)
		},
	}

	iamSession, err := session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				Endpoint:                      &objectstoreGateway,
				Region:                        region,
				CredentialsChainVerboseErrors: aws.Bool(true),
				HTTPClient:                    &x509Client,
			},
		},
	)

	iamSession.Handlers.Sign.RemoveByName(v4.SignRequestHandler.Name)
	swapped := iamSession.Handlers.Sign.SwapNamed(handler)

	if !swapped {
		iamSession.Handlers.Sign.PushFrontNamed(handler)
	}

	handler2 := request.NamedHandler{
		Name: iamObjectscale.AccountIDHandlerName,
		Fn: func(r *request.Request) {
			r.HTTPRequest.Header.Add(iamObjectscale.AccountIDHeaderName, namespace)
		},
	}

	swapped = iamSession.Handlers.Sign.SwapNamed(handler2)
	if !swapped {
		iamSession.Handlers.Sign.PushFrontNamed(handler2)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create new IAM session: %w", err)
	}

	iamClient := iam.New(iamSession)

	return &Server{
		mgmtClient:         mgmtClient,
		iamClient:          iamClient,
		x509Client:         x509Client,
		objClient:          objClient,
		backendID:          id,
		namespace:          namespace,
		emptyBucket:        config.EmptyBucket,
		username:           username,
		password:           password,
		region:             *region,
		objectScaleGateway: objectscaleGateway,
		objectStoreGateway: objectstoreGateway,
		objectScaleID:      objectscaleID,
		objectStoreID:      objectstoreID,
		s3Endpoint:         protocolS3Endpoint,
	}, nil
}

// ID extends COSI interface by adding ID method.
func (s *Server) ID() string {
	return s.backendID
}

// DriverDeleteBucket deletes Bucket on specific Object Storage Platform.
func (s *Server) DriverDeleteBucket(ctx context.Context,
	req *cosi.DriverDeleteBucketRequest,
) (*cosi.DriverDeleteBucketResponse, error) {
	_, span := otel.Tracer("DeleteBucketRequest").Start(ctx, "ObjectscaleDriverDeleteBucket")
	defer span.End()

	log.WithFields(log.Fields{
		"bucketID": req.BucketId,
	}).Info("bucket is being deleted")

	span.AddEvent("bucket is being deleted")

	// Check if bucketID is not empty.
	if req.GetBucketId() == "" {
		err := errors.New("empty bucketID")
		log.Error(err.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Extract bucket name from bucketID.
	bucketName := strings.SplitN(req.BucketId, "-", splitNumber)[1]

	// Delete bucket.
	err := s.mgmtClient.Buckets().Delete(ctx, bucketName, s.namespace, s.emptyBucket)

	if errors.Is(err, model.Error{Code: model.CodeResourceNotFound}) {
		log.WithFields(log.Fields{
			"bucket": bucketName,
		}).Warn("bucket does not exist")

		span.AddEvent("bucket does not exist")

		return &cosi.DriverDeleteBucketResponse{}, nil
	}

	if err != nil {
		log.WithFields(log.Fields{
			"bucket": bucketName,
			"error":  err,
		}).Error("failed to delete bucket")

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "failed to delete bucket")

		return nil, status.Error(codes.Internal, "bucket was not successfully deleted")
	}

	log.WithFields(log.Fields{
		"bucket": bucketName,
	}).Info("bucket successfully deleted")

	span.AddEvent("bucket successfully deleted")

	return &cosi.DriverDeleteBucketResponse{}, nil
}

type Principal struct {
	AWS    []string `json:"AWS"`
	Action []string `json:"Action"`
}

type UpdateBucketPolicyStatement struct {
	Resource  []string  `json:"Resource"`
	SID       string    `json:"Sid"`
	Effect    string    `json:"Effect"`
	Principal Principal `json:"Principal"`
}

type UpdateBucketPolicyRequest struct {
	PolicyID  string                        `json:"Id"`
	Version   string                        `json:"Version"`
	Statement []UpdateBucketPolicyStatement `json:"Statement"`
}

// DriverGrantBucketAccess provides access to Bucket on specific Object Storage Platform.
// TODO: how about splitting key and IAM mechanisms into different functions?
// TODO: this probably has to be refactored in order to meet the gocognit requirements (complexity < 30).
func (s *Server) DriverGrantBucketAccess( // nolint:gocognit
	ctx context.Context,
	req *cosi.DriverGrantBucketAccessRequest,
) (*cosi.DriverGrantBucketAccessResponse, error) {
	ctx, span := otel.Tracer("GrantBucketAccessRequest").Start(ctx, "ObjectscaleDriverGrantBucketAccess")
	defer span.End()

	// Check if bucketID is not empty.
	if req.GetBucketId() == "" {
		err := errors.New("empty bucketID")
		log.Error(err.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

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

	userName := fmt.Sprintf("%v-user-%v", s.namespace, bucketName)

	userGet, err := s.iamClient.GetUser(&iam.GetUserInput{UserName: &userName})
	if err != nil && err.Error() != iam.ErrCodeNoSuchEntityException {
		errMsg := errors.New("failed to check for user existence")
		log.WithFields(log.Fields{
			"user":  userName,
			"error": err,
		}).Error(errMsg.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, errMsg.Error())

		return nil, status.Error(codes.Internal, "failed to check for user existence")
	}

	if userGet != nil && *userGet.User.UserName == userName {
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
	awsBucketResourceARN := fmt.Sprintf("arn:aws:s3:%s:%s:%s/*", s.objectScaleID, s.objectStoreID, bucketName)
	awsPrincipalString := fmt.Sprintf("urn:osc:iam::%s:user/%s", s.namespace, userName)
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

		if statement.Principal.Action == nil {
			statement.Principal.Action = []string{}
		}

		for _, a := range statement.Principal.Action {
			if a == "*" {
				foundAction = true
			}
		}

		if !foundAction {
			statement.Principal.Action = append(statement.Principal.Action, "*")
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
