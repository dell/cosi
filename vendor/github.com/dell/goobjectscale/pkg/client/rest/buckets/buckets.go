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

package buckets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/dell/goobjectscale/pkg/client/model"
	"github.com/dell/goobjectscale/pkg/client/rest/client"
)

// Buckets is a REST implementation of the Buckets interface
type Buckets struct {
	Client client.RemoteCaller
}

// Get implements the buckets interface
func (b *Buckets) Get(name string, params map[string]string) (*model.Bucket, error) {
	req := client.Request{
		Method:      http.MethodGet,
		Path:        path.Join("object", "bucket", name, "info"),
		ContentType: client.ContentTypeXML,
		Params:      params,
	}
	bucket := &model.BucketInfo{}
	err := b.Client.MakeRemoteCall(req, bucket)
	if err != nil {
		return nil, err
	}
	return &bucket.Bucket, nil
}

// List implements the buckets interface
func (b *Buckets) List(params map[string]string) (*model.BucketList, error) {
	req := client.Request{
		Method:      http.MethodGet,
		Path:        "/object/bucket",
		ContentType: client.ContentTypeXML,
		Params:      params,
	}
	bucketList := &model.BucketList{}
	err := b.Client.MakeRemoteCall(req, bucketList)
	if err != nil {
		return nil, err
	}
	return bucketList, nil
}

// GetPolicy implements the buckets interface
func (b *Buckets) GetPolicy(bucketName string, param map[string]string) (string, error) {
	req := client.Request{
		Method:      http.MethodGet,
		Path:        fmt.Sprintf("object/bucket/%s/policy", bucketName),
		ContentType: client.ContentTypeJSON,
		Params:      param,
	}
	var bucketPolicy json.RawMessage
	err := b.Client.MakeRemoteCall(req, &bucketPolicy)
	if err != nil {
		return "", err
	}
	policy, err := bucketPolicy.MarshalJSON()
	return string(policy), err
}

// UpdatePolicy implements the buckets interface
func (b *Buckets) UpdatePolicy(bucketName string, policy string, param map[string]string) error {
	req := client.Request{
		Method:      http.MethodPut,
		Path:        fmt.Sprintf("object/bucket/%s/policy", bucketName),
		ContentType: client.ContentTypeJSON,
		Params:      param,
		Body:        json.RawMessage(policy),
	}
	return b.Client.MakeRemoteCall(req, nil)
}

// DeletePolicy implements the buckets interface
func (b *Buckets) DeletePolicy(bucketName string, param map[string]string) error {
	req := client.Request{
		Method:      http.MethodDelete,
		Path:        fmt.Sprintf("object/bucket/%s/policy", bucketName),
		ContentType: client.ContentTypeJSON,
		Params:      param,
	}
	return b.Client.MakeRemoteCall(req, nil)
}

// Create implements the buckets interface
func (b *Buckets) Create(createParam model.Bucket) (*model.Bucket, error) {
	req := client.Request{
		Method:      http.MethodPost,
		Path:        "/object/bucket",
		ContentType: client.ContentTypeXML,
		Body:        &model.BucketCreate{Bucket: createParam},
	}
	bucket := &model.Bucket{}
	err := b.Client.MakeRemoteCall(req, bucket)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

// Delete implements the buckets interface
func (b *Buckets) Delete(name string, namespace string, emptyBucket bool) error {
	req := client.Request{
		Method:      http.MethodPost,
		Path:        path.Join("object", "bucket", name, "deactivate"),
		Params:      map[string]string{"namespace": namespace, "emptyBucket": fmt.Sprint(emptyBucket)},
		ContentType: client.ContentTypeJSON,
	}
	err := b.Client.MakeRemoteCall(req, nil)
	if err != nil {
		return err
	}
	return nil
}

// GetQuota gets the quota for the given bucket and namespace.
func (b *Buckets) GetQuota(bucketName string, namespace string) (*model.BucketQuotaInfo, error) {
	req := client.Request{
		Method:      http.MethodGet,
		Path:        fmt.Sprintf("object/bucket/%s/quota", bucketName),
		ContentType: client.ContentTypeXML,
		Params:      map[string]string{"namespace": namespace},
	}
	bucketQuota := &model.BucketQuotaInfo{}
	err := b.Client.MakeRemoteCall(req, bucketQuota)
	if err != nil {
		return nil, err
	}
	return bucketQuota, err
}

// UpdateQuota updates the quota for the specified bucket.
func (b *Buckets) UpdateQuota(bucketQuota model.BucketQuotaUpdate) error {
	req := client.Request{
		Method:      http.MethodPut,
		Path:        fmt.Sprintf("object/bucket/%s/quota", bucketQuota.BucketName),
		ContentType: client.ContentTypeXML,
		Body:        bucketQuota,
	}
	return b.Client.MakeRemoteCall(req, nil)
}

// DeleteQuota deletes the quota setting for the given bucket and namespace.
func (b *Buckets) DeleteQuota(bucketName string, namespace string) error {
	req := client.Request{
		Method:      http.MethodDelete,
		Path:        fmt.Sprintf("object/bucket/%s/quota", bucketName),
		ContentType: client.ContentTypeXML,
		Params:      map[string]string{"namespace": namespace},
	}
	return b.Client.MakeRemoteCall(req, nil)
}
