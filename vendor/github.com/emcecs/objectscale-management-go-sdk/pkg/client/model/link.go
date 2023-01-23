package model

// Link is a HTTP hyperlink
type Link struct {
	// HREF is the hyperlink reference
	HREF string `json:"href" xml:"href"`

	// Rel is the relationship between the link references
	Rel string `json:"rel" xml:"rel"`
}
