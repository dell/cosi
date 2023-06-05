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

import "encoding/xml"

// FederatedObjectStoreList is a list of federated object stores
type FederatedObjectStoreList struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"ReplicationInfo"`

	// Items is the list of federated object stores in the list
	Items []FederatedObjectStore `xml:"ReplicationStoreInfo"`
}

// FederatedObjectStore is an ObjectStore which is Federated
type FederatedObjectStore struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"ReplicationStoreInfo"`

	CRRConfigured bool `xml:"CRRConfigured"`

	ObjectScaleID string `xml:"ObjectScaleId"`

	ObjectStoreID string `xml:"ObjectStoreId"`

	ObjectStoreName string `xml:"ObjectStoreName"`

	ReplicationStatus string `xml:"ReplicationStatus"`

	ObjectStoreRTO int64 `xml:"ObjectStoreRTO,omitempty"`

	FailedData int64 `xml:"FailedData,omitempty"`

	CRRControlParameters CRRControlParameters `xml:"CRRControlParameters"`
}

// CRRControlParameters represents parameters for Cross Region Replication
type CRRControlParameters struct {
	XMLName xml.Name `xml:"CRRControlParameters"`

	SuspendStartMills int64 `xml:"suspendStartMills,omitempty"`

	PauseStartMills int64 `xml:"pauseStartMills,omitempty"`

	PauseEndMills int64 `xml:"pauseEndMills,omitempty"`

	ThrottleBandwidth int `xml:"throttleBandwidth,omitempty"`
}
