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

package fake

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dell/goobjectscale/pkg/client/api"
	"github.com/dell/goobjectscale/pkg/client/model"
)

// ClientSet is a set of clients for each API section
type ClientSet struct {
	buckets               api.BucketsInterface
	objectUser            api.ObjectUserInterface
	tenants               api.TenantsInterface
	objectMt              api.ObjmtInterface
	crr                   api.CRRInterface
	alertPolicies         api.AlertPoliciesInterface
	status                api.StatusInterfaces
	federatedobjectstores api.FederatedObjectStoresInterface
}

// NewClientSet returns a new client set based on the provided REST client parameters
func NewClientSet(objs ...interface{}) *ClientSet {
	var (
		policy                      = make(map[string]string)
		bucketList                  []model.Bucket
		blobUsers                   []model.BlobUser
		userSecrets                 []UserSecret
		userInfoList                []UserInfo
		tenantList                  []model.Tenant
		accountBillingInfoList      *model.AccountBillingInfoList
		accountBillingSampleList    *model.AccountBillingSampleList
		bucketBillingInfoList       *model.BucketBillingInfoList
		bucketBillingSampleList     *model.BucketBillingSampleList
		bucketBillingPerfList       *model.BucketPerfDataList
		bucketReplicationInfoList   *model.BucketReplicationInfoList
		bucketReplicationSampleList *model.BucketReplicationSampleList
		storeBillingInfoList        *model.StoreBillingInfoList
		storeBillingSampleList      *model.StoreBillingSampleList
		storeReplicationDataList    *model.StoreReplicationDataList
		crr                         *model.CRR
		alertPolicies               []model.AlertPolicy
		rebuildInfo                 *model.RebuildInfo
		federatedObjectStoreList    []model.FederatedObjectStore
	)
	for _, o := range objs {
		switch object := o.(type) {
		case *model.Bucket:
			bucketList = append(bucketList, *object)
		case *BucketPolicy:
			policy[fmt.Sprintf("%s/%s", object.BucketName, object.Namespace)] = object.Policy
		case *model.BlobUser:
			blobUsers = append(blobUsers, *object)
		case *UserSecret:
			userSecrets = append(userSecrets, *object)
		case *UserInfo:
			userInfoList = append(userInfoList, *object)
		case *model.Tenant:
			tenantList = append(tenantList, *object)
		case *model.AccountBillingInfoList:
			accountBillingInfoList = object
		case *model.AccountBillingSampleList:
			accountBillingSampleList = object
		case *model.BucketBillingInfoList:
			bucketBillingInfoList = object
		case *model.BucketBillingSampleList:
			bucketBillingSampleList = object
		case *model.BucketPerfDataList:
			bucketBillingPerfList = object
		case *model.BucketReplicationInfoList:
			bucketReplicationInfoList = object
		case *model.BucketReplicationSampleList:
			bucketReplicationSampleList = object
		case *model.StoreBillingInfoList:
			storeBillingInfoList = object
		case *model.StoreBillingSampleList:
			storeBillingSampleList = object
		case *model.StoreReplicationDataList:
			storeReplicationDataList = object
		case *model.CRR:
			crr = object
		case *model.AlertPolicy:
			alertPolicies = append(alertPolicies, *object)
		case *model.RebuildInfo:
			rebuildInfo = object
		case *model.FederatedObjectStore:
			federatedObjectStoreList = append(federatedObjectStoreList, *object)
		default:
			panic(fmt.Sprintf("Fake client set doesn't support %T type", o))
		}
	}

	return &ClientSet{
		buckets: &Buckets{
			items:  bucketList,
			policy: policy,
		},
		objectUser: NewObjectUsers(blobUsers, userSecrets, userInfoList),
		tenants: &Tenants{
			items: tenantList,
		},
		objectMt: &Objmt{
			accountBillingInfoList:      accountBillingInfoList,
			accountBillingSampleList:    accountBillingSampleList,
			bucketBillingInfoList:       bucketBillingInfoList,
			bucketBillingSampleList:     bucketBillingSampleList,
			bucketBillingPerfList:       bucketBillingPerfList,
			bucketReplicationInfoList:   bucketReplicationInfoList,
			bucketReplicationSampleList: bucketReplicationSampleList,
			storeBillingInfoList:        storeBillingInfoList,
			storeBillingSampleList:      storeBillingSampleList,
			storeReplicationDataList:    storeReplicationDataList,
		},
		crr: &CRR{
			Config: crr,
		},

		alertPolicies: &AlertPolicies{
			items: alertPolicies,
		},
		status: &Status{
			RebuildInfo: rebuildInfo,
		},
		federatedobjectstores: &FederatedObjectStores{
			items: federatedObjectStoreList,
		},
	}
}

