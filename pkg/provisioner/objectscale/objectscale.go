// Copyright © 2022-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

package objectscale

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	cosi "sigs.k8s.io/container-object-storage-interface/proto"

	"github.com/aws/smithy-go/middleware"

	obsConfig "github.com/dell/cosi/pkg/config"
	"github.com/dell/cosi/pkg/internal/transport"
	logger "github.com/dell/cosi/pkg/logger"
	"github.com/pkg/errors"

	driver "github.com/dell/cosi/pkg/provisioner/virtualdriver"
	"github.com/dell/goobjectscale/pkg/client/api"
	"github.com/dell/goobjectscale/pkg/client/rest"
	"github.com/dell/goobjectscale/pkg/client/rest/client"
	obsconfig "github.com/dell/goobjectscale/pkg/config"
)

const (
	// bucketVersion is used when sending the bucket policy update request.
	bucketVersion = "2012-10-17"
	// allowEffect is used when updating the bucket policy, in order to grant permissions to user.
	allowEffect = "Allow"
	// maxUsernameLength is used to trim the username to specific length.
	maxUsernameLength = 64

	CreateBucketTraceName       = "CreateBucketRequest"
	DeleteBucketTraceName       = "DeleteBucketRequest"
	GrantBucketAccessTraceName  = "GrantBucketAccessRequest"
	RevokeBucketAccessTraceName = "RevokeBucketAccessRequest"
)

type Server struct {
	mgmtClient  api.ClientSet
	backendID   string
	emptyBucket bool
	namespace   string
	s3Endpoint  string
	iamClient   func(context.Context) (IAM, error)
	cosi.UnimplementedProvisionerServer
}

// Since aws-v2 doesn't provide interface, define our own for ease of mocking and unit testing
//
//go:generate go run github.com/vektra/mockery/v2@latest --all
type IAM interface {
	CreateAccessKey(ctx context.Context, params *iam.CreateAccessKeyInput, optFns ...func(*iam.Options)) (*iam.CreateAccessKeyOutput, error)
	CreateUser(ctx context.Context, params *iam.CreateUserInput, optFns ...func(*iam.Options)) (*iam.CreateUserOutput, error)
	DeleteAccessKey(ctx context.Context, params *iam.DeleteAccessKeyInput, optFns ...func(*iam.Options)) (*iam.DeleteAccessKeyOutput, error)
	DeleteUser(ctx context.Context, params *iam.DeleteUserInput, optFns ...func(*iam.Options)) (*iam.DeleteUserOutput, error)
	GetUser(ctx context.Context, params *iam.GetUserInput, optFns ...func(*iam.Options)) (*iam.GetUserOutput, error)
	ListAccessKeys(ctx context.Context, params *iam.ListAccessKeysInput, optFns ...func(*iam.Options)) (*iam.ListAccessKeysOutput, error)
}

var _ driver.Driver = (*Server)(nil)

func New(objConfig *obsConfig.Objectscale) (*Server, error) {
	log.Info("Initializing driver")
	id := objConfig.Id
	if id == "" {
		return nil, errors.New("empty driver id")
	}

	username := objConfig.Credentials.Username
	if username == "" {
		return nil, errors.New("empty username")
	}

	password := objConfig.Credentials.Password
	if password == "" {
		return nil, errors.New("empty password")
	}

	mgmtConfig := &obsconfig.MgmtConfig{
		Username:    username,
		EndpointURL: objConfig.MgmtEndpoint,
	}

	log.Debugf("Management config has been set %v", mgmtConfig)
	// set Password here to prevent it from being printed in the logs
	mgmtConfig.Password = password

	protocolS3Endpoint := objConfig.Protocols.S3.Endpoint
	if protocolS3Endpoint == "" {
		return nil, errors.New("empty protocol S3 endpoint")
	}

	baseTransport, err := transport.New(objConfig.Tls)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Transport: baseTransport,
	}

	objectscaleAuthUser := client.AuthUser{
		Gateway:  mgmtConfig.EndpointURL,
		Username: mgmtConfig.Username,
		Password: mgmtConfig.Password,
	}

	simpleClient := &client.Simple{
		Endpoint:       mgmtConfig.EndpointURL,
		Authenticator:  &objectscaleAuthUser,
		OverrideHeader: false,
		HTTPClient:     httpClient,
	}

	simpleClient.SetLogger(logger.Log())
	clientset := rest.NewClientSet(simpleClient)

	log.Info("Driver has been successfully initialized")
	iamFactory := &IAMClientFactory{
		namespace:     *objConfig.Namespace,
		client:        httpClient,
		endpoint:      objConfig.MgmtEndpoint,
		authenticator: &objectscaleAuthUser,
		username:      mgmtConfig.Username,
		password:      mgmtConfig.Password,
	}

	return &Server{
		mgmtClient:  clientset,
		backendID:   id,
		emptyBucket: objConfig.EmptyBucket,
		namespace:   *objConfig.Namespace,
		s3Endpoint:  protocolS3Endpoint,
		iamClient:   iamFactory.getIAMClient,
	}, nil
}

type IAMClientFactory struct {
	namespace     string
	username      string
	password      string
	endpoint      string
	authenticator client.Authenticator
	client        *http.Client
}

func (i IAMClientFactory) getIAMClient(ctx context.Context) (IAM, error) {
	if !i.authenticator.IsAuthenticated() {
		err := i.authenticator.Login(ctx, i.client)
		if err != nil {
			return nil, err
		}
	}

	iamConfig, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			i.username, i.password, "")),
		config.WithHTTPClient(i.client),
		config.WithAPIOptions([]func(*middleware.Stack) error{
			smithyhttp.AddHeaderValue("X-Emc-Namespace", i.namespace),
			smithyhttp.AddHeaderValue("X-Sds-Auth-Token", i.authenticator.Token()),
		}),
	)
	if err != nil {
		return nil, err
	}
	iamClient := iam.NewFromConfig(iamConfig, func(o *iam.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("%s/iam/", i.endpoint))
	})
	return iamClient, nil
}

// ID extends COSI interface by adding ID method.
func (s *Server) ID() string {
	return s.backendID
}

func BuildUsername(namespace, access string) string {
	raw := fmt.Sprintf("%v-user-%v", namespace, access)
	if len(raw) > maxUsernameLength {
		raw = raw[:maxUsernameLength]
	}

	return raw
}
