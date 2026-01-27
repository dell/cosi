// Copyright © 2022 - 2025 Dell Inc. or its subsidiaries. All Rights Reserved.
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
	"reflect"
	"strings"
	"testing"

	omocks "github.com/dell/cosi/pkg/provisioner/objectscale/mocks"
	"github.com/dell/cosi/pkg/provisioner/policy"

	"github.com/dell/cosi/pkg/internal/testcontext"
	"github.com/dell/goobjectscale/pkg/client/api/mocks"
	"github.com/dell/goobjectscale/pkg/client/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	cosi "sigs.k8s.io/container-object-storage-interface/proto"
)

var testBucketRevokeAccessRequest = &cosi.DriverRevokeBucketAccessRequest{
	BucketId:  strings.Join([]string{testID, testBucketName}, "-"),
	AccountId: "account-id",
}

func TestServerDriverRevokeAccess(t *testing.T) {
	t.Parallel()

	for scenario, fn := range map[string]func(t *testing.T){
		// happy path
		"RevokeAccess":                     testDriverRevokeBucketAccess,
		"RevokeAccessWithMultiplePolicies": testDriverRevokeBucketAccessWithMultiplePolicies,
	} {
		fn := fn

		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			fn(t)
		})
	}
}

func testDriverRevokeBucketAccess(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	resourceARNs := BuildResourceStrings("bucket-name")
	awsPrincipalString := BuildPrincipalString(testBucketRevokeAccessRequest.AccountId, testNamespace)

	bucketPolicy := policy.Document{
		Version: "2012-10-17",
		Statement: []policy.StatementEntry{
			{
				Effect:    "Allow",
				Action:    []string{"*"},
				Resource:  resourceARNs,
				Principal: map[string]string{"AWS": awsPrincipalString},
				Sid:       PolicySid,
			},
		},
	}
	bucketPolicyJSON, err := json.Marshal(bucketPolicy)
	assert.Nil(t, err)

	bucketsMock := mocks.NewBucketServiceInterface(t)
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{}, nil).Once()
	bucketsMock.On("GetPolicy", mock.Anything, mock.Anything, mock.Anything).Return(string(bucketPolicyJSON), nil).Once()
	bucketsMock.On("DeletePolicy", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	mgmtClientMock := mocks.NewClientSet(t)
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	iamMock := omocks.NewIAM(t)
	iamMock.On("GetUser", mock.Anything, mock.Anything).Return(nil, errors.New("error")).Once()
	iamMock.On("ListAccessKeys", mock.Anything, mock.Anything).Return(&iam.ListAccessKeysOutput{
		AccessKeyMetadata: []types.AccessKeyMetadata{
			{
				AccessKeyId: aws.String("key"),
			},
		},
	}, nil).Once()
	iamMock.On("DeleteAccessKey", mock.Anything, mock.Anything).Return(nil, nil).Once()
	iamMock.On("DeleteUser", mock.Anything, mock.Anything).Return(nil, nil).Once()

	val := func(context.Context) (IAM, error) {
		return iamMock, nil
	}

	server := Server{
		mgmtClient: mgmtClientMock,
		namespace:  testNamespace,
		backendID:  testID,
		iamClient:  val,
	}

	res, err := server.DriverRevokeBucketAccess(ctx, testBucketRevokeAccessRequest)

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func testDriverRevokeBucketAccessWithMultiplePolicies(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	resourceARNs := BuildResourceStrings("bucket-name")
	awsPrincipalString := BuildPrincipalString(testBucketRevokeAccessRequest.AccountId, testNamespace)

	// policy has additional entries, so we expect an Update of this policy and not to Delete
	bucketPolicy := policy.Document{
		Version: "2012-10-17",
		Statement: []policy.StatementEntry{
			{
				Effect:    "Allow",
				Action:    []string{"*"},
				Resource:  resourceARNs,
				Principal: map[string]string{"AWS": awsPrincipalString},
				Sid:       PolicySid,
			},
			{
				Effect:    "Allow",
				Action:    []string{"*"},
				Resource:  []string{"existing-resource"},
				Principal: map[string]string{"AWS": "existing-principal"},
				Sid:       PolicySid,
			},
		},
	}
	bucketPolicyJSON, err := json.Marshal(bucketPolicy)
	assert.Nil(t, err)

	bucketsMock := mocks.NewBucketServiceInterface(t)
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{}, nil).Once()
	bucketsMock.On("GetPolicy", mock.Anything, mock.Anything, mock.Anything).Return(string(bucketPolicyJSON), nil).Once()
	bucketsMock.On("UpdatePolicy", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	mgmtClientMock := mocks.NewClientSet(t)
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	iamMock := omocks.NewIAM(t)
	iamMock.On("GetUser", mock.Anything, mock.Anything).Return(nil, errors.New("error")).Once()
	iamMock.On("ListAccessKeys", mock.Anything, mock.Anything).Return(&iam.ListAccessKeysOutput{
		AccessKeyMetadata: []types.AccessKeyMetadata{
			{
				AccessKeyId: aws.String("key"),
			},
		},
	}, nil).Once()
	iamMock.On("DeleteAccessKey", mock.Anything, mock.Anything).Return(nil, nil).Once()
	iamMock.On("DeleteUser", mock.Anything, mock.Anything).Return(nil, nil).Once()

	val := func(context.Context) (IAM, error) {
		return iamMock, nil
	}

	server := Server{
		mgmtClient: mgmtClientMock,
		namespace:  testNamespace,
		backendID:  testID,
		iamClient:  val,
	}

	res, err := server.DriverRevokeBucketAccess(ctx, testBucketRevokeAccessRequest)

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestKvToFields(t *testing.T) {
	tests := []struct {
		name      string
		keyValues []any
		want      map[string]any
	}{
		{
			name:      "empty key values",
			keyValues: []any{},
			want:      map[string]any{},
		},
		{
			name:      "single key value pair",
			keyValues: []any{"a", 1},
			want:      map[string]any{"a": 1},
		},
		{
			name:      "multiple key value pairs with different values",
			keyValues: []any{"a", 1, "b", "two", "c", 3.0},
			want:      map[string]any{"a": 1, "b": "two", "c": 3.0},
		},
		{
			name:      "non-string key is skipped",
			keyValues: []any{123, "value", "ok", true},
			want:      map[string]any{"ok": true},
		},
		{
			name:      "odd length drops last value",
			keyValues: []any{"a", 1, "b"},
			want:      map[string]any{"a": 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := kvToFields(tt.keyValues...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
