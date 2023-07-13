package policy

import "encoding/json"

type StatementEntry struct {
	Effect    string
	Action    []string
	Resource  []string
	Principal PrincipalEntry
	Sid       string
}

type PrincipalEntry struct {
	AWS []string
}

type PolicyDocument struct {
	Version   string
	Id        string
	Statement []StatementEntry
}

// To JSON to string.
func (p *PolicyDocument) ToJSON() (string, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// JSON to PolicyDocument.
func NewFromJSON(jsonString string) (PolicyDocument, error) {
	p := PolicyDocument{}
	err := json.Unmarshal([]byte(jsonString), &p)
	if err != nil {
		return PolicyDocument{}, err
	}

	return p, nil
}

// Check equality between documents.
func (p *PolicyDocument) Equal(p2 *PolicyDocument) bool {
	if p.Version != p2.Version {
		return false
	}

	if p.Id != p2.Id {
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

	if !s.Principal.Equal(&s2.Principal) {
		return false
	}

	if s.Sid != s2.Sid {
		return false
	}

	return true
}

// Check equality between principals.
func (p *PrincipalEntry) Equal(p2 *PrincipalEntry) bool {
	if len(p.AWS) != len(p2.AWS) {
		return false
	}

	for i, a := range p.AWS {
		if a != p2.AWS[i] {
			return false
		}
	}

	return true
}
