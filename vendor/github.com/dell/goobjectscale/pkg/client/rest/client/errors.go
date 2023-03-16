//
//
//  Copyright © 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
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

import "errors"

var (
	// ErrAuthorization is returned when the client is unable to authenticate with the server.
	ErrAuthorization = errors.New("authorization")

	// ErrContentType is returned when the client or server responds with an unknown content type header.
	ErrContentType = errors.New("content type")
)
