package notion

import (
	"encoding/json"
)

// Database represents a notion database
//
// See https://developers.notion.com/reference/database
type Database struct {
	Object         string     `json:"object,omitempty"`
	ID             string     `json:"id,omitempty"`
	CreatedTime    string     `json:"created_time,omitempty"`
	LastEditedTime string     `json:"last_edited_time,omitempty"`
	Title          []RichText `json:"title,omitempty"`
	// TODO: properties
}

// DatabaseList is a response to list databases endpoint
//
// See https://developers.notion.com/reference/get-databases
// See https://developers.notion.com/reference/pagination
type DatabaseList struct {
	HasMore    bool       `json:"has_more,omitempty"`
	NextCursor string     `json:"next_cursor,omitempty"`
	Object     string     `json:"object,omitempty"`
	Results    []Database `json:"results,omitempty"`
}

// ListDatabases lists all databases shared with the authenticated integration.
//
// See https://developers.notion.com/reference/get-databases
func (c *Client) ListDatabases(page Pagination) (*DatabaseList, error) {
	r, err := c.request("GET", "/databases", page.query(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, err // TODO: introduce common error classes
	}

	defer resp.Body.Close()

	// TODO: Make a generic decode function
	if resp.StatusCode != 200 {
		return nil, parseErrorFromResponse(resp)
	}

	dbs := &DatabaseList{}
	err = json.NewDecoder(resp.Body).Decode(dbs)
	if err != nil {
		return nil, err // TODO: introduce common error classes
	}

	return dbs, nil
}
