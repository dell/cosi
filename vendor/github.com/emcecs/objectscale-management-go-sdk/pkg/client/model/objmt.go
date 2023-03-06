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

import (
	"encoding/xml"
)

// AccountBillingInfoList list of the billing info per users
type AccountBillingInfoList struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"account_billing_objmt_infos" json:"account_billing_objmt_infos"`

	// Status metering collection request status
	Status string `xml:"status,omitempty" json:"status,omitempty"`

	// SizeUnit size unit of metric values
	SizeUnit string `xml:"size_unit,omitempty" json:"size_unit,omitempty"`

	// DateTime request time (ISO 8601 format (2020-01-27T14:30:55Z))
	DateTime string `xml:"date_time,omitempty" json:"date_time,omitempty"`

	// Info list of billing metrics objects for each user
	Info []AccountBillingInfo `xml:"account_billing_objmt_info" json:"account_billing_objmt_info"`
}

// AccountBillingInfo contains latest billing metrics info for one user
type AccountBillingInfo struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"account_billing_objmt_info" json:"account_billing_objmt_info"`

	// AccountID IAM user account ID used for request
	AccountID string `xml:"account_id,omitempty"`

	// ConsistentTime metrics collection UTC timestamp (ISO 8601 format (2007-04-05T14:30:55Z))
	ConsistentTime string `xml:"consistent_time,omitempty"`

	// TotalLocalData total local logical capacity usage(local: objects + MPU + user-written metadata)
	TotalLocalData int64 `xml:"total_local_data,omitempty"`

	// TotalReplicaData total replica logical capacity usage (replicated: objects + MPU + user written metadata)
	TotalReplicaData int64 `xml:"total_replica_data,omitempty"`

	// HardQuotaInCount hard quota in count
	HardQuotaInCount int64 `xml:"hard_quota_in_count,omitempty"`

	// HardQuotaInGB hard quota in logical size
	HardQuotaInGB int64 `xml:"hard_quota_in_GB,omitempty"`

	// SoftQuotaInCount soft quota in count
	SoftQuotaInCount int64 `xml:"soft_quota_in_count,omitempty"`

	// SoftQuotaInGB soft quota in logical size
	SoftQuotaInGB int64 `xml:"soft_quota_in_GB,omitempty"`

	// TotalUserObjectMetric total object logical and physical size, count per user
	TotalUserObjectMetric []StorageClassBasedCountSize `xml:"total_user_object_metric>storage_class_counts,omitempty"`

	// TotalMPUMetric total MPU parts logical and physical size, count per user
	TotalMPUMetric []StorageClassBasedCountSize `xml:"total_mpu_metric>storage_class_counts,omitempty"`

	// TotalMPRMetric total MPR parts logical and physical size, count per user
	TotalMPRMetric []StorageClassBasedCountSize `xml:"total_mpr_metric>storage_class_counts,omitempty"`

	// TotalReplicaObjectMetric total replicated objects logical and physical size, count per user
	TotalReplicaObjectMetric []StorageClassBasedCountSize `xml:"total_replica_object_metric>storage_class_counts,omitempty"`

	// TotalUserMetadataMetric total object metadata logical and physical size, count per user
	TotalUserMetadataMetric []StorageClassBasedCountSize `xml:"total_user_metadata_metric>storage_class_counts,omitempty"`

	// TotalReplicaMetadataMetric total replicated object metadata logical and physical size, count per user
	TotalReplicaMetadataMetric []StorageClassBasedCountSize `xml:"total_replica_metadata_metric>storage_class_counts,omitempty"`

	// BucketBillingInfo metrics for the buckets managed by this account
	BucketBillingInfo []BucketBillingInfo `xml:"bucket_billing_info,omitempty"`
}

