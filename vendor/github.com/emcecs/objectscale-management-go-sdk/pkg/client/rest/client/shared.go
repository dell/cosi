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

package client

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/model"
)

// RemoteCaller interface is used to create backend calls.
// into represents type, _into_ which data will be unmarshalled.
// Naming follows Effective Go naming convention https://go.dev/doc/effective_go#interface-names
type RemoteCaller interface {
	MakeRemoteCall(r Request, into interface{}) error
}

// HandleResponse handles custom behavior based on server response
func HandleResponse(resp *http.Response) error {
	if resp.StatusCode > 399 {
		switch resp.Body {
		case nil:
			switch {
			case resp.Status != "":
				return fmt.Errorf("server error: %s", strings.ToLower(resp.Status))
			case resp.StatusCode != 0:
				return fmt.Errorf("server error: status code %d", resp.StatusCode)
			}
		default:
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			apiError := &model.Error{}
			err = xml.Unmarshal(body, apiError)
			if err != nil {
				return err
			}
			return fmt.Errorf("server error: %w", apiError)
		}
	}
	return nil
}