// Status implements the client API
func (c *ClientSet) Status() api.StatusInterfaces {
	return c.status
}

// CRR implements the client API
func (c *ClientSet) CRR() api.CRRInterface {
	return c.crr
}

// AlertPolicies implements the client API
func (c *ClientSet) AlertPolicies() api.AlertPoliciesInterface {
	return c.alertPolicies
}

// Buckets implements the client API
func (c *ClientSet) Buckets() api.BucketsInterface {
	return c.buckets
}

// FederatedObjectStores implements the client API.
func (c *ClientSet) FederatedObjectStores() api.FederatedObjectStoresInterface {
	return c.federatedobjectstores
}

// Tenants implements the client API.
func (c *ClientSet) Tenants() api.TenantsInterface {
	return c.tenants
}

// ObjectUser implements the client API.
func (c *ClientSet) ObjectUser() api.ObjectUserInterface {
	return c.objectUser
}

// ObjectMt implements the client API for objMT metrics.
func (c *ClientSet) ObjectMt() api.ObjmtInterface {
	return c.objectMt
}

// BucketPolicy contains information about bucket policy to be used in fake client set.
type BucketPolicy struct {
	BucketName string
	Policy     string
	Namespace  string
}

// UserSecret make easiest passing secret about users.
type UserSecret struct {
	UID    string
	Secret *model.ObjectUserSecret
}

// UserInfo make easiest passing info about users.
type UserInfo struct {
	UID  string
	Info *model.ObjectUserInfo
}

// ObjectUsers contains information about object users to be used in fake client set.
type ObjectUsers struct {
	Users    *model.ObjectUserList
	Secrets  map[string]*model.ObjectUserSecret
	InfoList map[string]*model.ObjectUserInfo
}

var _ api.ObjectUserInterface = (*ObjectUsers)(nil) // interface guard

// NewObjectUsers returns initialized ObjectUsers.
func NewObjectUsers(blobUsers []model.BlobUser, userSecrets []UserSecret, userInfoList []UserInfo) *ObjectUsers {
	mappedUserSecrets := map[string]*model.ObjectUserSecret{}
	mappedUserInfoList := map[string]*model.ObjectUserInfo{}
	for _, s := range userSecrets {
		mappedUserSecrets[s.UID] = s.Secret
	}
	for _, i := range userInfoList {
		mappedUserInfoList[i.UID] = i.Info
	}
	return &ObjectUsers{
		&model.ObjectUserList{
			BlobUser: blobUsers,
		},
		mappedUserSecrets,
		mappedUserInfoList,
	}
}

// List returns a list of object users.
func (o *ObjectUsers) List(_ map[string]string) (*model.ObjectUserList, error) {
	return o.Users, nil
}

// GetSecret returns information about object user secrets.
func (o *ObjectUsers) GetSecret(uid string, _ map[string]string) (*model.ObjectUserSecret, error) {
	if _, ok := o.Secrets[uid]; !ok {
		return nil, model.Error{
			Description: "secret not found",
			Details:     fmt.Sprintf("secret for %s is not found", uid),
			Code:        model.CodeResourceNotFound,
		}
	}
	return o.Secrets[uid], nil
}

