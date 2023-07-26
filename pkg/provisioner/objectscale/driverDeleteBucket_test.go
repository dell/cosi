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
	"strings"
	"testing"
	"time"

	"github.com/dell/goobjectscale/pkg/client/fake"
	"github.com/dell/goobjectscale/pkg/client/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

// testDriverCreateBucket tests bucket deletion functionality on ObjectScale platform.
func TestDriverDeleteBucket(t *testing.T) {
	const (
		namespace = "namespace"
		testID    = "test.id"
	)

	testCases := []struct {
		description   string
		inputBucketID string
		expectedError error
		server        Server
	}{
		{
			description:   "invalid bucketID",
			inputBucketID: "",
			expectedError: status.Error(codes.InvalidArgument, "empty bucketID"),
		},
		{
			description:   "bucket does not exist",
			inputBucketID: strings.Join([]string{testID, "bucket-valid"}, "-"),
			expectedError: nil,
			server: Server{
				mgmtClient: fake.NewClientSet(),
				namespace:  namespace,
				backendID:  testID,
			},
		},
		{
			description:   "failed to delete bucket",
			inputBucketID: strings.Join([]string{testID, "bucket-invalid-FORCEFAIL"}, "-"),
			expectedError: status.Error(codes.Internal, "bucket was not successfully deleted"),
			server: Server{
				mgmtClient: fake.NewClientSet(&model.Bucket{
					Name:      "bucket-valid",
					Namespace: namespace,
				}),
				namespace:   namespace,
				backendID:   testID,
				emptyBucket: true,
			},
		},
		{
			description:   "bucket successfully deleted",
			inputBucketID: strings.Join([]string{testID, "bucket-valid"}, "-"),
			expectedError: nil,
			server: Server{
				mgmtClient: fake.NewClientSet(&model.Bucket{
					Name:      "bucket-valid",
					Namespace: namespace,
				}),
				namespace:   namespace,
				backendID:   testID,
				emptyBucket: true,
			},
		},
	}

	for _, scenario := range testCases {
		t.Run(scenario.description, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err := scenario.server.DriverDeleteBucket(ctx, &cosi.DriverDeleteBucketRequest{BucketId: scenario.inputBucketID})
			assert.ErrorIs(t, err, scenario.expectedError, err)
		})
	}
}
