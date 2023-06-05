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

package model

// RebuildInfo is the information about the rebuild status
type RebuildInfo struct {
	// Status of the storage server
	Status string `json:"status"`

	// TotalBytes of the storage server
	TotalBytes int `json:"total_bytes,string"`

	// RemainingBytes of the storage server
	RemainingBytes int `json:"remaining_bytes,string"`

	// Level of the storage server
	Level int `json:"level,string"`

	// Disk of the storage node
	Disk string `json:"disk"`

	// Message from the storage server
	Message string `json:"message"`

	// Host of the storage server
	Host string `json:"host"`

	// Progress of the recovery
	Progress string `json:"progress"`
}