// CreateSecret will create a specific secret
func (o *ObjectUsers) CreateSecret(uid string, req model.ObjectUserSecretKeyCreateReq, _ map[string]string) (*model.ObjectUserSecretKeyCreateRes, error) {
	if _, ok := o.Secrets[uid]; !ok {
		o.Secrets[uid] = &model.ObjectUserSecret{
			SecretKey1: req.SecretKey,
		}
		return &model.ObjectUserSecretKeyCreateRes{
			SecretKey: req.SecretKey,
		}, nil
	}

	switch {
	case o.Secrets[uid].SecretKey1 != "" && o.Secrets[uid].SecretKey2 != "":
		return nil, model.Error{
			Description: "max keys reached",
			Details:     fmt.Sprintf("user %s already has 2 valid keys", uid),
			Code:        model.CodeExceedingLimit,
		}
	case o.Secrets[uid].SecretKey1 != "":
		o.Secrets[uid].SecretKey2 = req.SecretKey
		return &model.ObjectUserSecretKeyCreateRes{
			SecretKey: req.SecretKey,
		}, nil
	default:
		o.Secrets[uid].SecretKey1 = req.SecretKey
		return &model.ObjectUserSecretKeyCreateRes{
			SecretKey: req.SecretKey,
		}, nil
	}
}

// DeleteSecret will delete a specific secret
func (o *ObjectUsers) DeleteSecret(uid string, req model.ObjectUserSecretKeyDeleteReq, _ map[string]string) error {
	if _, ok := o.Secrets[uid]; !ok {
		return model.Error{
			Description: "user not found",
			Details:     fmt.Sprintf("user %s not found", uid),
			Code:        model.CodeResourceNotFound,
		}
	}
	switch req.SecretKey {
	case o.Secrets[uid].SecretKey1:
		o.Secrets[uid].SecretKey1 = o.Secrets[uid].SecretKey2
		clearSecretKey2(o.Secrets[uid])
		return nil
	case o.Secrets[uid].SecretKey2:
		clearSecretKey2(o.Secrets[uid])
		return nil
	default:
		return model.Error{
			Description: "not found",
			Details:     fmt.Sprintf("user %s secret key not found", uid),
			Code:        model.CodeResourceNotFound,
		}
	}
}

func clearSecretKey2(key *model.ObjectUserSecret) {
	key.SecretKey2 = ""
	key.KeyExpiryTimestamp2 = ""
	key.KeyExpiryTimestamp2 = ""
}

// GetInfo returns information about object user.
func (o *ObjectUsers) GetInfo(uid string, _ map[string]string) (*model.ObjectUserInfo, error) {
	if _, ok := o.InfoList[uid]; !ok {
		return nil, model.Error{
			Description: "info not found",
			Details:     fmt.Sprintf("info for %s is not found", uid),
			Code:        model.CodeResourceNotFound,
		}
	}
	return o.InfoList[uid], nil
}

// FederatedObjectStores implements the federated object stores API
type FederatedObjectStores struct {
	items []model.FederatedObjectStore
}

var _ api.FederatedObjectStoresInterface = (*FederatedObjectStores)(nil) // interface guard

// List implements the tenants API
func (f *FederatedObjectStores) List(_ map[string]string) (*model.FederatedObjectStoreList, error) {
	return &model.FederatedObjectStoreList{Items: f.items}, nil
}

// Tenants implements the tenants API
type Tenants struct {
	items []model.Tenant
}

var _ api.TenantsInterface = (*Tenants)(nil) // interface guard

// Create creates a tenant in an object store.
func (t *Tenants) Create(payload model.TenantCreate) (*model.Tenant, error) {
	newtenant := &model.Tenant{
		ID:                payload.AccountID,
		EncryptionEnabled: payload.EncryptionEnabled,
		ComplianceEnabled: payload.ComplianceEnabled,
		BucketBlockSize:   payload.BucketBlockSize,
	}
	t.items = append(t.items, *newtenant)
	return newtenant, nil
}