// AccountBillingSampleList contains time range based billing metrics for users
type AccountBillingSampleList struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"account_billing_objmt_samples" json:"account_billing_objmt_samples"`

	// Status metering collection request status
	Status string `xml:"status,omitempty" json:"status,omitempty"`

	// SizeUnit size unit of metric values
	SizeUnit string `xml:"size_unit,omitempty" json:"size_unit,omitempty"`

	// DateTime request time (ISO 8601 format (2020-01-27T14:30:55Z))
	DateTime string `xml:"date_time,omitempty" json:"date_time,omitempty"`

	// StartTime is a start of time window for data selection (ISO 8601 format (2020-01-27T14:30:55Z))
	StartTime string `xml:"start_time,omitempty" json:"start_time,omitempty"`

	// EndTime is an end of time window for data selection (ISO 8601 format (2020-01-27T14:30:55Z))
	EndTime string `xml:"end_time,omitempty" json:"end_time,omitempty"`

	// Samples list of time range based billing metrics objects for each user
	Samples []AccountBillingSample `json:"account_billing_objmt_sample" xml:"account_billing_objmt_sample"`
}

// AccountBillingSample contains time range based billing metrics info for one user
type AccountBillingSample struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"account_billing_objmt_sample" json:"account_billing_objmt_sample"`

	// AccountID IAM user account ID used for request
	AccountID string `xml:"account_id,omitempty"`

	// StartTime is a start of time window for data selection (ISO 8601 format (2020-01-27T14:30:55Z))
	StartTime string `xml:"start_time,omitempty"`

	// EndTime is an end of time window for data selection (ISO 8601 format (2020-01-27T14:30:55Z))
	EndTime string `xml:"end_time,omitempty"`

	// SampleTimeRange time window in UTC format
	SampleTimeRange int64 `xml:"sample_time_range,omitempty"`

	// ConsistentTime metrics collection UTC timestamp in ISO 8601 format (2007-04-05T14:30:55Z)
	ConsistentTime string `xml:"consistent_time,omitempty"`

	// AccountBillingInfo contains latest billing metrics info for the user
	AccountBillingInfo AccountBillingInfo `xml:"account_billing_objmt_info,omitempty"`

	// UserCreationDelta list of objects creation delta per storage classes
	UserCreationDelta []StorageClassBasedCountSize `xml:"user_creation_delta>storage_class_counts"`

	// UserDeletionDelta list of objects deletion delta per storage classes
	UserDeletionDelta []StorageClassBasedCountSize `xml:"user_deletion_delta>storage_class_counts"`

	// MpuCreateDelta list of MPU parts creation delta per storage classes
	MpuCreateDelta []StorageClassBasedCountSize `xml:"mpu_create_delta>storage_class_counts"`

	// MpuDeleteDelta list of MPU parts deletion delta per storage classes
	MpuDeleteDelta []StorageClassBasedCountSize `xml:"mpu_delete_delta>storage_class_counts"`

	// MprCreateDelta list of MPR parts creation delta per storage classes
	MprCreateDelta []StorageClassBasedCountSize `xml:"mpr_create_delta>storage_class_counts"`

	// MprDeleteDelta list of MPR parts deletion delta per storage classes
	MprDeleteDelta []StorageClassBasedCountSize `xml:"mpr_delete_delta>storage_class_counts"`

	// ReplicaCreationDelta list of replicated objects creation delta per storage classes
	ReplicaCreationDelta []StorageClassBasedCountSize `xml:"replica_creation_delta>storage_class_counts"`

	// ReplicaDeletionDelta list of replicated objects deletion delta per storage classes
	ReplicaDeletionDelta []StorageClassBasedCountSize `xml:"replica_deletion_delta>storage_class_counts"`

	// UserMetadataCreationDelta list of object metadata creation delta per storage classes
	UserMetadataCreationDelta []StorageClassBasedCountSize `xml:"user_metadata_create_delta>storage_class_counts"`

	// UserMetadataDeletionDelta list of object metadata deletion delta per storage classes
	UserMetadataDeletionDelta []StorageClassBasedCountSize `xml:"user_metadata_delete_delta>storage_class_counts"`

	// ReplicaMetadataCreationDelta list of object metadata creation delta per storage classes
	ReplicaMetadataCreationDelta []StorageClassBasedCountSize `xml:"replica_metadata_create_delta>storage_class_counts"`

	// ReplicaMetadataDeletionDelta list of object metadata deletion delta per storage classes
	ReplicaMetadataDeletionDelta []StorageClassBasedCountSize `xml:"replica_metadata_delete_delta>storage_class_counts"`

	// HardQuotaInCount hard quota in count
	HardQuotaInCount int64 `xml:"hard_quota_in_count,omitempty"`

	// HardQuotaInGB hard quota in logical size
	HardQuotaInGB int64 `xml:"hard_quota_in_GB,omitempty"`

	// SoftQuotaInCount soft quota in count
	SoftQuotaInCount int64 `xml:"soft_quota_in_count,omitempty"`

	// SoftQuotaInGB soft quota in logical size
	SoftQuotaInGB int64 `xml:"soft_quota_in_GB,omitempty"`

	// BucketBillingSample metrics for the buckets managed by this account
	BucketBillingSample []BucketBillingSample `xml:"bucket_billing_sample,omitempty"`
}

