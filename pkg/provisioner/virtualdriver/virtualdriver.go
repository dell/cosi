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

// TODO: write documentation comment for virtualdriver package
package virtualdriver

import (
	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

// Driver is a interface extending cosi.ProvisionerServer interface by ID method
//
// ProvisionerServer is the server API for Provisioner service, which controls the full lifecycle of
// buckets and bucket accesses on object storage provider.
type Driver interface {
	// each driver must implement default ProvisionerServer interface specified by COSI specification
	cosi.ProvisionerServer

	// additionally, driver should return ID, specific to the platform, that allows to identify which platform,
	// and which hardware OSP this driver is configured to support.
	//
	// E.g. for ObjectScale, this should be ObjectScaleID/ObjectStoreID
	//
	// ID value should be stored in:
	// - CreateBucketRequest.Parameters["id"]
	// - GrantBucketAccessRequest.Parameters["id"]
	//
	// ID value should be extracted from:
	// - DeleteBucketRequest.BucketID
	// - DriverRevokeBucketAccessRequest.BucketID
	ID() string
}
