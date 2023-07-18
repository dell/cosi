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

package fake

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	driver "github.com/dell/cosi-driver/pkg/provisioner/virtualdriver"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

// Driver is a mock implementation of virtualdriver.Driver interface.
type Driver struct {
	FakeID string
}

var _ driver.Driver = (*Driver)(nil) // interface guard

const (
	// ForceFail constant can be used to forcefully fail any of the Driver* method from the Driver.
	//
	// DriverCreateBucket, DriverGrantBucketAccess will take this in request.Parameters:
	//
	//	req := &cosi.DriverCreateBucketRequest{
	//		Parameters: map[string]string{
	//			fake.ForceFail: "x",
	//		}
	//	}
	//
	// DriverDeleteBucket, DriverRevokeBucketAccess expect this in BucketId:
	//
	//	req := &cosi.DriverDeleteBucketRequest{
	//		BucketId: fmt.Sprintf("fake-%s", fake.ForceFail)
	//	}
	//
	ForceFail = "X-TEST/force-fail"

	// KeyDriverID is used for retrieving the ID of the virtual driver a request comes to. It's retrieved from the
	// request's parameters, then a driver with such ID is being checked for existence. It's also used for controlling
	// tests flow, e.g. forcibly failing a test.
	KeyDriverID = "id"
)

// ID is implementation of method from virtualdriver.Driver interface.
func (d *Driver) ID() string {
	return d.FakeID
}

// DriverCreateBucket is implementation of method from virtualdriver.Driver interface.
//
// To forcefully fail it, add parameter with Key "X-TEST/force-fail" and any non-zero value.
func (d *Driver) DriverCreateBucket(_ context.Context, req *cosi.DriverCreateBucketRequest) (*cosi.DriverCreateBucketResponse, error) {
	if _, ok := req.Parameters[ForceFail]; ok {
		return nil, status.Error(codes.Internal, "an unexpected error occurred")
	}

	return &cosi.DriverCreateBucketResponse{
		BucketId: fmt.Sprintf("%s-bucket", d.ID()),
	}, nil
}

// DriverDeleteBucket is implementation of method from virtualdriver.Driver interface.
//
// To forcefully fail it set BucketId in request to contain string "X-TEST/force-fail".
func (d *Driver) DriverDeleteBucket(_ context.Context, req *cosi.DriverDeleteBucketRequest) (*cosi.DriverDeleteBucketResponse, error) {
	if strings.Contains(req.BucketId, ForceFail) {
		return nil, status.Error(codes.Internal, "an unexpected error occurred")
	}

	return &cosi.DriverDeleteBucketResponse{}, nil
}

// DriverGrantBucketAccess is implementation of method from virtualdriver.Driver interface.
//
// To forcefully fail it, add parameter with Key "X-TEST/force-fail" and any non-zero value.
func (d *Driver) DriverGrantBucketAccess(_ context.Context, req *cosi.DriverGrantBucketAccessRequest) (*cosi.DriverGrantBucketAccessResponse, error) {
	if _, ok := req.Parameters[ForceFail]; ok {
		return nil, status.Error(codes.Internal, "an unexpected error occurred")
	}

	return &cosi.DriverGrantBucketAccessResponse{
		AccountId: fmt.Sprintf("%s-account", d.ID()),
		Credentials: map[string]*cosi.CredentialDetails{
			"s3": {
				Secrets: map[string]string{
					"endpoint":        "test.endpoint",
					"accessKeyID":     "test-access-key",
					"accessSecretKey": "test-secret-key",
				},
			},
		},
	}, nil
}

// DriverRevokeBucketAccess is implementation of method from virtualdriver.Driver interface.
//
// To forcefully fail it set BucketId in request to contain string "X-TEST/force-fail".
func (d *Driver) DriverRevokeBucketAccess(_ context.Context, req *cosi.DriverRevokeBucketAccessRequest) (*cosi.DriverRevokeBucketAccessResponse, error) {
	if strings.Contains(req.BucketId, ForceFail) {
		return nil, status.Error(codes.Internal, "an unexpected error occurred")
	}

	return &cosi.DriverRevokeBucketAccessResponse{}, nil
}