// StorageClassBasedCountSize contains logical and physical size and count of objects per storage class
type StorageClassBasedCountSize struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `json:"storage_class_counts" xml:"storage_class_counts"`

	// StorageClass name of storage class provided metrics
	StorageClass string `xml:"storage_class,omitempty"`

	// Counts count of objects in requested metric
	Counts int64 `xml:"count_size>counts,omitempty"`

	// LogicalSize logical size of objects in requested metric
	LogicalSize int64 `xml:"count_size>logical_size,omitempty"`

	// CreateLogicalSize logical size of created objects using for metric data ingress
	CreateLogicalSize int64 `xml:"count_size>create_logical_size,omitempty"`

	// DeleteLogicalSize logical size of deleted objects using for metric data ingress
	DeleteLogicalSize int64 `xml:"count_size>delete_logical_size,omitempty"`

	// PhysicalSize physical size of objects in requested metric
	PhysicalSize int64 `xml:"count_size>physical_size,omitempty"`

	// CreatePhysicalSize physical size of created objects using for metric data ingress
	CreatePhysicalSize int64 `xml:"count_size>create_physical_size,omitempty"`

	// DeletePhysicalSize physical size of deleted objects using for metric data ingress
	DeletePhysicalSize int64 `xml:"count_size>delete_physical_size,omitempty"`
}

// BucketBillingInfoList list of the billing info per buckets
type BucketBillingInfoList struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"bucket_billing_objmt_infos" json:"bucket_billing_objmt_infos"`

	// Status metering collection request status
	Status string `xml:"status,omitempty" json:"status,omitempty"`

	// SizeUnit size unit of metric values
	SizeUnit string `xml:"size_unit,omitempty" json:"size_unit,omitempty"`

	// DateTime request time (ISO 8601 format (2020-01-27T14:30:55Z))
	DateTime string `xml:"date_time,omitempty" json:"date_time,omitempty"`

	// Info list of billing metrics objects for each bucket
	Info []BucketBillingInfo `xml:"bucket_billing_objmt_info" json:"bucket_billing_objmt_info"`
}

// BucketBillingInfo contains latest billing metrics info for one bucket
type BucketBillingInfo struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"bucket_billing_objmt_info" json:"bucket_billing_objmt_info"`

	// BucketName bucket name using for collecting metrics
	BucketName string `xml:"bucket_name,omitempty"`

	// CompressionRatio float value of compression ratio in bucket
	CompressionRatio float64 `xml:"compression_ratio,omitempty"`

	// ConsistentTime metrics collection UTC timestamp in ISO 8601 format (2007-04-05T14:30:55Z)
	ConsistentTime string `xml:"consistent_time,omitempty"`

	// HardQuotaInCount hard quota in count
	HardQuotaInCount int64 `xml:"hard_quota_in_count,omitempty"`

	// HardQuotaInGB hard quota in logical size
	HardQuotaInGB int64 `xml:"hard_quota_in_GB,omitempty"`

	// SoftQuotaInCount soft quota in count
	SoftQuotaInCount int64 `xml:"soft_quota_in_count,omitempty"`

	// SoftQuotaInGB soft quota in logical size
	SoftQuotaInGB int64 `xml:"soft_quota_in_GB,omitempty"`

	// ObjectDistribution metrics collection of objects grouped by size
	ObjectDistribution string `xml:"object_distribution,omitempty"`

	// TotalLocalData total local logical capacity usage(local: objects + MPU + user-written metadata)
	TotalLocalData int64 `xml:"total_local_data,omitempty"`

	// TotalReplicaData total replica logical capacity usage (replicated: objects + MPU + user written metadata)
	TotalReplicaData int64 `xml:"total_replica_data,omitempty"`

	// TotalUserObjectMetric total object logical and physical size, count per bucket
	TotalUserObjectMetric []StorageClassBasedCountSize `xml:"total_user_object_metric>storage_class_counts"`

	// TotalUserMetadataMetric total user metadata logical and physical size, count per bucket
	TotalUserMetadataMetric []StorageClassBasedCountSize `xml:"total_user_metadata_metric>storage_class_counts"`

	// TotalMPUMetric total MPU parts logical and physical size, count per bucket
	TotalMPUMetric []StorageClassBasedCountSize `xml:"total_mpu_metric>storage_class_counts"`

	// TotalMPRMetric total MPR parts logical and physical size, count per bucket
	TotalMPRMetric []StorageClassBasedCountSize `xml:"total_mpr_metric>storage_class_counts"`

	// TotalReplicaObjectMetric total replicated objects logical and physical size, count per bucket
	TotalReplicaObjectMetric []StorageClassBasedCountSize `xml:"total_replica_object_metric>storage_class_counts"`

	// TotalReplicaMetadataMetric total replicated object metadata logical and physical size, count per bucket
	TotalReplicaMetadataMetric []StorageClassBasedCountSize `xml:"total_replica_metadata_metric>storage_class_counts"`
}

