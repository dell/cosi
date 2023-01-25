package model

// TagSet is a list of tags
type TagSet struct {
	Tags []Tag `json:"Tag,omitempty" xml:"Tag,omitempty"`
}

// Tag is an arbitrary piece of metadata applied to an object
type Tag struct {
	// Key is the tag label or name of the tag
	Key string `json:"Key" xml:"Key"`

	// Value is the value of the tag
	Value string `json:"Value" xml:"Value"`
}
