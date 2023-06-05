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

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// Error implements custom error
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

var _ error = Error{}

// Error is a method that allows us to use the Error model as go error
func (err Error) Error() string {
	if err.Description == "" {
		err.Description = "Unknown"
	}

	if err.Details != "" {
		return fmt.Sprintf("%s: %s", err.Description, err.Details)
	}
	return err.Description
}

// StatusCode is there so we can reference the Code field in Is method
func (err Error) StatusCode() int64 {
	return err.Code
}

// Is compare errors status
func (err Error) Is(target error) bool {
	// create intermidiate interface
	type statusCoder interface {
		StatusCode() int64
	}

	// validate if target implements statusCoder interface,
	// and compare the status codes
	switch target := target.(type) {
	case statusCoder:
		return err.StatusCode() == target.StatusCode()

	default:
		// if someone is already relying on error message comparission, then don't break it
		return strings.EqualFold(err.Error(), target.Error())
	}
}

// Error Codes
const (
	// Request parameter cannot be found
	CodeParameterNotFound int64 = 1004
	// Required parameter is missing or empty
	CodeMissingParameter int64 = 1005
	// Resource not found
	CodeResourceNotFound int64 = 1019
	// Exceeding limit
	CodeExceedingLimit int64 = 1031
	// Internal exception occurred
	CodeInternalException int64 = 30024
	// Bucket already exists
	CodeBucketAlreadyExists int64 = 40008
)