// BucketBillingSampleList contains time range based billing metrics for buckets
type BucketBillingSampleList struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"bucket_billing_objmt_samples" json:"bucket_billing_objmt_samples"`

	// Status metering collection request status
	Status string `xml:"status,omitempty" json:"status,omitempty"`

	// SizeUnit size unit of metric values
	SizeUnit string `xml:"size_unit,omitempty" json:"size_unit,omitempty"`

	// DateTime request time (ISO 8601 format (2020-01-27T14:30:55Z))
	DateTime string `xml:"date_time,omitempty" json:"date_time,omitempty"`

	// StartTime is a start of time window for data selection (ISO 8601 format (2020-01-27T14:30:55Z))
	StartTime string `xml:"start_time,omitempty" json:"start_time,omitempty"`

	// EndTime is an end of time window for data selection (ISO 8601 format (2020-01-27T14:30:55Z))
	EndTime string `xml:"end_time,omitempty" json:"end_time,omitempty"`

	// Samples list of time range based billing metrics objects for each bucket
	Samples []BucketBillingSample `json:"bucket_billing_objmt_sample" xml:"bucket_billing_objmt_sample"`
}

// BucketBillingSample contains time range based billing metrics info for one bucket
type BucketBillingSample struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"bucket_billing_objmt_sample" json:"bucket_billing_objmt_sample"`

	// BucketName bucket name using for collecting metrics
	BucketName string `xml:"bucket_name,omitempty"`

	// SampleTimeRange time window in UTC format
	SampleTimeRange int64 `xml:"sample_time_range,omitempty"`

	// CrrThroughput Calculated by replicated object logical size/time range in seconds
	CrrThroughput int64 `xml:"crr_throughput,omitempty"`

	// HardQuotaInCount hard quota in count
	HardQuotaInCount int64 `xml:"hard_quota_in_count,omitempty"`

	// HardQuotaInGB hard quota in logical size
	HardQuotaInGB int64 `xml:"hard_quota_in_GB,omitempty"`

	// SoftQuotaInCount soft quota in count
	SoftQuotaInCount int64 `xml:"soft_quota_in_count,omitempty"`

	// SoftQuotaInGB soft quota in logical size
	SoftQuotaInGB int64 `xml:"soft_quota_in_GB,omitempty"`

	// ConsistentTime metrics collection UTC timestamp in ISO 8601 format (2007-04-05T14:30:55Z)
	ConsistentTime string `xml:"consistent_time,omitempty"`

	// BucketBillingInfo contains latest billing metrics info for the bucket
	BucketBillingInfo BucketBillingInfo `xml:"bucket_billing_objmt_info,omitempty"`

	// BucketBillingTags list of ingress nad egress billing values
	BucketBillingTags []BucketBillingTag `xml:"bucket_billing_tag,omitempty"`

	// UserCreationDelta list of objects creation delta per storage classes
	UserCreationDelta []StorageClassBasedCountSize `xml:"user_creation_delta>storage_class_counts"`

	// UserDeletionDelta list of objects deletion delta per storage classes
	UserDeletionDelta []StorageClassBasedCountSize `xml:"user_deletion_delta>storage_class_counts"`

	// MpuCreateDelta list of MPU parts creation delta per storage classes
	MpuCreateDelta []StorageClassBasedCountSize `xml:"mpu_create_delta>storage_class_counts"`

	// MpuDeleteDelta list of MPU parts deletion delta per storage classes
	MpuDeleteDelta []StorageClassBasedCountSize `xml:"mpu_delete_delta>storage_class_counts"`

	// MprCreateDelta list of MPU parts creation delta per storage classes
	MprCreateDelta []StorageClassBasedCountSize `xml:"mpr_create_delta>storage_class_counts"`

	// MprDeleteDelta list of MPU parts deletion delta per storage classes
	MprDeleteDelta []StorageClassBasedCountSize `xml:"mpr_delete_delta>storage_class_counts"`

	// ReplicaCreationDelta list of replicated objects creation delta per storage classes
	ReplicaCreationDelta []StorageClassBasedCountSize `xml:"replica_creation_delta>storage_class_counts"`

	// ReplicaDeletionDelta list of replicated objects deletion delta per storage classes
	ReplicaDeletionDelta []StorageClassBasedCountSize `xml:"replica_deletion_delta>storage_class_counts"`

	// UserMetadataCreationDelta list of object metadata creation delta per storage classes
	UserMetadataCreationDelta []StorageClassBasedCountSize `xml:"user_metadata_create_delta>storage_class_counts"`

	// UserMetadataDeletionDelta list of object metadata deletion delta per storage classes
	UserMetadataDeletionDelta []StorageClassBasedCountSize `xml:"user_metadata_delete_delta>storage_class_counts"`

	// ReplicaMetadataCreationDelta list of object metadata creation delta per storage classes
	ReplicaMetadataCreationDelta []StorageClassBasedCountSize `xml:"replica_metadata_create_delta>storage_class_counts"`

	// ReplicaMetadataDeletionDelta list of object metadata deletion delta per storage classes
	ReplicaMetadataDeletionDelta []StorageClassBasedCountSize `xml:"replica_metadata_delete_delta>storage_class_counts"`
}

