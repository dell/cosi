// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package config

import "fmt"
import "encoding/json"

// Credentials used for authentication to object storage provider
type Credentials struct {
	// Password for object storage provider
	Password string `json:"password" yaml:"password"`

	// Username for object storage provider
	Username string `json:"username" yaml:"username"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Credentials) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["password"]; !ok || v == nil {
		return fmt.Errorf("field password in Credentials: required")
	}
	if v, ok := raw["username"]; !ok || v == nil {
		return fmt.Errorf("field username in Credentials: required")
	}
	type Plain Credentials
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Credentials(plain)
	return nil
}

// S3 configuration
type S3 struct {
	// Endpoint of the ObjectStore S3 service
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *S3) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["endpoint"]; !ok || v == nil {
		return fmt.Errorf("field endpoint in S3: required")
	}
	type Plain S3
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = S3(plain)
	return nil
}

// Protocols supported by the connection
type Protocols struct {
	// S3 corresponds to the JSON schema field "s3".
	S3 *S3 `json:"s3,omitempty" yaml:"s3,omitempty"`
}

// TLS configuration details
type Tls struct {
	// Base64 encoded content of the clients's certificate file
	ClientCert *string `json:"client-cert,omitempty" yaml:"client-cert,omitempty"`

	// Base64 encoded content of the clients's key certificate file
	ClientKey *string `json:"client-key,omitempty" yaml:"client-key,omitempty"`

	// Controls whether a client verifies the server's certificate chain and host name
	Insecure bool `json:"insecure" yaml:"insecure"`

	// Base64 encoded content of the root certificate authority file
	RootCas *string `json:"root-cas,omitempty" yaml:"root-cas,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Tls) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["insecure"]; !ok || v == nil {
		return fmt.Errorf("field insecure in Tls: required")
	}
	type Plain Tls
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Tls(plain)
	return nil
}

// Configuration specific to the Dell ObjectScale platform
type Objectscale struct {
	// Credentials corresponds to the JSON schema field "credentials".
	Credentials Credentials `json:"credentials" yaml:"credentials"`

	// Indicates if the contents of the bucket should be emptied as part of the
	// deletion process
	EmptyBucket bool `json:"emptyBucket,omitempty" yaml:"emptyBucket,omitempty"`

	// Default, unique identifier for the single connection.
	Id string `json:"id" yaml:"id"`

	// Endpoint of the ObjectScale Gateway Internal service
	ObjectscaleGateway string `json:"objectscale-gateway" yaml:"objectscale-gateway"`

	// Endpoint of the ObjectScale ObjectStore Management Gateway service
	ObjectstoreGateway string `json:"objectstore-gateway" yaml:"objectstore-gateway"`

	// ID of objectstore retrieved from the ObjectScale Portal or directly from
	// objectstore-objmt-* pod
	ObjectstoreId string `json:"objectstoreId" yaml:"objectstoreId"`

	// Protocols corresponds to the JSON schema field "protocols".
	Protocols Protocols `json:"protocols" yaml:"protocols"`

	// Identity and Access Management (IAM) API specific field, points to the region
	// in which object storage provider is installed
	Region *string `json:"region,omitempty" yaml:"region,omitempty"`

	// Tls corresponds to the JSON schema field "tls".
	Tls Tls `json:"tls" yaml:"tls"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Objectscale) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["credentials"]; !ok || v == nil {
		return fmt.Errorf("field credentials in Objectscale: required")
	}
	if v, ok := raw["id"]; !ok || v == nil {
		return fmt.Errorf("field id in Objectscale: required")
	}
	if v, ok := raw["objectscale-gateway"]; !ok || v == nil {
		return fmt.Errorf("field objectscale-gateway in Objectscale: required")
	}
	if v, ok := raw["objectstore-gateway"]; !ok || v == nil {
		return fmt.Errorf("field objectstore-gateway in Objectscale: required")
	}
	if v, ok := raw["objectstoreId"]; !ok || v == nil {
		return fmt.Errorf("field objectstoreId in Objectscale: required")
	}
	if v, ok := raw["protocols"]; !ok || v == nil {
		return fmt.Errorf("field protocols in Objectscale: required")
	}
	if v, ok := raw["tls"]; !ok || v == nil {
		return fmt.Errorf("field tls in Objectscale: required")
	}
	type Plain Objectscale
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if v, ok := raw["emptyBucket"]; !ok || v == nil {
		plain.EmptyBucket = false
	}
	*j = Objectscale(plain)
	return nil
}

// this file contains JSON schema for Dell COSI Driver Configuration file
type ConfigSchemaJson struct {
	// List of connections to object storage platforms that can be used for object
	// storage provisioning.
	Connections []Configuration `json:"connections,omitempty" yaml:"connections,omitempty"`
}

// Configuration for single connection to object storage platform that is used for
// object storage provisioning
type Configuration struct {
	// Objectscale corresponds to the JSON schema field "objectscale".
	Objectscale *Objectscale `json:"objectscale,omitempty" yaml:"objectscale,omitempty"`
}