// Delete deletes a tenant in an object store.
func (t *Tenants) Delete(tenantID string) error {
	for i, tenant := range t.items {
		if tenant.ID == tenantID {
			t.items = append(t.items[:i], t.items[i+1:]...)
			return nil
		}
	}
	return model.Error{
		Description: "tenant not found",
		Code:        model.CodeResourceNotFound,
	}
}

// Update updates Tenant details default_bucket_size and alias
func (t *Tenants) Update(payload model.TenantUpdate, tenantID string) error {
	for i, tenant := range t.items {
		if tenant.ID == tenantID {
			t.items[i].BucketBlockSize = payload.BucketBlockSize
			t.items[i].Alias = payload.Alias
			return nil
		}
	}
	return model.Error{
		Description: "tenant not found",
		Code:        model.CodeResourceNotFound,
	}
}

// Get implements the tenants API
func (t *Tenants) Get(id string, _ map[string]string) (*model.Tenant, error) {
	for _, tenant := range t.items {
		if tenant.ID == id {
			return &tenant, nil
		}
	}
	return nil, model.Error{
		Description: "not found",
		Code:        model.CodeResourceNotFound,
	}
}

// List implements the tenants API
func (t *Tenants) List(_ map[string]string) (*model.TenantList, error) {
	return &model.TenantList{Items: t.items}, nil
}

// GetQuota retrieves the quota settings for the given tenant
func (t *Tenants) GetQuota(id string, _ map[string]string) (*model.TenantQuota, error) {
	for _, tenant := range t.items {
		if tenant.ID == id {
			return &model.TenantQuota{
				XMLName:                 tenant.XMLName,
				BlockSize:               tenant.BlockSize,
				NotificationSize:        tenant.NotificationSize,
				BlockSizeInCount:        tenant.BlockSizeInCount,
				NotificationSizeInCount: tenant.NotificationSizeInCount,
				ID:                      tenant.ID,
			}, nil
		}
	}
	return nil, model.Error{
		Description: "not found",
		Code:        model.CodeResourceNotFound,
	}
}

// SetQuota updates the quota settings for the given tenant
func (t *Tenants) SetQuota(id string, tenantQuota model.TenantQuotaSet) error {
	for i, tenant := range t.items {
		if tenant.ID == id {
			t.items[i].BlockSize = tenantQuota.BlockSize
			t.items[i].BlockSizeInCount = tenantQuota.BlockSizeInCount
			t.items[i].NotificationSize = tenantQuota.NotificationSize
			t.items[i].NotificationSizeInCount = tenantQuota.NotificationSizeInCount
			return nil
		}
	}

	return model.Error{
		Description: "not found",
		Code:        model.CodeResourceNotFound,
	}
}

// DeleteQuota deletes the quota settings for the given tenant
func (t *Tenants) DeleteQuota(id string) error {
	for i, tenant := range t.items {
		if tenant.ID == id {
			t.items[i].BlockSize = ""
			t.items[i].BlockSizeInCount = ""
			t.items[i].NotificationSize = ""
			t.items[i].NotificationSizeInCount = ""
			return nil
		}
	}
	return model.Error{
		Description: "not found",
		Code:        model.CodeResourceNotFound,
	}
}

// Buckets implements the buckets API
type Buckets struct {
	items  []model.Bucket
	policy map[string]string
}

var _ api.BucketsInterface = (*Buckets)(nil) // interface guard

// List implements the buckets API
func (b *Buckets) List(_ map[string]string) (*model.BucketList, error) {
	return &model.BucketList{Items: b.items}, nil
}