// BucketBillingTag contains ingress and egress metrics for bucket
type BucketBillingTag struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"bucket_billing_tag" json:"bucket_billing_tag"`

	// BucketName bucket name using for collecting metrics
	BucketName string `xml:"bucket_name,omitempty"`

	// Ingress size of ingress data in bucket
	Ingress int64 `xml:"ingress,omitempty"`

	// Egress size of ergess data in bucket
	Egress int64 `xml:"egress,omitempty"`
}

// BucketPerfDataList time range based list of performance bucket billing metrics
type BucketPerfDataList struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"bucket_perf_samples" json:"bucket_perf_samples"`

	// Status metering collection request status
	Status string `xml:"status,omitempty" json:"status,omitempty"`

	// SizeUnit size unit of metric values
	SizeUnit string `xml:"size_unit,omitempty" json:"size_unit,omitempty"`

	// DateTime request time (ISO 8601 format (2020-01-27T14:30:55Z))
	DateTime string `xml:"date_time,omitempty" json:"date_time,omitempty"`

	// StartTime is a start of time window for data selection (ISO 8601 format (2020-01-27T14:30:55Z))
	StartTime string `xml:"start_time,omitempty" json:"start_time,omitempty"`

	// EndTime is an end of time window for data selection (ISO 8601 format (2020-01-27T14:30:55Z))
	EndTime string `xml:"end_time,omitempty" json:"end_time,omitempty"`

	// Sample list of metrics per each bucket
	Samples []BucketPerfSample `xml:"bucket_perf_sample" json:"bucket_perf_sample" `
}

// BucketPerfSample time range based bucket performance billing metrics
type BucketPerfSample struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"bucket_perf_sample" json:"bucket_perf_sample"`

	// SampleTimeRange time window in UTC format BucketName      string `xml:"bucket_name,omitempty"`
	SampleTimeRange int64 `xml:"sample_time_range,omitempty"`

	// ConsistentTime metrics collection UTC timestamp in ISO 8601 format (2007-04-05T14:30:55Z)
	ConsistentTime string `xml:"consistent_time,omitempty"`

	// IngressLatency ingress latency value for bucket in requested time window
	IngressLatency int64 `xml:"ingress_latency,omitempty"`

	// IngressBytes ingress bytes for bucket in requested time window
	IngressBytes int64 `xml:"ingress_bytes,omitempty"`

	// IngressCounts ingress objects count for bucket in requested time window
	IngressCounts int64 `xml:"ingress_counts,omitempty"`

	// EgressLatecy egress latency for bucket in requested time window
	EgressLatency int64 `xml:"egress_latency,omitempty"`

	// EgressBytes egress bytes for bucket in requested time window
	EgressBytes int64 `xml:"egress_bytes,omitempty"`

	// EgressCounts egress objects count for bucket in requested time window
	EgressCounts int64 `xml:"egress_counts,omitempty"`
}

// BucketReplicationSampleList time range based list of bucket replication billing metrics
type BucketReplicationSampleList struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"bucket_replication_samples" json:"bucket_replication_samples"`

	// Status metering collection request status
	Status string `xml:"status,omitempty" json:"status,omitempty"`

	// SizeUnit size unit of metric values
	SizeUnit string `xml:"size_unit,omitempty" json:"size_unit,omitempty"`

	// DateTime request time (ISO 8601 format (2020-01-27T14:30:55Z))
	DateTime string `xml:"date_time,omitempty" json:"date_time,omitempty"`

	// StartTime is a start of time window for data selection (ISO 8601 format (2020-01-27T14:30:55Z))
	StartTime string `xml:"start_time,omitempty" json:"start_time,omitempty"`

	// EndTime is an end of time window for data selection (ISO 8601 format (2020-01-27T14:30:55Z))
	EndTime string `xml:"end_time,omitempty" json:"end_time,omitempty"`

	// Samples list of replication metrics per each bucket
	Samples []BucketReplicationSample `xml:"bucket_replication_sample,omitempty" json:"bucket_replication_sample,omitempty"`
}

