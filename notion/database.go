package notion

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
	Results    []Database `json:"results,omitempty"`
}

// ListDatabases lists all databases shared with the authenticated integration.
//
// See https://developers.notion.com/reference/get-databases
func (c *Client) ListDatabases(page Pagination) (*DatabaseList, error) {
	r, err := c.buildRequest("GET", "/databases", page.query(), nil)
	if err != nil {
		return nil, ClientError{Reason: "can't create a request", Inner: err} // TODO: test
	}

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, TransportError{URL: r.URL.String(), Inner: err} // TODO: test
	}

	dbs := &DatabaseList{}
	if err = decodeResponse(resp, dbs); err != nil {
		return nil, err
	}
	return dbs, nil
}