// Get implements the buckets API
func (b *Buckets) Get(name string, params map[string]string) (*model.Bucket, error) {
	// this is not path, it is used to quickly distinguish which function must fail
	_, ok := params["X-TEST/Buckets/Get/force-fail"]
	if ok {
		return nil, model.Error{
			Description: "An unexpected error occurred",
			Code:        model.CodeInternalException,
		}
	}

	for _, bucket := range b.items {
		if bucket.Name == name {
			return &bucket, nil
		}
	}
	return nil, model.Error{
		Description: "not found",
		Code:        model.CodeParameterNotFound,
	}
}

// GetPolicy implements the buckets API
func (b *Buckets) GetPolicy(bucketName string, params map[string]string) (string, error) {
	// this is not path, it is used to quickly distinguish which function must fail
	_, ok := params["X-TEST/Buckets/GetPolicy/force-fail"]
	if ok {
		return "", model.Error{
			Description: "An unexpected error occurred",
			Code:        model.CodeInternalException,
		}
	}
	if policy, ok := b.policy[fmt.Sprintf("%s/%s", bucketName, params["namespace"])]; ok {
		return policy, nil
	}
	return "", nil
}

// DeletePolicy implements the buckets API
func (b *Buckets) DeletePolicy(bucketName string, params map[string]string) error {
	// this is not path, it is used to quickly distinguish which function must fail
	_, ok := params["X-TEST/Buckets/DeletePolicy/force-fail"]
	if ok {
		return model.Error{
			Description: "An unexpected error occurred",
			Code:        model.CodeInternalException,
		}
	}
	found := false
	if found {
		delete(b.policy, fmt.Sprintf("%s/%s", bucketName, params["namespace"]))
		return nil
	}
	return model.Error{
		Description: "bucket not found",
		Code:        model.CodeResourceNotFound,
	}
}

// UpdatePolicy implements the buckets API
func (b *Buckets) UpdatePolicy(bucketName string, policy string, params map[string]string) error {
	// this is not path, it is used to quickly distinguish which function must fail
	_, ok := params["X-TEST/Buckets/UpdatePolicy/force-fail"]
	if ok {
		return model.Error{
			Description: "An unexpected error occurred",
			Code:        model.CodeInternalException,
		}
	}
	found := false
	if found {
		b.policy[fmt.Sprintf("%s/%s", bucketName, params["namespace"])] = policy
		return nil
	}
	return model.Error{
		Description: "bucket not found",
		Code:        model.CodeResourceNotFound,
	}
}

// Create implements the buckets API
func (b *Buckets) Create(createParams model.Bucket) (*model.Bucket, error) {
	// This piece of code verifies if the incoming request is for forcing an unexpected error.
	if strings.Contains(createParams.Name, "FORCEFAIL") {
		return &createParams, model.Error{
			Description: "Bucket was not sucessfully created",
			Code:        model.CodeInternalException,
		}
	}

	for _, existingBucket := range b.items {
		if existingBucket.Namespace == createParams.Namespace && existingBucket.Name == createParams.Name {
			return nil, model.Error{
				Description: "duplicate found",
				Code:        model.CodeBucketAlreadyExists,
			}
		}
	}
	b.items = append(b.items, createParams)
	return &createParams, nil
}

// Delete implements the buckets API
func (b *Buckets) Delete(name string, namespace string, emptyBucket bool) error {
	// This piece of code verifies if the incoming request is for forcing an unexpected error.
	if strings.Contains(name, "FORCEFAIL") {
		return model.Error{
			Description: "Bucket was not sucessfully deleted",
			Code:        model.CodeInternalException,
		}
	}
	for i, existingBucket := range b.items {
		if existingBucket.Name == name && existingBucket.Namespace == namespace {
			b.items = append(b.items[:i], b.items[i+1:]...)
			return nil
		}
	}
	return model.Error{
		Description: "not found",
		Code:        model.CodeResourceNotFound,
	}
}

