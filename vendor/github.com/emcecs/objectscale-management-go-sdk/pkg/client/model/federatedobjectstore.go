package model

import "encoding/xml"

// FederatedObjectStoreList is a list of federated object stores
type FederatedObjectStoreList struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"ReplicationInfo"`

	// Items is the list of federated object stores in the list
	Items []FederatedObjectStore `xml:"ReplicationStoreInfo"`
}

type FederatedObjectStore struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"ReplicationStoreInfo"`

	CRRConfigured bool `xml:"CRRConfigured"`

	ObjectScaleId string `xml:"ObjectScaleId"`

	ObjectStoreId string `xml:"ObjectStoreId"`

	ObjectStoreName string `xml:"ObjectStoreName"`

	ReplicationStatus string `xml:"ReplicationStatus"`

	ObjectStoreRTO int64 `xml:"ObjectStoreRTO,omitempty"`

	FailedData int64 `xml:"FailedData,omitempty"`

	CRRControlParameters CRRControlParameters `xml:"CRRControlParameters"`
}

type CRRControlParameters struct {
	XMLName xml.Name `xml:"CRRControlParameters"`

	SuspendStartMills int64 `xml:"suspendStartMills,omitempty"`

	PauseStartMills int64 `xml:"pauseStartMills,omitempty"`

	PauseEndMills int64 `xml:"pauseEndMills,omitempty"`

	ThrottleBandwidth int `xml:"throttleBandwidth,omitempty"`
}
