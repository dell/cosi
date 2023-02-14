//Copyright © 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
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

package pkg

import (
	"context"
	"testing"
)

// FIXME: those are only smoke tests, no real testing is done here
func TestNewDriver(t *testing.T) {
	idSrv, provSrv, err := NewDriver(context.TODO(), "smoke-test")
	if err != nil {
		t.Errorf("should not return error, got: %s", err.Error())
	}
	if idSrv == nil {
		t.Errorf("identity server should not be nil")
	}
	if provSrv == nil {
		t.Errorf("provisioner server should not be nil")
	}
}
