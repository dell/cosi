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

// Package iamapi contains interface that is used to generate mock for AWS IAMAPI.
package iamapi

import "github.com/aws/aws-sdk-go/service/iam/iamiface"

//go:generate go run github.com/vektra/mockery/v2@latest

// IAMAPI interface is an interface based on which the Mock is generated.
type IAMAPI interface {
	iamiface.IAMAPI
}