// BucketReplicationSample time range based replication billing metrics for one bucket
type BucketReplicationSample struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"bucket_replication_sample" json:"bucket_replication_sample"`

	// SourceBucket name of replication source bucket
	SourceBucket string `xml:"replication_source_destination>source_bucket,omitempty"`

	// DestinationBucket name of replication destination bucket
	DestinationBucket string `xml:"replication_source_destination>destination_bucket_arn,omitempty"`

	// SampleTimeRange time window in UTC format
	SampleTimeRange int64 `xml:"sample_time_range,omitempty"`

	// ConsistentTime metrics collection UTC timestamp (ISO 8601 format (2020-01-27T14:30:55Z))
	ConsistentTime string `xml:"consistent_time,omitempty"`

	ReplicationBillingInfo ReplicationBillingInfo `xml:"replication_billing_info,omitempty"`

	// PendingToReplicateDelta metrics of data pending replication per storage class
	PendingToReplicateDelta []StorageClassBasedCountSize `xml:"pending_to_replicate_delta>storage_class_counts,omitempty"`

	// ReplicatedDelta replicated data metrics per storage class
	ReplicatedDelta []StorageClassBasedCountSize `xml:"replicated_delta>storage_class_counts,omitempty"`

	// PendingToReplicateDelta metrics of data failed to replicate per storage class
	ReplicatedFailureDelta []StorageClassBasedCountSize `xml:"replicated_failure_delta>storage_class_counts,omitempty"`
}

// BucketReplicationInfoList list of latest bucket replication billing metrics
type BucketReplicationInfoList struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"replication_info_list" json:"replication_info_list"`

	// Status metering collection request status
	Status string `xml:"status,omitempty" json:"status,omitempty"`

	// SizeUnit size unit of metric values
	SizeUnit string `xml:"size_unit,omitempty" json:"size_unit,omitempty"`

	// DateTime request time (ISO 8601 format (2020-01-27T14:30:55Z))
	DateTime string `xml:"date_time,omitempty" json:"date_time,omitempty"`

	// Info list of replication metrics per each bucket
	Info []ReplicationBillingInfo `xml:"replication_billing_info" json:"replication_billing_info"`
}

// ReplicationBillingInfo latest replication billing metrics for one bucket
type ReplicationBillingInfo struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"replication_billing_info" json:"replication_billing_info"`

	// SourceBucket name of replication source bucket
	SourceBucket string `xml:"replication_source_destination>source_bucket,omitempty"`

	// DestinationBucket name of replication destination bucket
	DestinationBucket string `xml:"replication_source_destination>destination_bucket_arn,omitempty"`

	// ConsistentTime metrics collection UTC timestamp (ISO 8601 format (2020-01-27T14:30:55Z))
	ConsistentTime string `xml:"consistent_time,omitempty"`

	// PendingToReplicate metrics of data pending replication per storage class
	PendingToReplicate []StorageClassBasedCountSize `xml:"pending_to_replicate>storage_class_counts,omitempty"`
}

// StoreBillingInfoList list of object store latest billing info metrics
type StoreBillingInfoList struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"store_billing_info_list" json:"store_billing_info_list"`

	// Status metering collection request status
	Status string `xml:"status,omitempty" json:"status,omitempty"`

	// SizeUnit size unit of metric values
	SizeUnit string `xml:"size_unit,omitempty" json:"size_unit,omitempty"`

	// DateTime request time (ISO 8601 format (2020-01-27T14:30:55Z))
	DateTime string `xml:"date_time,omitempty" json:"date_time,omitempty"`

	// StoreBillingInfo object store latest billing info metrics
	Info StoreBillingInfo `xml:"store_billing_info,omitempty" json:"store_billing_info,omitempty"`
}

