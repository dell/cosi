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
	"testing"

	"github.com/dell/goobjectscale/pkg/client/model"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

const (
	testBucketName = "test_bucket"
	testNamespace  = "namespace"
	testID         = "test.id"
	objectScaleID  = "gateway.objectscale.test"
	objectStoreID  = "objectstore"
)

var (
	testBucket = &model.Bucket{
		Namespace: testNamespace,
		Name:      testBucketName,
	}

	testRequest = &cosi.DriverCreateBucketRequest{
		Name: testBucketName,
	}
)

// testContext creates new context with deadline equal to test deadline, or (if the deadline is empty),
// with a timeout equal to the default timeout.
func testContext(t *testing.T) (context.Context, context.CancelFunc) {
	deadline, ok := t.Deadline()

	if ok {
		return context.WithDeadline(context.Background(), deadline)
	}

	return context.WithTimeout(context.Background(), defaultTimeout)
}
