package notion

import (
	"context"
	"net/http"
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
	Results    []Database `json:"results,omitempty"`
}

// ListDatabases lists all databases shared with the authenticated integration.
//
// See https://developers.notion.com/reference/get-databases
func (s *Service) ListDatabases(ctx context.Context, page Pagination) (*DatabaseList, error) {
	dbs := &DatabaseList{}
	apiErr := &Error{}
	if err := s.client.Do(ctx, http.MethodGet, "/databases", page.query(), nil, dbs, apiErr); err != nil {
		return nil, err
	}
	return dbs, nil
}