// GetQuota gets the quota for the given bucket and namespace.
func (b *Buckets) GetQuota(bucketName string, _ string) (*model.BucketQuotaInfo, error) {
	for _, bucket := range b.items {
		if bucket.Name == bucketName {
			return &model.BucketQuotaInfo{
				BucketQuota: model.BucketQuota{
					BucketName:            bucket.Name,
					Namespace:             bucket.Namespace,
					NotificationSize:      bucket.NotificationSize,
					NotificationSizeCount: bucket.NotificationSizeCount,
					BlockSize:             bucket.BlockSize,
					BlockSizeCount:        bucket.BlockSizeCount,
				},
			}, nil
		}
	}

	return nil, model.Error{
		Description: "not found",
		Code:        model.CodeResourceNotFound,
	}
}

// UpdateQuota updates the quota for the specified bucket.
func (b *Buckets) UpdateQuota(bucketQuota model.BucketQuotaUpdate) error {
	for i := 0; i < len(b.items); i++ {
		if b.items[i].Name == bucketQuota.BucketName {
			b.items[i].BlockSize = bucketQuota.BlockSize
			b.items[i].NotificationSize = bucketQuota.NotificationSize
			b.items[i].NotificationSizeCount = bucketQuota.NotificationSizeCount
			b.items[i].BlockSizeCount = bucketQuota.BlockSizeCount
			return nil
		}
	}
	return model.Error{
		Description: "not found",
		Code:        model.CodeResourceNotFound,
	}
}

// DeleteQuota deletes the quota setting for the given bucket and namespace.
func (b *Buckets) DeleteQuota(bucketName string, _ string) error {
	for i := 0; i < len(b.items); i++ {
		if b.items[i].Name == bucketName {
			b.items[i].BlockSize = -1
			b.items[i].BlockSizeCount = -1
			b.items[i].NotificationSize = -1
			b.items[i].NotificationSizeCount = -1
			return nil
		}
	}
	return model.Error{
		Description: "not found",
		Code:        model.CodeResourceNotFound,
	}
}

// Objmt is a fake (mocked) implementation of the Objmt interface
type Objmt struct {
	accountBillingInfoList      *model.AccountBillingInfoList
	accountBillingSampleList    *model.AccountBillingSampleList
	bucketBillingInfoList       *model.BucketBillingInfoList
	bucketBillingSampleList     *model.BucketBillingSampleList
	bucketBillingPerfList       *model.BucketPerfDataList
	bucketReplicationInfoList   *model.BucketReplicationInfoList
	bucketReplicationSampleList *model.BucketReplicationSampleList
	storeBillingInfoList        *model.StoreBillingInfoList
	storeBillingSampleList      *model.StoreBillingSampleList
	storeReplicationDataList    *model.StoreReplicationDataList
}

var _ api.ObjmtInterface = (*Objmt)(nil) // interface guard

// GetStoreBillingInfo returns billing info metrics for object store
func (mt *Objmt) GetStoreBillingInfo(_ map[string]string) (*model.StoreBillingInfoList, error) {
	return mt.storeBillingInfoList, nil
}

// GetStoreBillingSample returns billing sample (time-window) metrics for object store
func (mt *Objmt) GetStoreBillingSample(_ map[string]string) (*model.StoreBillingSampleList, error) {
	return mt.storeBillingSampleList, nil
}

// GetStoreReplicationData returns CRR metrics for defined object stores
func (mt *Objmt) GetStoreReplicationData(_ []string, _ map[string]string) (*model.StoreReplicationDataList, error) {
	return mt.storeReplicationDataList, nil
}

// GetAccountBillingInfo returns billing info metrics for defined accounts
func (mt *Objmt) GetAccountBillingInfo(_ []string, _ map[string]string) (*model.AccountBillingInfoList, error) {
	return mt.accountBillingInfoList, nil
}

