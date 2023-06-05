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

package rest

import (
	"github.com/dell/goobjectscale/pkg/client/api"
	"github.com/dell/goobjectscale/pkg/client/rest/alertpolicies"
	"github.com/dell/goobjectscale/pkg/client/rest/buckets"
	"github.com/dell/goobjectscale/pkg/client/rest/client"
	"github.com/dell/goobjectscale/pkg/client/rest/crr"
	"github.com/dell/goobjectscale/pkg/client/rest/federatedobjectstores"
	"github.com/dell/goobjectscale/pkg/client/rest/objectuser"
	"github.com/dell/goobjectscale/pkg/client/rest/objmt"
	"github.com/dell/goobjectscale/pkg/client/rest/status"
	"github.com/dell/goobjectscale/pkg/client/rest/tenants"
)

// ClientSet is a set of clients for each API section
type ClientSet struct {
	client                client.RemoteCaller
	buckets               api.BucketsInterface
	objectUser            api.ObjectUserInterface
	tenants               api.TenantsInterface
	objmt                 api.ObjmtInterface
	crr                   api.CRRInterface
	status                api.StatusInterfaces
	alertPolicies         api.AlertPoliciesInterface
	federatedObjectStores api.FederatedObjectStoresInterface
}

// NewClientSet returns a new client set based on the provided REST client parameters
func NewClientSet(c client.RemoteCaller) *ClientSet {
	return &ClientSet{
		client:                c,
		buckets:               &buckets.Buckets{Client: c},
		objectUser:            &objectuser.ObjectUser{Client: c},
		tenants:               &tenants.Tenants{Client: c},
		objmt:                 &objmt.Objmt{Client: c},
		crr:                   &crr.CRR{Client: c},
		alertPolicies:         &alertpolicies.AlertPolicies{Client: c},
		status:                &status.Status{Client: c},
		federatedObjectStores: &federatedobjectstores.FederatedObjectStores{Client: c},
	}
}

// Client returns the REST client used in the ClientSet
func (c *ClientSet) Client() client.RemoteCaller {
	return c.client
}

// Buckets implements the client API
func (c *ClientSet) Buckets() api.BucketsInterface {
	return c.buckets
}

// ObjectUser implements the client API
func (c *ClientSet) ObjectUser() api.ObjectUserInterface {
	return c.objectUser
}

// AlertPolicies implements the client API
func (c *ClientSet) AlertPolicies() api.AlertPoliciesInterface {
	return c.alertPolicies
}

// Tenants implements the client API
func (c *ClientSet) Tenants() api.TenantsInterface {
	return c.tenants
}

// ObjectMt implements the client API for objMT metrics
func (c *ClientSet) ObjectMt() api.ObjmtInterface {
	return c.objmt
}

// CRR implements the client API for Cross Region Replication
func (c *ClientSet) CRR() api.CRRInterface {
	return c.crr
}

// Status implements the client API
func (c *ClientSet) Status() api.StatusInterfaces {
	return c.status
}

// FederatedObjectStores implements the client API
func (c *ClientSet) FederatedObjectStores() api.FederatedObjectStoresInterface {
	return c.federatedObjectStores
}
