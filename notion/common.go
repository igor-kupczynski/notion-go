package notion

// Annotation contains style information which applies to the whole rich text object.
//
// See https://developers.notion.com/reference/rich-text#annotations
type Annotations struct {
	Bold          bool   `json:"bold,omitempty"`
	Italic        bool   `json:"italic,omitempty"`
	Strikethrough bool   `json:"strikethrough,omitempty"`
	Underline     bool   `json:"underline,omitempty"`
	Code          bool   `json:"code,omitempty"`
	Color         string `json:"color,omitempty"`
}

// RichText objects combine a text content with syle information
//
// See https://developers.notion.com/reference/rich-text
type RichText struct {
	PlainText   string      `json:"plain_text,omitempty"`
	Href        string      `json:"href,omitempty"`
	Annotations Annotations `json:"annotations,omitempty"`
	Type        string      `json:"type,omitempty"`
	Content     string      `json:"content,omitempty"`
	// TODO: links
	// TODO: mentions
	// TODO: equations
}