// StoreBillingInfo object store latest billing info metrics
type StoreBillingInfo struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"store_billing_info" json:"store_billing_info"`

	// CompressionRatio float value of compression ratio in object store
	CompressionRatio float64 `xml:"compression_ratio,omitempty"`

	// ConsistentTime metrics collection UTC timestamp (ISO 8601 format (2020-01-27T14:30:55Z))
	ConsistentTime string `xml:"consistent_time,omitempty"`

	// TotalLocalData total local logical capacity usage(local: objects + MPU + user-written metadata)
	TotalLocalData int64 `xml:"total_local_data,omitempty"`

	// TotalReplicaData total replica logical capacity usage (replicated: objects + MPU + user written metadata)
	TotalReplicaData int64 `xml:"total_replica_data,omitempty"`

	// TotalUserObjectMetric total object logical and physical size, count in object store
	TotalUserObjectMetric []StorageClassBasedCountSize `xml:"total_user_object_metric>storage_class_counts,omitempty"`

	// TotalMPUMetric total MPU parts logical and physical size, count in object store
	TotalMPUMetric []StorageClassBasedCountSize `xml:"total_mpu_metric>storage_class_counts,omitempty"`

	// TotalMPRMetric total MPR parts logical and physical size, count per user
	TotalMPRMetric []StorageClassBasedCountSize `xml:"total_mpr_metric>storage_class_counts,omitempty"`

	// TotalUserMetadataMetric total object metadata logical and physical size, count per user
	TotalUserMetadataMetric []StorageClassBasedCountSize `xml:"total_user_metadata_metric>storage_class_counts,omitempty"`

	// TotalReplicaMetadataMetric total replicated object metadata logical and physical size, count per user
	TotalReplicaMetadataMetric []StorageClassBasedCountSize `xml:"total_replica_metadata_metric>storage_class_counts,omitempty"`

	// TotalReplicaObjectMetric total replicated objects logical and physical size, count in object store
	TotalReplicaObjectMetric []StorageClassBasedCountSize `xml:"total_replica_object_metric>storage_class_counts,omitempty"`

	// TopBucketsByObjectCount list of top buckets by objects count in object store
	TopBucketsByObjectCount []TopNBucket `xml:"top_n_buckets_by_object_count>top_n_bucket"`

	// TopBucketsByObjectSize list of top buckets by objects size in object store
	TopBucketsByObjectSize []TopNBucket `xml:"top_n_buckets_by_object_size>top_n_bucket"`
}

// StoreBillingSampleList time range based list of billing metrics in object store
type StoreBillingSampleList struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"store_billing_samples" json:"store_billing_samples"`

	// Status metering collection request status
	Status string `xml:"status,omitempty" json:"status,omitempty"`

	// SizeUnit size unit of metric values
	SizeUnit string `xml:"size_unit,omitempty" json:"size_unit,omitempty"`

	// DateTime request time (ISO 8601 format (2020-01-27T14:30:55Z))
	DateTime string `xml:"date_time,omitempty" json:"date_time,omitempty"`

	// StartTime is a start of time window for data selection (ISO 8601 format (2020-01-27T14:30:55Z))
	StartTime string `xml:"start_time,omitempty" json:"start_time,omitempty"`

	// EndTime is an end of time window for data selection (ISO 8601 format (2020-01-27T14:30:55Z))
	EndTime string `xml:"end_time,omitempty" json:"end_time,omitempty"`

	// Samples list of object store billing metrics
	Samples []StoreBillingSample `xml:"store_billing_sample,omitempty" json:"store_billing_sample,omitempty"`
}

