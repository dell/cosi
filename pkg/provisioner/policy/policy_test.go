// Copyright © 2025 Dell Inc. or its subsidiaries. All Rights Reserved.
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

// Package policy exists for handling operations on AWS Policies,
// defines structures and functions for processing and comparing policies.
package policy

import (
	"testing"
)

func TestToJSON(t *testing.T) {
	tests := []struct {
		name     string
		document Document
		want     string
		wantErr  bool
	}{
		{
			name: "valid policy",
			document: Document{
				Version: "1",
				ID:      "id",
				Statement: []StatementEntry{
					{
						Effect: "Allow",
						Action: []string{"s3:*"},
						Resource: []string{
							"arn:aws:s3:::*",
							"arn:aws:s3:::bucket/*",
						},
					},
				},
			},
			want:    `{"Version":"1","Id":"id","Statement":[{"Effect":"Allow","Action":["s3:*"],"Resource":["arn:aws:s3:::*","arn:aws:s3:::bucket/*"],"Principal":{"AWS":null},"Sid":""}]}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.document.ToJSON()
			if err != nil {
				t.Errorf("ToJSON() error = %v", err)
			}

			if got != tt.want {
				t.Errorf("ToJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFromJSON(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		want    Document
		wantErr bool
	}{
		{
			name:    "valid JSON",
			jsonStr: `{"Version":"1","Id":"id","Statement":[{"Effect":"Allow","Action":["s3:*"],"Resource":["arn:aws:s3:::*","arn:aws:s3:::bucket/*"],"Principal":{"AWS":null},"Sid":""}]}`,
			want: Document{
				Version: "1",
				ID:      "id",
				Statement: []StatementEntry{
					{
						Effect: "Allow",
						Action: []string{"s3:*"},
						Resource: []string{
							"arn:aws:s3:::*",
							"arn:aws:s3:::bucket/*",
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFromJSON(tt.jsonStr)
			if err != nil {
				t.Errorf("NewFromJSON() error = %v", err)
				return
			}

			if !got.Equal(&tt.want) {
				t.Errorf("NewFromJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocumentEqual(t *testing.T) {
	tests := []struct {
		name string
		doc1 Document
		doc2 Document
		want bool
	}{
		{
			name: "equal documents",
			doc1: Document{
				Version: "1",
				ID:      "id",
				Statement: []StatementEntry{
					{
						Effect: "Allow",
						Action: []string{"s3:*"},
						Resource: []string{
							"arn:aws:s3:::*",
							"arn:aws:s3:::bucket/*",
						},
					},
				},
			},
			doc2: Document{
				Version: "1",
				ID:      "id",
				Statement: []StatementEntry{
					{
						Effect: "Allow",
						Action: []string{"s3:*"},
						Resource: []string{
							"arn:aws:s3:::*",
							"arn:aws:s3:::bucket/*",
						},
					},
				},
			},
			want: true,
		},

		{
			name: "different versions",
			doc1: Document{
				Version: "1",
				ID:      "id",
				Statement: []StatementEntry{
					{
						Effect: "Allow",
						Action: []string{"s3:*"},
						Resource: []string{
							"arn:aws:s3:::*",
							"arn:aws:s3:::bucket/*",
						},
					},
				},
			},
			doc2: Document{
				Version: "2",
				ID:      "id",
				Statement: []StatementEntry{
					{
						Effect: "Allow",
						Action: []string{"s3:*"},
						Resource: []string{
							"arn:aws:s3:::*",
							"arn:aws:s3:::bucket/*",
						},
					},
				},
			},
			want: false,
		},
		{
			name: "different IDs",
			doc1: Document{
				Version: "1",
				ID:      "id1",
				Statement: []StatementEntry{
					{
						Effect: "Allow",
						Action: []string{"s3:*"},
						Resource: []string{
							"arn:aws:s3:::*",
							"arn:aws:s3:::bucket/*",
						},
					},
				},
			},
			doc2: Document{
				Version: "1",
				ID:      "id2",
				Statement: []StatementEntry{
					{
						Effect: "Allow",
						Action: []string{"s3:*"},
						Resource: []string{
							"arn:aws:s3:::*",
							"arn:aws:s3:::bucket/*",
						},
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.doc1.Equal(&tt.doc2)
			if got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatementEntryEqual(t *testing.T) {
	tests := []struct {
		name   string
		entry1 StatementEntry
		entry2 StatementEntry
		want   bool
	}{
		{
			name: "equal documents",
			entry1: StatementEntry{
				Effect: "Allow",
				Action: []string{"s3:*"},
				Resource: []string{
					"arn:aws:s3:::*",
					"arn:aws:s3:::bucket/*",
				},
			},
			entry2: StatementEntry{
				Effect: "Allow",
				Action: []string{"s3:*"},
				Resource: []string{
					"arn:aws:s3:::*",
					"arn:aws:s3:::bucket/*",
				},
			},
			want: true,
		},
		{
			name: "different Effects",
			entry1: StatementEntry{
				Effect: "Allow",
				Action: []string{"s3:*"},
				Resource: []string{
					"arn:aws:s3:::*",
					"arn:aws:s3:::bucket/*",
				},
			},
			entry2: StatementEntry{
				Effect: "Deny",
				Action: []string{"s3:*"},
				Resource: []string{
					"arn:aws:s3:::*",
					"arn:aws:s3:::bucket/*",
				},
			},
			want: false,
		},
		{
			name: "different Actions",
			entry1: StatementEntry{
				Effect: "Allow",
				Action: []string{"s3:*"},
				Resource: []string{
					"arn:aws:s3:::*",
					"arn:aws:s3:::bucket/*",
				},
			},
			entry2: StatementEntry{
				Effect: "Allow",
				Action: []string{"s3:GetObject"},
				Resource: []string{
					"arn:aws:s3:::*",
					"arn:aws:s3:::bucket/*",
				},
			},
			want: false,
		},
		{
			name: "different length Actions",
			entry1: StatementEntry{
				Effect: "Allow",
				Action: []string{},
				Resource: []string{
					"arn:aws:s3:::*",
					"arn:aws:s3:::bucket/*",
				},
			},
			entry2: StatementEntry{
				Effect: "Allow",
				Action: []string{"s3:GetObject"},
				Resource: []string{
					"arn:aws:s3:::*",
					"arn:aws:s3:::bucket/*",
				},
			},
			want: false,
		},
		{
			name: "different Resource",
			entry1: StatementEntry{
				Effect: "Allow",
				Action: []string{"s3:*"},
				Resource: []string{
					"arn:aws:s3:::*",
					"arn:aws:s3:::bucket/*",
				},
			},
			entry2: StatementEntry{
				Effect: "Allow",
				Action: []string{"s3:*"},
				Resource: []string{
					"arn:aws:s3:::example-bucket/*",
					"arn:aws:s3:::bucket/*",
				},
			},
			want: false,
		},
		{
			name: "different length Resource",
			entry1: StatementEntry{
				Effect:   "Allow",
				Action:   []string{"s3:*"},
				Resource: []string{},
			},
			entry2: StatementEntry{
				Effect: "Allow",
				Action: []string{"s3:*"},
				Resource: []string{
					"arn:aws:s3:::example-bucket/*",
					"arn:aws:s3:::bucket/*",
				},
			},
			want: false,
		},
		{
			name: "different Principal",
			entry1: StatementEntry{
				Effect: "Allow",
				Action: []string{"s3:*"},
				Resource: []string{
					"arn:aws:s3:::*",
					"arn:aws:s3:::bucket/*",
				},
				Principal: PrincipalEntry{
					AWS: []string{"arn:aws:iam::123456789012:user/csm"},
				},
			},
			entry2: StatementEntry{
				Effect: "Allow",
				Action: []string{"s3:*"},
				Resource: []string{
					"arn:aws:s3:::*",
					"arn:aws:s3:::bucket/*",
				},
				Principal: PrincipalEntry{
					AWS: []string{"arn:aws:iam::123456789012:user/admin"},
				},
			},
			want: false,
		},

		{
			name: "different Sid",
			entry1: StatementEntry{
				Effect: "Allow",
				Action: []string{"s3:*"},
				Resource: []string{
					"arn:aws:s3:::*",
					"arn:aws:s3:::bucket/*",
				},
				Principal: PrincipalEntry{
					AWS: []string{"arn:aws:iam::123456789012:user/csm"},
				},
				Sid: "statementID",
			},
			entry2: StatementEntry{
				Effect: "Allow",
				Action: []string{"s3:*"},
				Resource: []string{
					"arn:aws:s3:::*",
					"arn:aws:s3:::bucket/*",
				},
				Principal: PrincipalEntry{
					AWS: []string{"arn:aws:iam::123456789012:user/csm"},
				},
				Sid: "statementID2",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.entry1.Equal(&tt.entry2)
			if got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrincipalEntryEqual(t *testing.T) {
	tests := []struct {
		name   string
		entry1 PrincipalEntry
		entry2 PrincipalEntry
		want   bool
	}{
		{
			name: "equal Principal entry",
			entry1: PrincipalEntry{
				AWS: []string{"arn:aws:iam::123456789012:user/csm"},
			},
			entry2: PrincipalEntry{
				AWS: []string{"arn:aws:iam::123456789012:user/csm"},
			},
			want: true,
		},
		{
			name: "different AWS",
			entry1: PrincipalEntry{
				AWS: []string{"arn:aws:iam::123456789012:user/csm"},
			},
			entry2: PrincipalEntry{
				AWS: []string{"arn:aws:iam::123456789012:user/admin"},
			},
			want: false,
		},
		{
			name: "different length AWS",
			entry1: PrincipalEntry{
				AWS: []string{},
			},
			entry2: PrincipalEntry{
				AWS: []string{"arn:aws:iam::123456789012:user/admin"},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.entry1.Equal(&tt.entry2)
			if got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
