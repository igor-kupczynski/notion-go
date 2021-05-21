package notion

// Page represents the properties of a single page
//
// See also https://developers.notion.com/reference/page
type Page struct {
	Object         string                   `json:"object,omitempty"`
	ID             string                   `json:"id,omitempty"`
	CreatedTime    string                   `json:"created_time,omitempty"`
	LastEditedTime string                   `json:"last_edited_time,omitempty"`
	Parent         Parent                   `json:"parent"`
	Archived       bool                     `json:"archived,omitempty"`
	Properties     map[string]PropertyValue `json:"properties,omitempty"`
}

// Parent points to a page parent
//
// See also https://developers.notion.com/reference/page#database-parent
type Parent struct {
	Type       string `json:"type,omitempty"`
	DatabaseID string `json:"database_id,omitempty"`
	PageID     string `json:"page_id,omitempty"`
}

// PropertyValue describes the identifier, type, and value of a page property
//
// See also https://developers.notion.com/reference/page#all-property-values
type PropertyValue struct {
	ID             string                     `json:"id,omitempty"`
	Type           string                     `json:"type,omitempty"`
	Title          []RichText                 `json:"title,omitempty"`
	RichText       []RichText                 `json:"rich_text,omitempty"`
	Number         int                        `json:"number,omitempty"`
	Select         *SelectPropertyValue       `json:"select,omitempty"`
	MultiSelect    []MultiSelectPropertyValue `json:"multi_select,omitempty"`
	Checkbox       bool                       `json:"checkbox,omitempty"`
	CreatedTime    string                     `json:"created_time,omitempty"`
	LastEditedTime string                     `json:"last_edited_time,omitempty"`
	// TODO: add the other property types
}

// SelectPropertyValue represents the value of a select property
//
// See also https://developers.notion.com/reference/page#select-property-values
type SelectPropertyValue struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
}

// MultiSelectPropertyValue represents the value of a select property
//
// See also https://developers.notion.com/reference/page#multi-select-option-values
type MultiSelectPropertyValue struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
}
