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

// Package objectuser contains API to work with Object user resource of object scale.
package objectuser

import (
	"fmt"
	"net/http"

	"github.com/dell/goobjectscale/pkg/client/model"
	"github.com/dell/goobjectscale/pkg/client/rest/client"
)

// ObjectUser is a REST implementation of the object user interface
type ObjectUser struct {
	Client client.RemoteCaller
}

// GetInfo returns information about an object user within the ObjectScale object store.
func (o *ObjectUser) GetInfo(uid string, params map[string]string) (*model.ObjectUserInfo, error) {
	req := client.Request{
		Method:      http.MethodGet,
		Path:        fmt.Sprintf("object/users/%s/info", uid),
		ContentType: client.ContentTypeJSON,
		Params:      params,
	}
	ou := &model.ObjectUserInfo{}
	err := o.Client.MakeRemoteCall(req, ou)
	if err != nil {
		return nil, err
	}
	return ou, nil
}

// GetSecret returns information about object user secrets.
func (o *ObjectUser) GetSecret(uid string, params map[string]string) (*model.ObjectUserSecret, error) {
	req := client.Request{
		Method:      http.MethodGet,
		Path:        fmt.Sprintf("/object/user-secret-keys/%s", uid),
		ContentType: client.ContentTypeJSON,
		Params:      params,
	}
	ou := &model.ObjectUserSecret{}
	err := o.Client.MakeRemoteCall(req, ou)
	if err != nil {
		return nil, err
	}
	return ou, nil
}

// CreateSecret creates secret for a user.
func (o *ObjectUser) CreateSecret(uid string, key model.ObjectUserSecretKeyCreateReq, params map[string]string) (*model.ObjectUserSecretKeyCreateRes, error) {
	req := client.Request{
		Method:      http.MethodPost,
		Path:        fmt.Sprintf("/object/user-secret-keys/%s", uid),
		ContentType: client.ContentTypeJSON,
		Body:        &key,
		Params:      params,
	}
	resp := &model.ObjectUserSecretKeyCreateRes{}
	err := o.Client.MakeRemoteCall(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteSecret deletes secret for a user.
func (o *ObjectUser) DeleteSecret(uid string, key model.ObjectUserSecretKeyDeleteReq, params map[string]string) error {
	req := client.Request{
		Method:      http.MethodPost,
		Path:        fmt.Sprintf("/object/user-secret-keys/%s/deactivate", uid),
		ContentType: client.ContentTypeJSON,
		Body:        &key,
		Params:      params,
	}
	return o.Client.MakeRemoteCall(req, nil)
}

// List returns a list of object users within the ObjectScale object store.
func (o *ObjectUser) List(params map[string]string) (*model.ObjectUserList, error) {
	req := client.Request{
		Method:      http.MethodGet,
		Path:        "object/users",
		ContentType: client.ContentTypeJSON,
		Params:      params,
	}
	ouList := &model.ObjectUserList{}
	err := o.Client.MakeRemoteCall(req, ouList)
	if err != nil {
		return nil, err
	}
	return ouList, nil
}
