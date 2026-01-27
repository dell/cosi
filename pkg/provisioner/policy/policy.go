// Copyright © 2023-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

// Package policy exists for handling operations on AWS Policies,
// defines structures and functions for processing and comparing policies.
package policy

import "encoding/json"

type StatementEntry struct {
	Effect    string            `json:"Effect"`
	Action    []string          `json:"Action"`
	Resource  []string          `json:"Resource"`
	Principal map[string]string `json:"Principal,omitempty"`
	Sid       string            `json:"Sid,omitempty"`
}

type Document struct {
	Version   string           `json:"Version"`
	ID        string           `json:"Id,omitempty"`
	Statement []StatementEntry `json:"Statement"`
}

// To JSON to string.
func (p *Document) ToJSON() (string, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// JSON to PolicyDocument.
func NewFromJSON(jsonString string) (Document, error) {
	p := Document{}
	if err := json.Unmarshal([]byte(jsonString), &p); err != nil {
		return Document{}, err
	}

	return p, nil
}

// Check equality between documents.
func (p *Document) Equal(p2 *Document) bool {
	if p.Version != p2.Version {
		return false
	}

	if p.ID != p2.ID {
		return false
	}

	if len(p.Statement) != len(p2.Statement) {
		return false
	}

	for i, s := range p.Statement {
		if !s.Equal(&p2.Statement[i]) {
			return false
		}
	}

	return true
}

// Check equality between statements.
func (s *StatementEntry) Equal(s2 *StatementEntry) bool {
	if s.Effect != s2.Effect {
		return false
	}

	if len(s.Action) != len(s2.Action) {
		return false
	}

	for i, a := range s.Action {
		if a != s2.Action[i] {
			return false
		}
	}

	if len(s.Resource) != len(s2.Resource) {
		return false
	}

	for i, r := range s.Resource {
		if r != s2.Resource[i] {
			return false
		}
	}

	return true
}
