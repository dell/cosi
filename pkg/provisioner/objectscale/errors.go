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

package objectscale

import "github.com/dell/goobjectscale/pkg/client/model"

var (
	// ErrParameterNotFound is general instance of the model.Error with the CodeParameterNotFound.
	// It indicates that request parameter cannot be found - e.g. requested bucket does not exist.
	ErrParameterNotFound = model.Error{Code: model.CodeParameterNotFound}

	// ErrParameterNotFound is general instance of the model.Error with the CodeInternalException.
	// It indicates that internal exception occurred, and user should look at ObjectScale logs to find the cause.
	ErrInternalException = model.Error{Code: model.CodeInternalException}
)
