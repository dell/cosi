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
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/dell/cosi-driver/pkg/iamfaketoo"
	"github.com/dell/cosi-driver/pkg/internal/testcontext"
	"github.com/dell/goobjectscale/pkg/client/fake"
	"github.com/dell/goobjectscale/pkg/client/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

var _ iamiface.IAMAPI = (*iamfaketoo.IAMAPI)(nil)

func TestServerBucketAccessRevoke(t *testing.T) {
	t.Parallel()

	for scenario, fn := range map[string]func(t *testing.T){
		"testValidAccessRevoking": testValidAccessRevoking,
	} {
		fn := fn

		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			fn(t)
		})
	}
}

func testValidAccessRevoking(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	// That's how we can mock the objectscale IAM api client
	IAMClient := iamfaketoo.NewIAMAPI(t)
	IAMClient.On("CreateUserWithContext", mock.Anything, mock.Anything).Return(
		&iam.CreateUserOutput{
			User: &iam.User{
				UserName: aws.String("namespace-user-valid"), // This mocked response is based on `namesapce` from server and bucketId from request
			},
		}, nil).Once()
	IAMClient.On("GetUser", mock.Anything).Return(nil, nil).Once()
	IAMClient.On("CreateAccessKey", mock.Anything).Return(&iam.CreateAccessKeyOutput{AccessKey: &iam.AccessKey{AccessKeyId: aws.String("acc"), SecretAccessKey: aws.String("sec")}}, nil).Once()

	server := Server{
		mgmtClient: fake.NewClientSet(&model.Bucket{ // That's how we can mock the objectscale bucket api client
			Name:      "valid", // This is based on "bucket-valid" BucketId from request
			Namespace: testNamespace,
		}),
		iamClient:     IAMClient, // Inject mocked IAM client
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverRevokeBucketAccessRequest{
		BucketId:  "valid-bucket",
		AccountId: "namespace-user-valid",
	}

	response, err := server.DriverRevokeBucketAccess(ctx, req)
	assert.ErrorIs(t, err, nil, err)
	assert.NotNil(t, response)
}
