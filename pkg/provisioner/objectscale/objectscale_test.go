// Copyright © 2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

package objectscale

import (
	"context"
	"errors"
	"testing"

	"github.com/dell/cosi/pkg/config"
	"github.com/dell/goobjectscale/pkg/client/rest/client/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	namespace     = "test-namespace"
	testCred      = "test-unittest"
	sampleRootCAs = "dGVzdC1yb290LWNhCg=="
)

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		config     *config.Objectscale
		wantErr    bool
		errMessage string
	}{
		{
			name: "Successful creation of new server",
			config: &config.Objectscale{
				Id: "test-id",
				Credentials: config.Credentials{
					Username: "test-username",
					Password: testCred,
				},
				Namespace: &namespace,
				Protocols: config.Protocols{
					S3: &config.S3{
						Endpoint: "s3.objectstore.test",
					},
				},
				Tls: config.Tls{
					Insecure: true,
				},
			},
			wantErr: false,
		},
		{
			name: "Successful creation of new server with certificate",
			config: &config.Objectscale{
				Id: "test-id",
				Credentials: config.Credentials{
					Username: "test-username",
					Password: testCred,
				},
				Namespace: &namespace,
				Protocols: config.Protocols{
					S3: &config.S3{
						Endpoint: "s3.objectstore.test",
					},
				},
				Tls: config.Tls{
					Insecure: false,
					RootCas:  &sampleRootCAs,
				},
			},
			wantErr: false,
		},
		{
			name: "Error with insecure false without providing root CAs",
			config: &config.Objectscale{
				Id: "test-id",
				Credentials: config.Credentials{
					Username: "test-username",
					Password: testCred,
				},
				Namespace: &namespace,
				Protocols: config.Protocols{
					S3: &config.S3{
						Endpoint: "s3.objectstore.test",
					},
				},
				Tls: config.Tls{
					Insecure: false,
				},
			},
			wantErr:    true,
			errMessage: "root certificate authority is missing",
		},
		{
			name: "Error when id is empty",
			config: &config.Objectscale{
				Credentials: config.Credentials{
					Username: "test-username",
				},
			},
			wantErr:    true,
			errMessage: "empty driver id",
		},
		{
			name: "Error when username is empty",
			config: &config.Objectscale{
				Id: "test-id",
			},
			wantErr:    true,
			errMessage: "empty username",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errMessage {
				t.Errorf("New() error message = %v, want %v", err.Error(), tt.errMessage)
			}
		})
	}
}

func TestGetIAMClient(t *testing.T) {
	tests := []struct {
		name             string
		iamClientFactory func() IAMClientFactory
		expectError      bool
	}{
		{
			name: "Get IAM client successfully when logged in",
			iamClientFactory: func() IAMClientFactory {
				authenticator := mocks.NewAuthenticator(t)
				authenticator.On("IsAuthenticated").Return(true)
				authenticator.On("Token").Return("test-token")
				return IAMClientFactory{
					namespace:     "test-namespace",
					username:      "test-unittest",
					password:      "test-password",
					authenticator: authenticator,
				}
			},
			expectError: false,
		},
		{
			name: "Get IAM client successfully when not logged in",
			iamClientFactory: func() IAMClientFactory {
				authenticator := mocks.NewAuthenticator(t)
				authenticator.On("IsAuthenticated").Return(false)
				authenticator.On("Login", mock.Anything, mock.Anything).Return(nil)
				authenticator.On("Token").Return("test-token")
				return IAMClientFactory{
					namespace:     "test-namespace",
					username:      "test-unittest",
					password:      "test-password",
					authenticator: authenticator,
				}
			},
			expectError: false,
		},
		{
			name: "Failed to get IAM client due to login error",
			iamClientFactory: func() IAMClientFactory {
				authenticator := mocks.NewAuthenticator(t)
				authenticator.On("IsAuthenticated").Return(false)
				authenticator.On("Login", mock.Anything, mock.Anything).Return(errors.New("login failure"))
				return IAMClientFactory{
					namespace:     "test-namespace",
					username:      "test-unittest",
					password:      "test-password",
					authenticator: authenticator,
				}
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iamClient, err := tt.iamClientFactory().getIAMClient(context.Background())

			if tt.expectError {
				assert.NotNil(t, err)
				assert.Nil(t, iamClient)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, iamClient)
			}
		})
	}
}

func TestID(t *testing.T) {
	s := Server{
		backendID: "test-backend-id",
	}
	assert.Equal(t, s.ID(), "test-backend-id")
}

func TestBuildUser(t *testing.T) {
	tests := []struct {
		name           string
		namespace      string
		bucketName     string
		expectedOutput string
	}{
		{
			name:           "Successful build user name",
			namespace:      "ns1",
			bucketName:     "bucket1",
			expectedOutput: "ns1-user-bucket1",
		},
		{
			name:           "Success when 64 character limit is exceeded",
			namespace:      "ns1",
			bucketName:     "bucket123456789012345678901234567890123456789012345678901234567890",
			expectedOutput: "ns1-user-bucket1234567890123456789012345678901234567890123456789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := BuildUsername(tt.namespace, tt.bucketName)
			assert.Equal(t, tt.expectedOutput, output)
		})
	}
}
