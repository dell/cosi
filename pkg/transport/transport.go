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

package transport

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/dell/cosi-driver/pkg/config"
)

func New(cfg config.Tls) (*http.Transport, error) {
	var tlsConfig *tls.Config
	if cfg.Insecure {
		/* #nosec */
		tlsConfig = &tls.Config{InsecureSkipVerify: true}
	} else {
		// TODO: this need reevaluation, and valid implementation
		return nil, fmt.Errorf("secure transport not implemented")
	}

	return &http.Transport{
		TLSClientConfig: tlsConfig,
	}, nil
}
