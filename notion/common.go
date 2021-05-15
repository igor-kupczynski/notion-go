package notion

import (
	"encoding/json"
	"log"
	"net/http"
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
	PlainText   string      `json:"plain_text,omitempty"`
	Href        string      `json:"href,omitempty"`
	Annotations Annotations `json:"annotations,omitempty"`
	Type        string      `json:"type,omitempty"`
	Content     string      `json:"content,omitempty"`
	// TODO: links
	// TODO: mentions
	// TODO: equations
}

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

func decodeResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		var apiErr ApplicationError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			log.Printf("can't decode the response: %v", err)
		}
		apiErr.HttpStatusCode = resp.StatusCode
		return apiErr
	}

	err := json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return ClientError{Reason: "can't parse the response", Inner: err}
	}
	return nil
}
