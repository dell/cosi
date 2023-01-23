package model

import "encoding/xml"

type CRR struct {
	XMLName xml.Name `xml:"ReplicationAdminConfiguration"`

	DestObjectScale string `xml:"destinationObjectScale"`

	DestObjectStore string `xml:"destinationObjectStore"`

	PauseStartMills int64 `xml:"pauseStartMills"`

	PauseEndMills int64 `xml:"pauseEndMills"`

	SuspendStartMills int64 `xml:"suspendStartMills"`

	ThrottleBandwidth int `xml:"throttleBandwidth"`
}