// GetAccountBillingSample returns billing sample (time-window) metrics for defined accounts
func (mt *Objmt) GetAccountBillingSample(_ []string, _ map[string]string) (*model.AccountBillingSampleList, error) {
	return mt.accountBillingSampleList, nil
}

// GetBucketBillingInfo returns billing info metrics for defined buckets and account
func (mt *Objmt) GetBucketBillingInfo(_ string, _ []string, _ map[string]string) (*model.BucketBillingInfoList, error) {
	return mt.bucketBillingInfoList, nil
}

// GetBucketBillingSample returns billing sample (time-window) metrics for defined buckets and account
func (mt *Objmt) GetBucketBillingSample(_ string, _ []string, _ map[string]string) (*model.BucketBillingSampleList, error) {
	return mt.bucketBillingSampleList, nil
}

// GetBucketBillingPerf returns performance metrics for defined buckets and account
func (mt *Objmt) GetBucketBillingPerf(_ string, _ []string, _ map[string]string) (*model.BucketPerfDataList, error) {
	return mt.bucketBillingPerfList, nil
}

// GetReplicationInfo returns billing info metrics for defined replication pairs and account
func (mt *Objmt) GetReplicationInfo(_ string, _ [][]string, _ map[string]string) (*model.BucketReplicationInfoList, error) {
	return mt.bucketReplicationInfoList, nil
}

// GetReplicationSample returns billing sample (time-window) metrics for defined replication pairs and account
func (mt *Objmt) GetReplicationSample(_ string, _ [][]string, _ map[string]string) (*model.BucketReplicationSampleList, error) {
	return mt.bucketReplicationSampleList, nil
}

// CRR implements the crr API
type CRR struct {
	Config *model.CRR
}

var _ api.CRRInterface = (*CRR)(nil) // interface guard

// PauseReplication implements the CRR API
func (c *CRR) PauseReplication(destObjectScale string, destObjectStore string, params map[string]string) error {
	//resume, _ := strconv.Atoi(params["pauseEndMills"])
	resume, _ := strconv.ParseInt(params["pauseEndMills"], 10, 64)
	c.Config.DestObjectScale = destObjectScale
	c.Config.DestObjectStore = destObjectStore
	c.Config.PauseStartMills = int64(time.Millisecond)
	c.Config.PauseEndMills = resume
	return nil
}

// SuspendReplication implements the CRR API
func (c *CRR) SuspendReplication(destObjectScale string, destObjectStore string, _ map[string]string) error {
	c.Config.DestObjectScale = destObjectScale
	c.Config.DestObjectStore = destObjectStore
	return nil
}

// ResumeReplication implements the CRR API
func (c *CRR) ResumeReplication(destObjectScale string, destObjectStore string, _ map[string]string) error {
	c.Config.DestObjectScale = destObjectScale
	c.Config.DestObjectStore = destObjectStore
	return nil
}

// UnthrottleReplication implements the CRR API
func (c *CRR) UnthrottleReplication(destObjectScale string, destObjectStore string, _ map[string]string) error {
	c.Config.DestObjectScale = destObjectScale
	c.Config.DestObjectStore = destObjectStore
	return nil
}

// ThrottleReplication implements the CRR API
func (c *CRR) ThrottleReplication(destObjectScale string, destObjectStore string, param map[string]string) error {
	throttle, _ := strconv.Atoi(param["throttlePerSecond"])
	c.Config.DestObjectScale = destObjectScale
	c.Config.DestObjectStore = destObjectStore
	c.Config.ThrottleBandwidth = throttle
	return nil
}

// Get implements the CRR API
func (c *CRR) Get(destObjectScale string, destObjectStore string, _ map[string]string) (*model.CRR, error) {
	c.Config.DestObjectScale = destObjectScale
	c.Config.DestObjectStore = destObjectStore
	return c.Config, nil
}

// AlertPolicy implements the AlertPolicy API
type AlertPolicy struct {
	AlertPolicy *model.AlertPolicy
}

