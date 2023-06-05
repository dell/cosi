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

package model

import "encoding/xml"

// TenantInfo is an object store tenant with an alternate XML tag name
type TenantInfo struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"tenant_info"`

	Tenant
}

// Tenant is an object store tenant
type Tenant struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"tenant"`

	// ID is the id of the tenant scoped to the cluster instance
	ID string `json:"id,omitempty" xml:"id,omitempty"`

	// EncryptionEnabled displays if this tenant is configured for encryption at rest
	EncryptionEnabled bool `json:"is_encryption_enabled,omitempty" xml:"is_encryption_enabled,omitempty"`

	// ComplianceEnabled displays if this tenant is configured for compliance retention
	ComplianceEnabled bool `json:"is_compliance_enabled,omitempty" xml:"is_compliance_enabled,omitempty"`

	// ReplicationGroup is the default replication group id of the tenant
	ReplicationGroup string `json:"default_data_services_vpool,omitempty" xml:"default_data_services_vpool,omitempty"`

	// BucketBlockSize is the default bucket size at which new object creations will be blocked
	BucketBlockSize int64 `json:"default_bucket_block_size,omitempty" xml:"default_bucket_block_size,omitempty"`

	RetentionClasses        string `xml:"retention_classes"`
	NotificationSize        string `xml:"notificationSize"`
	BlockSize               string `xml:"blockSize"`
	BlockSizeInCount        string `xml:"blockSizeInCount"`
	NotificationSizeInCount string `xml:"notificationSizeInCount"`
	Alias                   string `xml:"alias"`
}

// TenantCreate is an object store tenant creation input
type TenantCreate struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"tenant_create"`

	// Alias is the tenant alias to set
	Alias string `xml:"alias"`
	// ID is the id of the tenant scoped to the cluster instance
	AccountID string `json:"account_id" xml:"account_id"`

	// EncryptionEnabled displays if this tenant is configured for encryption at rest
	EncryptionEnabled bool `json:"is_encryption_enabled,omitempty" xml:"is_encryption_enabled"`

	// ComplianceEnabled displays if this tenant is configured for compliance retention
	ComplianceEnabled bool `json:"is_compliance_enabled,omitempty" xml:"is_compliance_enabled"`

	// BucketBlockSize is the default bucket size at which new object creations will be blocked
	BucketBlockSize int64 `json:"default_bucket_block_size,omitempty" xml:"default_bucket_block_size"`
}

// TenantUpdate is an object store tenant update input
type TenantUpdate struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"tenant_update"`

	// BucketBlockSize is the default bucket size at which new object creations will be blocked
	BucketBlockSize int64 `json:"default_bucket_block_size,omitempty" xml:"default_bucket_block_size"`

	// Alias is the tenant alias to set
	Alias string `xml:"alias"`
}

// TenantList is a list of object store tenants
type TenantList struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `json:"tenants" xml:"tenants"`

	// Items is the list of tenants in the list
	Items []Tenant `json:"tenant" xml:"tenant"`
}

// TenantQuota is an object store tenant quota
type TenantQuota struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"tenant_quota_details"`

	BlockSize string `xml:"blockSize"`

	NotificationSize string `xml:"notificationSize"`

	BlockSizeInCount string `xml:"blockSizeInCount"`

	NotificationSizeInCount string `xml:"notificationSizeInCount"`

	ID string `json:"id,omitempty" xml:"id,omitempty"`
}

// TenantQuotaSet is an object store tenant quota
type TenantQuotaSet struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"tenant_quota_details"`

	BlockSize string `xml:"blockSize"`

	NotificationSize string `xml:"notificationSize"`

	BlockSizeInCount string `xml:"blockSizeInCount"`

	NotificationSizeInCount string `xml:"notificationSizeInCount"`
}
