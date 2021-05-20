package notion

import (
	"strconv"
)

// Annotations contains style information which applies to the whole rich text object.
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
	Type        string       `json:"type,omitempty"`
	Text        *Text        `json:"text,omitempty"`
	Annotations *Annotations `json:"annotations,omitempty"`
	PlainText   string       `json:"plain_text,omitempty"`
	Href        string       `json:"href,omitempty"`
	Content     string       `json:"content,omitempty"`
	// TODO: links
	// TODO: mentions
	// TODO: equations
}

// Text object
//
// See https://developers.notion.com/reference/rich-text#text-objects
type Text struct {
	Content string `json:"content,omitempty"`
	// TODO: link
}

// Property represents any type of the property object
//
// See https://developers.notion.com/reference/database#database-properties
type Property struct {
	ID             string                  `json:"id,omitempty"`
	Type           string                  `json:"type,omitempty"`
	Title          *TitleProperty          `json:"title,omitempty"`
	Select         *SelectProperty         `json:"select,omitempty"`
	MultiSelect    *MultiSelectProperty    `json:"multi_select,omitempty"`
	Checkbox       *CheckboxProperty       `json:"checkbox,omitempty"`
	CreatedTime    *CreatedTimeProperty    `json:"created_time,omitempty"`
	LastEditedTime *LastEditedTimeProperty `json:"last_edited_time,omitempty"`
}

// TitleProperty represents the title property
//
// See https://developers.notion.com/reference/database#title-configuration
type TitleProperty struct{}

// SelectProperty represents the select property
//
// See https://developers.notion.com/reference/database#select-configuration
type SelectProperty struct {
	Options []SelectOption `json:"options,omitempty"`
}

// SelectOption represents the options to SelectProperty
//
// See https://developers.notion.com/reference/database#select-options
type SelectOption struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
}

// MultiSelectProperty represents the select property
//
// See https://developers.notion.com/reference/database#multi-select-configuration
type MultiSelectProperty struct {
	Options []MultiSelectOption `json:"options,omitempty"`
}

// MultiSelectOption represents the options to MultiSelectProperty
//
// See https://developers.notion.com/reference/database#multi-select-options
type MultiSelectOption struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
}

// CheckboxProperty represents the checkbox property
//
// See https://developers.notion.com/reference/database#checkbox-configuration
type CheckboxProperty struct{}

// CreatedTimeProperty represents the created time property
//
// See https://developers.notion.com/reference/database#created-time-configuration
type CreatedTimeProperty struct{}

// LastEditedTimeProperty represents the last edited time property
//
// See https://developers.notion.com/reference/database#last-edited-time-configuration
type LastEditedTimeProperty struct{}

// Pagination represents a request pagination params
//
// See https://developers.notion.com/reference/pagination
type Pagination struct {
	StartCursor string
	PageSize    int
}

func (p *Pagination) query() map[string]string {
	if p == nil {
		return nil
	}
	query := map[string]string{
		"page_size": strconv.Itoa(p.PageSize),
	}

	if p.StartCursor != "" {
		query["start_cursor"] = p.StartCursor
	}

	return query
}

// Error represents an error returned by the API
//
// See https://developers.notion.com/reference/errors
type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
