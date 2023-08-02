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
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/bombsimon/logrusr/v4"

	internalLogger "github.com/dell/cosi/pkg/internal/logger"
	driver "github.com/dell/cosi/pkg/provisioner/virtualdriver"
	objectscaleRest "github.com/dell/goobjectscale/pkg/client/rest"
	objectscaleClient "github.com/dell/goobjectscale/pkg/client/rest/client"
	iamObjectscale "github.com/dell/goobjectscale/pkg/client/rest/iam"
	log "github.com/sirupsen/logrus"

	"github.com/dell/cosi/pkg/config"
	"github.com/dell/cosi/pkg/transport"
	"github.com/dell/goobjectscale/pkg/client/api"
)

const (
	splitNumber = 2
	// bucketVersion is used when sending the bucket policy update request.
	bucketVersion = "2012-10-17"
	// allowEffect is used when updating the bucket policy, in order to grant permissions to user.
	allowEffect = "Allow"
	// maxUsernameLength is used to trim the username to specific length.
	maxUsernameLength = 64
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

	x509Client := http.Client{Transport: transport}

	objectscaleAuthUser := objectscaleClient.AuthUser{
		Gateway:  objectscaleGateway,
		Username: username,
		Password: password,
	}
	mgmtClient := objectscaleRest.NewClientSet(
		&objectscaleClient.Simple{
			Endpoint:       objectstoreGateway,
			Authenticator:  &objectscaleAuthUser,
			HTTPClient:     &x509Client,
			OverrideHeader: false,
		},
	)

	objClient := objectscaleClient.AuthUser{
		Gateway:  objectscaleGateway,
		Username: username,
		Password: password,
	}

	iamSession, err := session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				CredentialsChainVerboseErrors: aws.Bool(true),
				Credentials:                   credentials.AnonymousCredentials,
				Endpoint:                      aws.String(objectscaleGateway + "/iam"),
				Region:                        region,
				HTTPClient:                    &x509Client,
				Logger:                        internalLogger.New(logrusr.New(log.StandardLogger())),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create new IAM session: %w", err)
	}

	iamClient := iam.New(iamSession)

	err = iamObjectscale.InjectAccountIDToIAMClient(iamClient, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to inject account ID to IAM client: %w", err)
	}

	err = iamObjectscale.InjectTokenToIAMClient(iamClient, &objectscaleAuthUser, x509Client)
	if err != nil {
		return nil, fmt.Errorf("failed to inject token to IAM client: %w", err)
	}

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

func BuildUsername(namespace, bucketName string) string {
	raw := fmt.Sprintf("%v-user-%v", namespace, bucketName)
	if len(raw) > maxUsernameLength {
		raw = raw[:maxUsernameLength]
	}

	return raw
}

// BuildResourceString constructs policy resource string.
func BuildResourceString(objectScaleID, objectStoreID, bucketName string) string {
	return fmt.Sprintf("arn:aws:s3:%s:%s:%s/*", objectScaleID, objectStoreID, bucketName)
}

// BuildPrincipalString constructs policy principal string.
func BuildPrincipalString(namespace, bucketName string) string {
	return fmt.Sprintf("urn:osc:iam::%s:user/%s", namespace, BuildUsername(namespace, bucketName))
}

// GetBucketName splits BucketID by -, the first element is backendID, the second element is bucketName.
func GetBucketName(bucketID string) (string, error) {
	list := strings.SplitN(bucketID, "-", splitNumber)

	if len(list) != 2 || list[1] == "" { //nolint:gomnd
		return "", errors.New("invalid bucketId")
	}

	return list[1], nil
}
