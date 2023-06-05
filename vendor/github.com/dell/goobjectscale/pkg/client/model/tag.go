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

// TagSet is a list of tags
type TagSet struct {
	Tags []Tag `json:"Tag,omitempty" xml:"Tag,omitempty"`
}

// Tag is an arbitrary piece of metadata applied to an object
type Tag struct {
	// Key is the tag label or name of the tag
	Key string `json:"Key" xml:"Key"`

	// Value is the value of the tag
	Value string `json:"Value" xml:"Value"`
}
