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

package federatedobjectstores

import (
	"net/http"

	"github.com/dell/goobjectscale/pkg/client/api"
	"github.com/dell/goobjectscale/pkg/client/model"
	"github.com/dell/goobjectscale/pkg/client/rest/client"
)

// FederatedObjectStores is a REST implementation of the FederateObjectStores interface
type FederatedObjectStores struct {
	Client client.RemoteCaller
}

var _ api.FederatedObjectStoresInterface = &FederatedObjectStores{} // interface guard

// List implements the federatedobjectstores interface
func (t *FederatedObjectStores) List(params map[string]string) (*model.FederatedObjectStoreList, error) {
	req := client.Request{
		Method:      http.MethodGet,
		Path:        "/replication/info",
		ContentType: client.ContentTypeXML,
		Params:      params,
	}
	federatedStoreList := &model.FederatedObjectStoreList{}
	err := t.Client.MakeRemoteCall(req, federatedStoreList)
	if err != nil {
		return nil, err
	}
	return federatedStoreList, nil
}
