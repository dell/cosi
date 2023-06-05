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

package alertpolicies

import (
	"net/http"
	"path"

	"github.com/dell/goobjectscale/pkg/client/api"
	"github.com/dell/goobjectscale/pkg/client/model"
	"github.com/dell/goobjectscale/pkg/client/rest/client"
)

// AlertPolicies is a REST implementation of the AlertPolicies interface
type AlertPolicies struct {
	Client client.RemoteCaller
}

var _ api.AlertPoliciesInterface = &AlertPolicies{} // interface guard

// Get implements the AlertPolicy interface
func (ap *AlertPolicies) Get(policyName string) (*model.AlertPolicy, error) {
	req := client.Request{
		Method:      http.MethodGet,
		Path:        path.Join("vdc", "alertpolicy", policyName),
		ContentType: client.ContentTypeXML,
	}
	alertpolicy := &model.AlertPolicy{}
	err := ap.Client.MakeRemoteCall(req, alertpolicy)
	if err != nil {
		return nil, err
	}
	return alertpolicy, nil
}

// List implements the AlertPolicy interface
func (ap *AlertPolicies) List(params map[string]string) (*model.AlertPolicies, error) {
	req := client.Request{
		Method:      http.MethodGet,
		Path:        path.Join("vdc", "alertpolicy", "list"),
		ContentType: client.ContentTypeXML,
		Params:      params,
	}
	alertpolicies := &model.AlertPolicies{}
	err := ap.Client.MakeRemoteCall(req, alertpolicies)
	if err != nil {
		return nil, err
	}

	return alertpolicies, nil
}

// Create implements the AlertPolicy interface
func (ap *AlertPolicies) Create(payload model.AlertPolicy) (*model.AlertPolicy, error) {
	req := client.Request{
		Method:      http.MethodPost,
		Path:        path.Join("vdc", "alertpolicy"),
		ContentType: client.ContentTypeXML,
		Body:        &payload,
	}
	alertpolicy := &model.AlertPolicy{}
	err := ap.Client.MakeRemoteCall(req, alertpolicy)
	if err != nil {
		return nil, err
	}
	return alertpolicy, nil
}

// Delete implements the AlertPolicy interface
func (ap *AlertPolicies) Delete(policyName string) error {
	req := client.Request{
		Method:      http.MethodDelete,
		Path:        path.Join("vdc", "alertpolicy", policyName),
		ContentType: client.ContentTypeXML,
	}
	alertpolicy := &model.AlertPolicy{}
	err := ap.Client.MakeRemoteCall(req, alertpolicy)
	if err != nil {
		return err
	}
	return nil
}

// Update implements the AlertPolicy interface
func (ap *AlertPolicies) Update(payload model.AlertPolicy, policyName string) (*model.AlertPolicy, error) {
	req := client.Request{
		Method:      http.MethodPut,
		Path:        path.Join("vdc", "alertpolicy", policyName),
		ContentType: client.ContentTypeXML,
		Body:        &payload,
	}
	alertpolicy := &model.AlertPolicy{}
	err := ap.Client.MakeRemoteCall(req, alertpolicy)
	if err != nil {
		return nil, err
	}
	return alertpolicy, nil
}
