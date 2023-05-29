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
	"github.com/google/uuid"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	driver "github.com/dell/cosi-driver/pkg/provisioner/virtualdriver"
	objectscaleRest "github.com/dell/goobjectscale/pkg/client/rest"
	objectscaleClient "github.com/dell/goobjectscale/pkg/client/rest/client"
	iamObjectscale "github.com/dell/goobjectscale/pkg/client/rest/iam"
	log "github.com/sirupsen/logrus"
	otelCodes "go.opentelemetry.io/otel/codes"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"

	"github.com/dell/cosi-driver/pkg/config"
	"github.com/dell/cosi-driver/pkg/transport"
	"github.com/dell/goobjectscale/pkg/client/api"
	"github.com/dell/goobjectscale/pkg/client/model"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	splitNumber = 2
	// bucketVersion is used when sending the bucket policy update request
	bucketVersion = "2012-10-17"
	// allowEffect is used when updating the bucket policy, in order to grant permissions to user
	allowEffect = "Allow"
)

// Server is implementation of driver.Driver interface for ObjectScale platform.
type Server struct {
	mgmtClient         api.ClientSet
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
	s3Endpont          string
}

var _ driver.Driver = (*Server)(nil)

// TODO: verify if emptiness verification can be moved to a separate function
// New initializes server based on the config file.
func New(config *config.Objectscale) (*Server, error) {
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

	return &Server{
		mgmtClient:         mgmtClient,
		backendID:          id,
		namespace:          namespace,
		emptyBucket:        config.EmptyBucket,
		username:           username,
		password:           password,
		region:             *region,
		objectScaleGateway: objectscaleGateway,
		objectStoreGateway: objectstoreGateway,
		s3Endpont:          protocolS3Endpoint,
	}, nil
}

// ID extends COSI interface by adding ID method.
func (s *Server) ID() string {
	return s.backendID
}

// DriverCreateBucket creates Bucket on specific Object Storage Platform.
func (s *Server) DriverCreateBucket(
	ctx context.Context,
	req *cosi.DriverCreateBucketRequest,
) (*cosi.DriverCreateBucketResponse, error) {
	_, span := otel.Tracer("CreateBucketRequest").Start(ctx, "ObjectscaleDriverCreateBucket")
	defer span.End()

	log.WithFields(log.Fields{
		"bucket": req.GetName(),
	}).Info("bucket is being created")

	span.AddEvent("bucket is being created")

	// Create bucket model.
	bucket := &model.Bucket{}
	bucket.Name = req.GetName()
	bucket.Namespace = s.namespace

	// Check if bucket name is not empty.
	if bucket.Name == "" {
		err := errors.New("empty bucket name")
		log.Error(err.Error())

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Display all request parameters.
	parameters := ""
	parametersCopy := make(map[string]string)

	for key, value := range req.GetParameters() {
		parameters += key + ":" + value + ";"
		parametersCopy[key] = value
	}

	// TODO: is this good way of doing this?
	parametersCopy["namespace"] = s.namespace

	log.WithFields(log.Fields{
		"parameters": parameters,
	}).Info("parameters of the bucket")

	// Remove backendID, as this is not valid parameter for bucket creation in ObjectScale.
	delete(parametersCopy, "backendID")

	// Check if bucket with specific name and parameters already exists.
	_, err := s.mgmtClient.Buckets().Get(bucket.Name, parametersCopy)
	if err != nil && !errors.Is(err, model.Error{Code: model.CodeParameterNotFound}) {
		log.WithFields(log.Fields{
			"bucket": bucket.Name,
			"error":  err,
		}).Error("failed to check bucket existence")

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "failed to check bucket existence")

		return nil, status.Error(codes.Internal, "an unexpected error occurred")
	} else if err == nil {
		log.WithFields(log.Fields{
			"bucket": bucket.Name,
		}).Warn("bucket already exists")

		span.AddEvent("bucket already exists")

		return &cosi.DriverCreateBucketResponse{
			BucketId: strings.Join([]string{s.backendID, bucket.Name}, "-"),
		}, nil
	}

	// Create bucket.
	bucket, err = s.mgmtClient.Buckets().Create(*bucket)
	if err != nil {
		log.WithFields(log.Fields{
			"bucket": bucket.Name,
			"error":  err,
		}).Error("failed to create bucket")

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "failed to create bucket")

		return nil, status.Error(codes.Internal, "bucket was not successfully created")
	}

	log.WithFields(log.Fields{
		"bucket": bucket.Name,
	}).Info("bucket successfully created")

	span.AddEvent("bucket successfully created")

	// Return response.
	return &cosi.DriverCreateBucketResponse{
		BucketId: strings.Join([]string{s.backendID, bucket.Name}, "-"),
	}, nil
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
	err := s.mgmtClient.Buckets().Delete(bucketName, s.namespace, s.emptyBucket)

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

