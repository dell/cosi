//
//
//  Copyright © 2021 - 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//       http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//
//

package model

import "encoding/xml"

// BucketInfo is an object storage bucket with an alternate XML tag name
type BucketInfo struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"bucket_info"`

	Bucket
}

// BucketCreate is an object storage bucket with an alternate XML tag name
type BucketCreate struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"object_bucket_create"`

	Bucket
}

// BucketQuotaUpdate is the struct of quota updating
type BucketQuotaUpdate struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"bucket_quota_param"`

	BucketQuota
}

// BucketQuotaInfo is the struct of quota information
type BucketQuotaInfo struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"bucket_quota_details"`

	BucketQuota
}

// BucketQuota is quota struct
type BucketQuota struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"bucket_quota"`

	// Name is the name of the cluster instance
	BucketName string `json:"bucketname,omitempty" xml:"bucketname"`

	// Namespace is the namespace of the bucket
	Namespace string `json:"namespace,omitempty" xml:"namespace,omitempty"`

	// BlockSize is the bucket size at which new object creations will be blocked
	BlockSize int64 `json:"blockSize,omitempty" xml:"blockSize,omitempty"`

	// BlockSizeCount is the bucket size, in counts, at which new object creations will be blocked
	BlockSizeCount int64 `json:"blockSizeInCount,omitempty" xml:"blockSizeInCount,omitempty"`

	// NotificationSize is the bucket size at which the users will be notified
	NotificationSize int64 `json:"notificationSize,omitempty" xml:"notificationSize,omitempty"`

	// NotificationSize is the bucket size, in counts, at which the users will be notified
	NotificationSizeCount int64 `json:"notificationSizeInCount,omitempty" xml:"notificationSizeInCount,omitempty"`
}

// Bucket is an object storage bucket
type Bucket struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"object_bucket"`

	// APIType is the object API type used by the bucket
	APIType string `json:"api_type,omitempty" xml:"api_type,omitempty"`

	// AuditDeleteExpiration is the amount of time to retain deletion audit
	// entries
	AuditDeleteExpiration int `json:"audit_delete_expiration,omitempty" xml:"audit_delete_expiration,omitempty"`

	// Created is the date and time that the bucket was created
	Created string `json:"created,omitempty" xml:"created,omitempty"`

	// ID is the id of the bucket scoped to the cluster instance
	ID string `json:"id,omitempty" xml:"id,omitempty"`

	// Name is the name of the cluster instance
	Name string `json:"name" xml:"name"`

	// EncryptionEnabled displays if this bucket is configured for encryption
	// at rest
	EncryptionEnabled bool `json:"is_encryption_enabled,omitempty" xml:"is_encryption_enabled,omitempty"`

	// SoftQuota is the warning quota level for the bucket
	SoftQuota string `json:"softquota,omitempty" xml:"softquota,omitempty"`

	// FSEnabled indicates if the bucket has file-system support enabled
	FSEnabled bool `json:"fs_access_enabled,omitempty" xml:"fs_access_enabled,omitempty"`

	// Locked indicates if the bucket is locked
	Locked bool `json:"locked,omitempty" xml:"locked,omitempty"`

	// ReplicationGroup is the replication group id of the bucket
	ReplicationGroup string `json:"vpool,omitempty" xml:"vpool,omitempty"`

	// Namespace is the namespace of the bucket
	Namespace string `json:"namespace,omitempty" xml:"namespace,omitempty"`

	// Owner is the s3 object user owner of the bucket
	Owner string `json:"owner,omitempty" xml:"owner,omitempty"`

	// StaleAllowed indicates if access to the bucket is allowed during an
	// outage
	StaleAllowed bool `json:"is_stale_allowed,omitempty" xml:"is_stale_allowed,omitempty"`

	// TSOReadOnly indicates if access to the bucket is allowed during an
	// outage
	TSOReadOnly bool `json:"is_tso_read_only,omitempty" xml:"is_tso_read_only,omitempty"`

	// DefaultRetention is the default retention period for objects in bucket
	DefaultRetention int64 `json:"default_retention,omitempty" xml:"default_retention,omitempty"`

	// BlockSize is the bucket size at which new object creations will be blocked
	BlockSize int64 `json:"block_size,omitempty" xml:"block_size,omitempty"`

	// BlockSizeCount is the bucket size, in counts, at which new object creations will be blocked
	BlockSizeCount int64 `json:"block_size_in_count,omitempty" xml:"block_size_in_count,omitempty"`

	// NotificationSize is the bucket size at which the users will be notified
	NotificationSize int64 `json:"notification_size,omitempty" xml:"notification_size,omitempty"`

	// NotificationSize is the bucket size, in counts, at which the users will be notified
	NotificationSizeCount int64 `json:"notification_size_in_count,omitempty" xml:"notification_size_in_count,omitempty"`

	// BlockSizeInput is the input of bucket size, support CreateBucket method
	BlockSizeInput int64 `json:"blockSize,omitempty" xml:"blockSize,omitempty"`

	// BlockSizeCountInput is the bucket size, in counts, at which new object creations will be blocked
	BlockSizeCountInput int64 `json:"blockSizeInCount,omitempty" xml:"blockSizeInCount,omitempty"`

	// NotificationSizeInput is the input of notification size, support CreateBucket method
	NotificationSizeInput int64 `json:"notificationSize,omitempty" xml:"notificationSize,omitempty"`

	// NotificationSizeCountInput is the bucket size, in counts, at which the users will be notified
	NotificationSizeCountInput int64 `json:"notificationSizeInCount,omitempty" xml:"notificationSizeInCount,omitempty"`

	// Tags is a list of arbitrary metadata keys and values applied to the
	// bucket
	Tags TagSet `json:"TagSet,omitempty" xml:"TagSet,omitempty"`

	// Retention is the default retention value for the bucket
	Retention int64 `json:"retention,omitempty" xml:"retention,omitempty"`

	// DefaultGroupFileReadPermission is a flag indicating the Read permission
	// for default group
	DefaultGroupFileReadPermission bool `json:"default_group_file_read_permission,omitempty" xml:"default_group_file_read_permission,omitempty"`

	// DefaultGroupFileWritePermission is a flag indicating the Execute permission
	// for default group
	DefaultGroupFileExecutePermission bool `json:"default_group_file_execute_permission,omitempty" xml:"default_group_file_execute_permission,omitempty"`

	// DefaultGroupFileExecutePermission is a flag indicating the Write permission
	// for default group
	DefaultGroupFileWritePermission bool `json:"default_group_file_write_permission,omitempty" xml:"default_group_file_write_permission,omitempty"`

	// DefaultGroupDirReadPermission is a flag indicating the Read permission
	// for default group
	DefaultGroupDirReadPermission bool `json:"default_group_dir_read_permission,omitempty" xml:"default_group_dir_read_permission,omitempty"`

	// DefaultGroupDirWritePermission is a flag indicating the Execute permission
	// for default group
	DefaultGroupDirExecutePermission bool `json:"default_group_dir_execute_permission,omitempty" xml:"default_group_dir_execute_permission,omitempty"`

	// DefaultGroupDirExecutePermission is a flag indicating the Write permission
	// for default group
	DefaultGroupDirWritePermission bool `json:"default_group_dir_write_permission,omitempty" xml:"default_group_dir_write_permission,omitempty"`

	// DefaultGroup is the bucket's default group
	DefaultGroup string `json:"default_group,omitempty" xml:"default_group,omitempty"`

	// SearchMetadata is the custom metadata for enabled for querying on the
	// bucket
	SearchMetadata `json:"search_metadata,omitempty" xml:"search_metadata,omitempty"`

	// MinMaxGovenor enforces minimum and maximum retention for bucket objects
	MinMaxGovenor `json:"min_max_govenor,omitempty" xml:"min_max_govenor,omitempty"`

	// StoragePolicy is the default storage policy of the bucket
	StoragePolicy string `json:"storagePolicy,omitempty" xml:"storage_policy,omitempty"`
}

