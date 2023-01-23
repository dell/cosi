package model

import "encoding/xml"

type Error struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `xml:"error"`

	// Code is the error code returned by the management API
	Code int64 `xml:"code" json:"code"`

	// Description is a human readable description of the error
	Description string `xml:"description" json:"description"`

	// Details are additional error information
	Details string `xml:"details" json:"details"`

	// Retryable signifies if a request returning an error of this type should
	// be retried
	Retryable bool `xml:"retryable" json:"retryable"`
}
