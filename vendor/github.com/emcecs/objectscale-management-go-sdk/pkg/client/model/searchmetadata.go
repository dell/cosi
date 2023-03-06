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

package model

// SearchMetadata is the custom metadata for enabled for querying on the
// bucket
type SearchMetadata struct {
	// MaxKeys indicates the maximum number of metadata search keys for the
	// bucket
	MaxKeys int `json:"maxKeys" xml:"maxKeys"`

	// Enabled indicates if search metadata is enabled for the bucket
	Enabled bool `json:"isEnabled" xml:"isEnabled"`

	// Metadata defines the fields that can be searched upon
	Metadata `json:"metadata" xml:"metadata"`
}

// Metadata defines the fields that can be searched upon
type Metadata struct {
	// Type is the metadata key type
	Type string `json:"type" xml:"type"`

	// Name is the metadata key name
	Name string `json:"name" xml:"name"`

	// Datatype is the data type of the metadata value
	Datatype string `json:"datatype" xml:"datatype"`
}
