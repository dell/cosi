package model

// SearchMetadata is the custom metadata for enabled for querying on the
// bucket
type SearchMetadata struct {
	// MaxKeys indicates the maximum number of metadata search keys for the
	// bucket
	MaxKeys int `json:"maxKeys" xml:"maxKeys"`

	// Enabled indicates if search metadata is enabled for the bucket
	Enabled bool `json:"isEnabled" xml:"isEnabled"`

	// Metadata defines the fields that can be searched upon
	Metadata `json:"metadata" xml:"metadata"`
}

// Metadata defines the fields that can be searched upon
type Metadata struct {
	// Type is the metadata key type
	Type string `json:"type" xml:"type"`

	// Name is the metadata key name
	Name string `json:"name" xml:"name"`

	// Datatype is the data type of the metadata value
	Datatype string `json:"datatype" xml:"datatype"`
}
