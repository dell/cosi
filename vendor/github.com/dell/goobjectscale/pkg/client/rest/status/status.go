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

package status

import (
	"net/http"

	"github.com/dell/goobjectscale/pkg/client/model"
	"github.com/dell/goobjectscale/pkg/client/rest/client"
)

// Status is a REST implementation of the Status interface
type Status struct {
	Client client.RemoteCaller
}

// GetRebuildStatus implements the status interface
func (b *Status) GetRebuildStatus(objStoreName, ssPodName, ssPodNameSpace, level string, params map[string]string) (*model.RebuildInfo, error) {
	// URL Example: https://10.240.117.5:4443/vdc/recovery-status/devices/youmin-test-1-ss-0.youmin-test-1-ss.dellemc-globalmarcu-domain-c10.svc.cluster.local/levels/1
	requestURL := "vdc/recovery-status/devices/" + ssPodName + "." +
		objStoreName + "-ss." + ssPodNameSpace + ".svc.cluster.local/levels/" + level
	req := client.Request{
		Method:      http.MethodGet,
		Path:        requestURL,
		ContentType: client.ContentTypeJSON,
		Params:      params,
	}
	rebuildInfo := &model.RebuildInfo{}
	err := b.Client.MakeRemoteCall(req, rebuildInfo)
	if err != nil {
		return nil, err
	}

	return rebuildInfo, nil
}