// AlertPolicies implements the AlertPolicies API
type AlertPolicies struct {
	items []model.AlertPolicy
}

var _ api.AlertPoliciesInterface = (*AlertPolicies)(nil) // interface guard

// Get implements the AlertPolicy API
func (ap *AlertPolicies) Get(policyName string) (*model.AlertPolicy, error) {
	for _, AlertPolicy := range ap.items {
		if AlertPolicy.PolicyName == policyName {
			return &AlertPolicy, nil
		}
	}
	return nil, model.Error{
		Description: "not found",
		Code:        model.CodeResourceNotFound,
	}
}

// List implements the buckets API
func (ap *AlertPolicies) List(_ map[string]string) (*model.AlertPolicies, error) {
	return &model.AlertPolicies{Items: ap.items}, nil
}

// Create implements the AlertPolicy API
func (ap *AlertPolicies) Create(payload model.AlertPolicy) (*model.AlertPolicy, error) {
	newAlertPolicy := &model.AlertPolicy{
		PolicyName:           payload.PolicyName,
		MetricType:           payload.MetricType,
		MetricName:           payload.MetricName,
		CreatedBy:            payload.CreatedBy,
		IsEnabled:            payload.IsEnabled,
		IsPerInstanceMetric:  payload.IsPerInstanceMetric,
		Period:               payload.Period,
		PeriodUnits:          payload.PeriodUnits,
		DatapointsToConsider: payload.DatapointsToConsider,
		DatapointsToAlert:    payload.DatapointsToAlert,
		Statistic:            payload.Statistic,
		Operator:             payload.Operator,
		Condition:            payload.Condition,
	}
	ap.items = append(ap.items, *newAlertPolicy)
	return newAlertPolicy, nil
}

// Delete implements the AlertPolicy API
func (ap *AlertPolicies) Delete(policyName string) error {
	for i, alertpolicy := range ap.items {
		if alertpolicy.PolicyName == policyName {
			ap.items = append(ap.items[:i], ap.items[i+1:]...)
			return nil
		}
	}
	return model.Error{
		Description: "alert policy not found",
		Code:        model.CodeResourceNotFound,
	}
}

// Update implements the AlertPolicy API
func (ap *AlertPolicies) Update(payload model.AlertPolicy, policyName string) (*model.AlertPolicy, error) {
	for i, alertpolicy := range ap.items {
		if alertpolicy.PolicyName == policyName {
			ap.items[i].PolicyName = payload.PolicyName
			ap.items[i].MetricType = payload.MetricType
			ap.items[i].MetricName = payload.MetricName
			ap.items[i].CreatedBy = payload.CreatedBy
			ap.items[i].IsEnabled = payload.IsEnabled
			ap.items[i].IsPerInstanceMetric = payload.IsPerInstanceMetric
			ap.items[i].Period = payload.Period
			ap.items[i].PeriodUnits = payload.PeriodUnits
			ap.items[i].DatapointsToConsider = payload.DatapointsToConsider
			ap.items[i].DatapointsToAlert = payload.DatapointsToAlert
			ap.items[i].Statistic = payload.Statistic
			ap.items[i].Operator = payload.Operator
			ap.items[i].Condition = payload.Condition
			return &alertpolicy, nil
		}
	}
	return nil, model.Error{
		Description: "alert policy not found",
		Code:        model.CodeResourceNotFound,
	}
}

// Status implements the Status API
type Status struct {
	RebuildInfo *model.RebuildInfo
}

var _ api.StatusInterfaces = (*Status)(nil) // interface guard

// GetRebuildStatus implements the Status API
func (s *Status) GetRebuildStatus(objStoreName, ssPodName, ssPodNameSpace, level string, params map[string]string) (*model.RebuildInfo, error) {
	s.RebuildInfo.TotalBytes = 2048
	s.RebuildInfo.RemainingBytes = 1024
	return s.RebuildInfo, nil
}
