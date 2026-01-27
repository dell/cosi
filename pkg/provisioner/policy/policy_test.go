// Copyright © 2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

package policy_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dell/cosi/pkg/provisioner/policy"
	"github.com/stretchr/testify/assert"
)

func TestNewFromJSON(t *testing.T) {
	tests := []struct {
		name        string
		jsonString  string
		expected    policy.Document
		expectedErr bool
	}{
		{
			name: "valid JSON",
			jsonString: `{
				"Version": "2012-10-17",
				"Id": "policy-id",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": ["s3:GetObject"],
						"Resource": ["arn:aws:s3:::my-bucket/*"]
					}
				]
			}`,
			expected: policy.Document{
				Version: "2012-10-17",
				ID:      "policy-id",
				Statement: []policy.StatementEntry{
					{
						Effect:   "Allow",
						Action:   []string{"s3:GetObject"},
						Resource: []string{"arn:aws:s3:::my-bucket/*"},
					},
				},
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := policy.NewFromJSON(tt.jsonString)
			if (err != nil) != tt.expectedErr {
				t.Errorf("NewFromJSON() error = %v, wantErr %v", err, tt.expectedErr)
				return
			}
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("NewFromJSON() = %v, want %v", actual, tt.expected)
			}
		})
	}
}

func TestToJSON(t *testing.T) {
	tests := []struct {
		name        string
		jsonString  string
		expected    policy.Document
		expectedErr bool
	}{
		{
			name: "toJSON equals",
			jsonString: `{
				"Version": "2012-10-17",
				"Id": "policy-id",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": ["s3:GetObject"],
						"Resource": ["arn:aws:s3:::my-bucket/*"]
					}
				]
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := policy.NewFromJSON(tt.jsonString)
			assert.Nil(t, err)
			generatedJSON, err := doc.ToJSON()

			condensed := strings.ReplaceAll(tt.jsonString, " ", "")
			condensed = strings.ReplaceAll(condensed, "\n", "")
			condensed = strings.ReplaceAll(condensed, "\t", "")

			assert.Nil(t, err)
			assert.Equal(t, condensed, generatedJSON)
		})
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name        string
		jsonString1 string
		jsonString2 string
		isEqual     bool
	}{
		{
			name: "valid JSON",
			jsonString1: `{
				"Version": "2012-10-17",
				"Id": "policy-id",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": ["s3:GetObject"],
						"Resource": ["arn:aws:s3:::my-bucket/*"]
					}
				]
			}`,
			jsonString2: `{
				"Version": "2012-10-17",
				"Id": "policy-id",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": ["s3:GetObject"],
						"Resource": ["arn:aws:s3:::my-bucket/*"]
					}
				]
			}`,
			isEqual: true,
		},
		{
			name: "different version",
			jsonString1: `{
				"Version": "2012-10-17",
				"Id": "policy-id",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": ["s3:GetObject"],
						"Resource": ["arn:aws:s3:::my-bucket/*"]
					}
				]
			}`,
			jsonString2: `{
				"Version": "2012-10-18",
				"Id": "policy-id",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": ["s3:GetObject"],
						"Resource": ["arn:aws:s3:::my-bucket/*"]
					}
				]
			}`,
			isEqual: false,
		},
		{
			name: "different id",
			jsonString1: `{
				"Version": "2012-10-17",
				"Id": "policy-id-1",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": ["s3:GetObject"],
						"Resource": ["arn:aws:s3:::my-bucket/*"]
					}
				]
			}`,
			jsonString2: `{
				"Version": "2012-10-17",
				"Id": "policy-id-2",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": ["s3:GetObject"],
						"Resource": ["arn:aws:s3:::my-bucket/*"]
					}
				]
			}`,
			isEqual: false,
		},
		{
			name: "different statement effect",
			jsonString1: `{
				"Version": "2012-10-17",
				"Id": "policy-id",
				"Statement": [
					{
						"Effect": "Deny",
						"Action": ["s3:GetObject"],
						"Resource": ["arn:aws:s3:::my-bucket/*"]
					}
				]
			}`,
			jsonString2: `{
				"Version": "2012-10-17",
				"Id": "policy-id",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": ["s3:GetObject"],
						"Resource": ["arn:aws:s3:::my-bucket/*"]
					}
				]
			}`,
			isEqual: false,
		},
		{
			name: "different action",
			jsonString1: `{
				"Version": "2012-10-17",
				"Id": "policy-id",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": ["s3:PutObject"],
						"Resource": ["arn:aws:s3:::my-bucket/*"]
					}
				]
			}`,
			jsonString2: `{
				"Version": "2012-10-17",
				"Id": "policy-id",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": ["s3:GetObject"],
						"Resource": ["arn:aws:s3:::my-bucket/*"]
					}
				]
			}`,
			isEqual: false,
		},
		{
			name: "different resource",
			jsonString1: `{
				"Version": "2012-10-17",
				"Id": "policy-id",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": ["s3:GetObject"],
						"Resource": ["arn:aws:s3:::my-bucket-1/*"]
					}
				]
			}`,
			jsonString2: `{
				"Version": "2012-10-17",
				"Id": "policy-id",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": ["s3:GetObject"],
						"Resource": ["arn:aws:s3:::my-bucket-2/*"]
					}
				]
			}`,
			isEqual: false,
		},
		{
			name: "different action length",
			jsonString1: `{
				"Version": "2012-10-17",
				"Id": "policy-id",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": ["s3:GetObject", "s3:PutObject"],
						"Resource": ["arn:aws:s3:::my-bucket/*"]
					}
				]
			}`,
			jsonString2: `{
				"Version": "2012-10-17",
				"Id": "policy-id",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": ["s3:GetObject"],
						"Resource": ["arn:aws:s3:::my-bucket/*"]
					}
				]
			}`,
			isEqual: false,
		},
		{
			name: "different resource length",
			jsonString1: `{
				"Version": "2012-10-17",
				"Id": "policy-id",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": ["s3:GetObject"],
						"Resource": ["arn:aws:s3:::my-bucket/*", "arn:aws:s3:::my-bucket-2/*"]
					}
				]
			}`,
			jsonString2: `{
				"Version": "2012-10-17",
				"Id": "policy-id",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": ["s3:GetObject"],
						"Resource": ["arn:aws:s3:::my-bucket/*"]
					}
				]
			}`,
			isEqual: false,
		},
		{
			name: "different statement length",
			jsonString1: `{
				"Version": "2012-10-17",
				"Id": "policy-id",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": ["s3:GetObject"],
						"Resource": ["arn:aws:s3:::my-bucket/*"]
					}
				]
			}`,
			jsonString2: `{
				"Version": "2012-10-17",
				"Id": "policy-id",
				"Statement": []
			}`,
			isEqual: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc1, err := policy.NewFromJSON(tt.jsonString1)
			assert.Nil(t, err)
			doc2, err := policy.NewFromJSON(tt.jsonString2)
			assert.Nil(t, err)
			assert.Equal(t, tt.isEqual, doc1.Equal(&doc2))
		})
	}
}
