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
	"strings"

	"github.com/dell/goobjectscale/pkg/client/model"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

const (
	testBucketName = "test_bucket"
	testNamespace  = "namespace"
	testID         = "test.id"
	objectScaleID  = "gateway.objectscale.test"
	objectStoreID  = "objectstore"
	testPolicy     = `{"Id":"S3PolicyId1","Version":"2012-10-17","Statement":[{"Resource":["arn:aws:s3:osci5b022e718aa7e0ff:osti202e682782ebcbfd:lynxbucket/*"],"Sid":"GetObject_permission","Effect":"Allow","Principal":{"AWS":["urn:osc:iam::osai07c2ae318ae9d6f2:user/iam_user20230523061025118"]},"Action":["s3:GetObjectVersion"]}]}`
	testUserName   = "test_user"
)

var (
	testBucket = &model.Bucket{
		Namespace: testNamespace,
		Name:      testBucketName,
	}

	testBucketCreationRequest = &cosi.DriverCreateBucketRequest{
		Name: testBucketName,
	}

	testBucketDeletionRequest = &cosi.DriverDeleteBucketRequest{
		BucketId: strings.Join([]string{testID, testBucketName}, "-"),
	}

	testBucketDeletionRequestEmptyBucketID = &cosi.DriverDeleteBucketRequest{
		BucketId: "",
	}
)