//{
//   "Id" : "S3PolicyId1",
//   "Version" : "2012-10-17",
//   "Statement" : [ {
//   "Resource" : [ "arn:aws:s3:osci5b022e718aa7e0ff:osti202e682782ebcbfd:lynxbucket/*" ],
//   "Sid" : "GetObject_permission",
//   "Effect" : "Allow",
//   "Principal" : {
//      "AWS" : [ "urn:osc:iam::osai07c2ae318ae9d6f2:user/iam_user20230523061025118" ]
//    },
//    "Action" : [ "s3:GetObjectVersion" ]
//} ]
// }

type principal struct {
	AWS    []string `json:"AWS"`
	Action []string `json:"Action"`
}

type updateBucketPolicyStatement struct {
	Resource  []string  `json:"Resource"`
	SID       string    `json:"Sid"`
	Effect    string    `json:"Effect"`
	Principal principal `json:"Principal"`
}

type updateBucketPolicyRequest struct {
	PolicyID  string                        `json:"Id"`
	Version   string                        `json:"Version"`
	Statement []updateBucketPolicyStatement `json:"Statement"`
}

// DriverGrantBucketAccess provides access to Bucket on specific Object Storage Platform.
// TODO: how about splitting key and IAM mechanisms into different functions?
func (s *Server) DriverGrantBucketAccess(
	ctx context.Context,
	req *cosi.DriverGrantBucketAccessRequest,
) (*cosi.DriverGrantBucketAccessResponse, error) {
	// TODO: think about more spans' info
	_, span := otel.Tracer("GrantBucketAccessRequest").Start(ctx, "ObjectscaleDriverGrantBucketAccess")
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

	//TODO: after adding IAM the flow here could be if auth type key -> run key, if auth type iam -> run IAM, else error
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
	// Display all request parameters.

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
	_, err := s.mgmtClient.Buckets().Get(bucketName, parametersCopy)
	if err != nil && !errors.Is(err, model.Error{Code: model.CodeResourceNotFound}) {
		log.WithFields(log.Fields{
			"bucket": bucketName,
		}).Error("failed to check bucket existence")

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "failed to check bucket existence")

		return nil, status.Error(codes.Internal, "an unexpected error occurred")
	} else if err != nil {
		log.WithFields(log.Fields{
			"bucket": bucketName,
		}).Error("bucket not found")

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "bucket not found")

		return nil, status.Error(codes.NotFound, "bucket not found")
	}
	// Create user.
	x509Client := *http.DefaultClient
	objClient := objectscaleClient.AuthUser{
		Gateway:  s.objectScaleGateway,
		Username: s.username,
		Password: s.password,
	}
	iamSession, err := session.NewSession(&aws.Config{
		Endpoint:                      &s.objectStoreGateway,
		Region:                        &s.region,
		CredentialsChainVerboseErrors: aws.Bool(true),
		HTTPClient:                    &x509Client,
	},
	)
	if err != nil {
		log.WithFields(log.Fields{
			"endpoint": s.objectStoreGateway,
			"region":   s.region,
		}).Error("cannot create session")

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "cannot create session")

		return nil, status.Error(codes.Internal, "cannot create session")
	}

	iamClient := iam.New(iamSession)
	iamObjectscale.InjectTokenToIAMClient(iamClient, &objClient, x509Client)
	// TODO: error handling
	iamObjectscale.InjectAccountIDToIAMClient(iamClient, s.namespace)
	// TODO: error handling
	userName := fmt.Sprintf("%v-user-%v", s.namespace, bucketName)
	user, err := iamClient.CreateUser(&iam.CreateUserInput{
		UserName: &userName,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"user": userName,
		}).Error("cannot create user")

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "cannot create user")

		return nil, status.Error(codes.Internal, fmt.Sprintf("cannot create user %s", userName))
	}

	log.WithFields(log.Fields{
		"user":   userName,
		"userId": user.User.UserId,
	}).Info("ObjectScale IAM user was created")

	// Check if policy for specific bucket exists.
	_, err = s.mgmtClient.Buckets().GetPolicy(bucketName, parametersCopy)
	if err != nil && !errors.Is(err, model.Error{Code: model.CodeResourceNotFound}) {
		log.WithFields(log.Fields{
			"bucket": bucketName,
		}).Error("failed to check bucket policy existence")

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "failed to check bucket policy existence")

		return nil, status.Error(codes.Internal, "an unexpected error occurred")
		// Even if we get no error, the policy can be empty, e.g. we get a 200 OK response and empty policy
	} else if err == nil {
		log.WithFields(log.Fields{
			"bucket": bucketName,
		}).Info("bucket policy already exists")
	}

	policyID, err := uuid.NewUUID()
	if err != nil {
		log.WithFields(log.Fields{
			"bucket": bucketName,
		}).Error("failed to generate PolicyID UUID")

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "failed to generate PolicyID UUID")

		return nil, status.Error(codes.Internal, "bucket policy was not updated")
	}

	if policyID.String() == "" {
		errMsg := errors.New("generated PolicyID was empty")
		log.WithFields(log.Fields{
			"bucket":   bucketName,
			"PolicyID": policyID,
		}).Error(errMsg.Error())

		span.RecordError(errMsg)
		span.SetStatus(otelCodes.Error, errMsg.Error())

		return nil, status.Error(codes.Internal, errMsg.Error())
	}

	awsBucketResourceARN := fmt.Sprintf("arn:aws:s3:%s:%s:%s/*", s.objectScaleID, s.objectStoreID, bucketName)
	awsPrincipalString := fmt.Sprintf("urn:osc:iam::%s:user/%s", s.namespace, userName)

	bucketPolicy := updateBucketPolicyRequest{
		PolicyID: policyID.String(),
		Version:  bucketVersion,
		Statement: []updateBucketPolicyStatement{
			{
				Resource: []string{awsBucketResourceARN},
				Effect:   allowEffect,
				Principal: principal{
					AWS:    []string{awsPrincipalString},
					Action: []string{"*"},
				},
			},
		},
	}

	// Marshal the struct to JSON to confirm JSON validity
	updateBucketPolicyJson, err := json.Marshal(bucketPolicy)
	if err != nil {
		errMsg := errors.New("could not marshal updateBucketPolicyRequest into JSON")
		log.WithFields(log.Fields{
			"bucket":   bucketName,
			"PolicyID": policyID,
		}).Error(errMsg)

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, errMsg.Error())

		return nil, status.Error(codes.Internal, fmt.Sprintf("%s: %e", errMsg.Error(), err))
	}

	err = s.mgmtClient.Buckets().UpdatePolicy(bucketName, string(updateBucketPolicyJson), parametersCopy)
	if err != nil {
		log.WithFields(log.Fields{
			"bucket": bucketName,
			"policy": updateBucketPolicyJson,
		}).Error("failed to update policy")

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "failed to update policy")

		return nil, status.Error(codes.Internal, fmt.Sprintf("bucket policy was not updated: %e", err))
	}

	// Create access key.
	requestModel := &model.ObjectUserSecretKeyCreateReq{
		SecretKey: "",          // ?
		Namespace: s.namespace, // TODO: variables regarding the namespace should be renamed to smth like AccountID
	}

	secret, err := s.mgmtClient.ObjectUser().CreateSecret(userName, *requestModel, parametersCopy)
	if err != nil {
		log.WithFields(log.Fields{
			"user":   userName,
			"secret": requestModel.SecretKey,
		}).Error("failed to create secret key")

		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "failed to create secret key")

		return nil, status.Error(codes.Internal, "secret key was not successfully created")
	}

	// Assemble credential details and add to credentialRepo
	secretsMap := make(map[string]string)
	secretsMap["accessKeyID"] = userName
	secretsMap["accessSecretKey"] = secret.SecretKey
	secretsMap["endpoint"] = s.s3Endpont

	log.WithFields(log.Fields{
		"user":      userName,
		"secretKey": secret.SecretKey,
		"endpoint":  s.s3Endpont,
	}).Info("secret access key for user with endpoint was created")

	credentialDetails := cosi.CredentialDetails{Secrets: secretsMap}
	credentials := make(map[string]*cosi.CredentialDetails)
	credentials["s3"] = &credentialDetails // ?

	log.WithFields(log.Fields{
		"bucket": bucketName,
		"user":   userName,
	}).Info("access to the bucket for user successfully granted")

	return &cosi.DriverGrantBucketAccessResponse{AccountId: userName, Credentials: credentials}, nil
}

// DriverRevokeBucketAccess revokes access from Bucket on specific Object Storage Platform.
func (s *Server) DriverRevokeBucketAccess(ctx context.Context,
	req *cosi.DriverRevokeBucketAccessRequest,
) (*cosi.DriverRevokeBucketAccessResponse, error) {
	_, span := otel.Tracer("RevokeBucketAccessRequest").Start(ctx, "ObjectscaleDriverRevokeBucketAccess")
	defer span.End()

	err := errors.New("not implemented")
	span.RecordError(err)
	span.SetStatus(otelCodes.Error, err.Error())

	return nil, status.Error(codes.Unimplemented, err.Error())
}