// StoreBillingSample time range based billing metrics in object store
type StoreBillingSample struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"store_billing_sample" json:"store_billing_sample"`

	// SampleTimeRange time window in UTC format
	SampleTimeRange int64 `xml:"sample_time_range,omitempty"`

	// ConsistentTime metrics collection UTC timestamp (ISO 8601 format (2020-01-27T14:30:55Z))
	ConsistentTime string `xml:"consistent_time,omitempty"`

	// Info object store latest billing info metrics
	Info StoreBillingInfo `xml:"store_billing_info,omitempty" json:"store_billing_info,omitempty"`

	// UserCreationDelta list of objects creation delta per storage classes
	UserCreationDelta []StorageClassBasedCountSize `xml:"user_creation_delta>storage_class_counts"`

	// UserDeletionDelta list of objects deletion delta per storage classes
	UserDeletionDelta []StorageClassBasedCountSize `xml:"user_deletion_delta>storage_class_counts"`

	// MpuCreateDelta list of MPU parts creation delta per storage classes
	MpuCreateDelta []StorageClassBasedCountSize `xml:"mpu_create_delta>storage_class_counts"`

	// MpuDeleteDelta list of MPU parts deletion delta per storage classes
	MpuDeleteDelta []StorageClassBasedCountSize `xml:"mpu_delete_delta>storage_class_counts"`

	// MprCreateDelta list of MPU parts creation delta per storage classes
	MprCreateDelta []StorageClassBasedCountSize `xml:"mpr_create_delta>storage_class_counts"`

	// MprDeleteDelta list of MPU parts deletion delta per storage classes
	MprDeleteDelta []StorageClassBasedCountSize `xml:"mpr_delete_delta>storage_class_counts"`

	// ReplicaCreationDelta list of replicated objects creation delta per storage classes
	ReplicaCreationDelta []StorageClassBasedCountSize `xml:"replica_creation_delta>storage_class_counts"`

	// ReplicaDeletionDelta list of replicated objects deletion delta per storage classes
	ReplicaDeletionDelta []StorageClassBasedCountSize `xml:"replica_deletion_delta>storage_class_counts"`

	// UserMetadataCreationDelta list of object metadata creation delta per storage classes
	UserMetadataCreationDelta []StorageClassBasedCountSize `xml:"user_metadata_create_delta>storage_class_counts"`

	// UserMetadataDeletionDelta list of object metadata deletion delta per storage classes
	UserMetadataDeletionDelta []StorageClassBasedCountSize `xml:"user_metadata_delete_delta>storage_class_counts"`

	// ReplicaMetadataCreationDelta list of object metadata creation delta per storage classes
	ReplicaMetadataCreationDelta []StorageClassBasedCountSize `xml:"replica_metadata_create_delta>storage_class_counts"`

	// ReplicaMetadataDeletionDelta list of object metadata deletion delta per storage classes
	ReplicaMetadataDeletionDelta []StorageClassBasedCountSize `xml:"replica_metadata_delete_delta>storage_class_counts"`
}

// TopNBucket top bucket metric
type TopNBucket struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"top_n_bucket" json:"top_n_bucket"`

	// BucketName top bucket name
	BucketName string `xml:"bucket_name,omitempty" json:"bucket_name,omitempty"`

	// MetricNumber bucket metric value
	MetricNumber int64 `xml:"metric_number,omitempty" json:"metric_number,omitempty"`
}

// StoreReplicationDataList time range based list of replication data metrics in object stores
type StoreReplicationDataList struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"store_replication_list" json:"store_replication_list"`

	// Status metering collection request status
	Status string `xml:"status,omitempty" json:"status,omitempty"`

	// SizeUnit size unit of metric values
	SizeUnit string `xml:"size_unit,omitempty" json:"size_unit,omitempty"`

	// DateTime request time (ISO 8601 format (2020-01-27T14:30:55Z))
	DateTime string `xml:"date_time,omitempty" json:"date_time,omitempty"`

	// StartTime is a start of time window for data selection (ISO 8601 format (2020-01-27T14:30:55Z))
	StartTime string `xml:"start_time,omitempty" json:"start_time,omitempty"`

	// EndTime is an end of time window for data selection (ISO 8601 format (2020-01-27T14:30:55Z))
	EndTime string `xml:"end_time,omitempty" json:"end_time,omitempty"`

	// Samples list of replication data metrics in object stores
	Samples []StoreReplicationThroughputRto `xml:"store_replication_throughput_rto,omitempty" json:"store_replication_throughput_rto,omitempty"`
}

// StoreReplicationThroughputRto replication data metrics per one object store
type StoreReplicationThroughputRto struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"store_replication_throughput_rto" json:"store_replication_throughput_rto"`

	// SampleTimeRange time window in UTC format
	SampleTimeRange int64 `xml:"sample_time_range,omitempty"`

	// ConsistentTime metrics collection UTC timestamp (ISO 8601 format (2020-01-27T14:30:55Z))
	ConsistentTime string `xml:"consistent_time,omitempty"`

	// DestinationStore object store URN identifier
	DestinationStore string `xml:"destination_store,omitempty" json:"destination_store,omitempty"`

	// Throughput CRR throughput value in object store
	Throughput int64 `xml:"throughput,omitempty" json:"throughput,omitempty"`

	// RTO CRR RTO value in object store
	RTO int64 `xml:"rto,omitempty" json:"rto,omitempty"`

	// PendingToReplicate metrics of data pending to replicate in object store per storage classes
	PendingToReplicate []StorageClassBasedCountSize `xml:"pending_to_replicate>storage_class_counts,omitempty"`

	// ReplicatedDelta metrics of replicated data in object store per storage classes
	ReplicatedDelta []StorageClassBasedCountSize `xml:"replicated_delta>storage_class_counts,omitempty"`

	// ReplicatedFailedDelta metrics of failed replicated data in object store per storage classes
	ReplicatedFailedDelta []StorageClassBasedCountSize `xml:"replicate_failed_delta>storage_class_counts,omitempty"`
}
