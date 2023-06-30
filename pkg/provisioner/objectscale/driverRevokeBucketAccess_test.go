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
	"github.com/dell/goobjectscale/pkg/client/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/dell/goobjectscale/pkg/client/api/mocks"
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

	accessKeyList := make([]*iam.AccessKeyMetadata, 1)
	accessKeyList[0] = &iam.AccessKeyMetadata{
		AccessKeyId: aws.String("abc"),
		UserName:    aws.String("namespace-user-valid"),
	}

	IAMClient.On("ListAccessKeys", mock.Anything).Return(&iam.ListAccessKeysOutput{
		AccessKeyMetadata: accessKeyList}, nil).Once()
	IAMClient.On("DeleteAccessKey", mock.Anything).Return(nil, nil).Once()
	IAMClient.On("DeleteUser", mock.Anything).Return(nil, nil).Once()
	IAMClient.On("GetUser", mock.Anything).Return(&iam.GetUserOutput{
		User: &iam.User{
			UserName: aws.String("namespace-user-valid"),
		},
	}, nil).Once()

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{
		Name:      "valid",
		Namespace: testNamespace,
	}, nil).Once()
	bucketsMock.On("GetPolicy", mock.Anything, mock.Anything, mock.Anything).Return(`{"Id":"S3PolicyId1","Version":"2012-10-17","Statement":[{"Resource":["arn:aws:s3:osci5b022e718aa7e0ff:osti202e682782ebcbfd:lynxbucket/*"],"Sid":"GetObject_permission","Effect":"Allow","Principal":{"AWS":["urn:osc:iam::osai07c2ae318ae9d6f2:user/iam_user20230523061025118"]},"Action":["s3:GetObjectVersion"]}]}`, nil).Once()
	bucketsMock.On("UpdatePolicy", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	server := Server{
		mgmtClient:    mgmtClientMock,
		iamClient:     IAMClient,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverRevokeBucketAccessRequest{
		BucketId:  "bucket-valid",
		AccountId: "namespace-user-valid",
	}

	response, err := server.DriverRevokeBucketAccess(ctx, req)
	assert.ErrorIs(t, err, nil, err)
	assert.NotNil(t, response)
}

// func testEmptyBucketID(t *testing.T) {

// }

func testEmptyAccountID(t *testing.T) {

}

// 1. empty accountID
// 2. empty bucketID
// 3. failed to check bucket existence
// 4. bucket not found
// 5. failed to check for user existence
// 6. failed to get user
// 7. failed to get access key list
// 8. failed to delete access key
// 9. failed to check bucket policy existence
// 10. empty policy
// 11. failed to marshal updatePolicy into JSON
// 12. failed to update bucket policy
// 13. failed to delete user

// &UpdateBucketPolicyRequest{
// 	PolicyID: "policy1",
// 	Version: "v1",
// 	Statement: &UpdateBucketPolicyStatement{
// 		Resource: []string{"arn:aws:s3:osci5b022e718aa7e0ff:osti202e682782ebcbfd:valid/*"},
// 		SID: "sid",
// 		Effect: "effect",
// 		Principal: &Principal{
// 			AWS: []string{"urn:osc:iam::osai07c2ae318ae9d6f2:user/iam_user20230523061025118"},
// 			Action: []string{"s3:GetObjectVersion"},
// 		},
// 	},
// },