// BucketList is a list of object storage buckets
type BucketList struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `json:"object_buckets" xml:"object_buckets"`

	// Items is the list of buckets in the list
	Items []Bucket `json:"object_bucket" xml:"object_bucket"`

	// MaxBuckets is the maximum number of buckets requested in the listing
	MaxBuckets int `json:"max_buckets,omitempty" xml:"MaxBuckets,omitempty"`

	// NextMarker is a reference object to receive the next set of buckets
	NextMarker string `json:"next_marker,omitempty" xml:"NextMarker,omitempty"`

	// Filter is a string query used to limit the returned buckets in the
	// listing
	Filter string `json:"Filter,omitempty" xml:"Filter,omitempty"`

	// NextPageLink is a hyperlink to the next page in the bucket listing
	NextPageLink string `json:"next_page_link,omitempty" xml:"NextPageLink,omitempty"`
}

// MinMaxGovenor enforces minimum and maximum retention for bucket objects
type MinMaxGovenor struct {

	// EnforceRetention indicates if retention should be enforced for this
	// min-max-govenor
	EnforceRetention bool `json:"enforce_retention" xml:"enforce_retention"`

	// MinimumFixedRetention  is the minimum fixed retention for objects within
	// a bucket
	MinimumFixedRetention int64 `json:"minimum_fixed_retention" xml:"minimum_fixed_retention"`

	// MinimumVariableRetention  is the minimum variable retention for objects
	// within a bucket
	MinimumVariableRetention int64 `json:"minimum_variable_retention" xml:"minimum_variable_retention"`

	// MaximumFixedRetention  is the maximum fixed retention for objects within
	// a bucket
	MaximumFixedRetention int64 `json:"maximum_fixed_retention" xml:"maximum_fixed_retention"`

	// MaximumVariableRetention  is the maximum variable retention for objects
	// within a bucket
	MaximumVariableRetention int64 `json:"maximum_variable_retention" xml:"maximum_variable_retention"`

	// Link is the hyperlink to this resource
	Link `json:"link" xml:"link"`

	// Inactive indicates if the bucket has been placed into an inactive state,
	// typically prior to deletion
	Inactive bool `json:"inactive" xml:"inactive"`

	// Global indicates if the resource is global
	Global bool `json:"global" xml:"global"`

	// Remote indicates if the resource is remote to the current API instance
	Remote bool `json:"remote" xml:"remote"`

	// Internal indicates if the resource is an internal resource
	Internal bool `json:"internal" xml:"internal"`

	// VDCLink is a link from a bucket to a VDC
	VDCLink `json:"vdc" xml:"vdc"`
}

// VDCLink is a link from a bucket to a VDC
type VDCLink struct {
	// ID is the identifier for the VDC
	ID string `json:"id" xml:"id"`

	// Link is a hyperlink to the VDC
	Link `json:"link" xml:"link"`
}
