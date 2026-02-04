// Copyright © 2023-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

// Package virtualdriver implements extension of provisioner server
// allowing for usage with multiple platforms.
package virtualdriver

import (
	cosi "sigs.k8s.io/container-object-storage-interface/proto"
)

// Driver is an interface extending cosi.ProvisionerServer interface by ID method
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
