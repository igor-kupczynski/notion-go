package notion

import (
	"context"
	"fmt"
	"net/http"
)

// Database represents a notion database
//
// See https://developers.notion.com/reference/database
type Database struct {
	Object         string              `json:"object,omitempty"`
	ID             string              `json:"id,omitempty"`
	CreatedTime    string              `json:"created_time,omitempty"`
	LastEditedTime string              `json:"last_edited_time,omitempty"`
	Title          []RichText          `json:"title,omitempty"`
	Properties     map[string]Property `json:"properties,omitempty"`
}

// PageList is a response to the query database endpoint
//
// See https://developers.notion.com/reference/post-database-query
// See https://developers.notion.com/reference/page
// See https://developers.notion.com/reference/pagination
type PageList struct {
	Object     string
	Results    []Page `json:"results,omitempty"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more,omitempty"`
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

// Filter describes conditions on page property values to include in the results from a database query
//
// See also https://developers.notion.com/reference/post-database-query#post-database-query-filter
type Filter struct {
	Property string                   `json:"property,omitempty"`
	Checkbox *CheckboxFilterCondition `json:"checkbox,omitempty"`
	// TODO: add more filter types
}

// CheckboxFilterCondition applies to database properties of type "checkbox".
//
// See also https://developers.notion.com/reference/post-database-query#checkbox-filter-condition
type CheckboxFilterCondition struct {
	Equals       bool `json:"equals,omitempty"`
	DoesNotEqual bool `json:"does_not_equal,omitempty"`
}

const (
	SortAsc  = "ascending"
	SortDesc = "descending"
)

// Sort objects describe the order of database query results
//
// See also https://developers.notion.com/reference/post-database-query (bottom of the page)
type Sort struct {
	Property  string `json:"property,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Direction string `json:"direction,omitempty"`
}

// RetrieveDatabase retrieves a Database object using the ID specified
//
// See https://developers.notion.com/reference/get-database
func (s *Service) RetrieveDatabase(ctx context.Context, databaseID string) (*Database, error) {
	db := &Database{}
	apiErr := &Error{}
	if err := s.client.Do(ctx, http.MethodGet, fmt.Sprintf("/databases/%s", databaseID), nil, nil, db, apiErr); err != nil {
		return nil, err
	}
	return db, nil
}

// QueryDatabase returns a list of pages from the given database
//
// The pages are filtered per given criteria.
//
// See https://developers.notion.com/reference/post-database-query#post-database-query-filter
func (s *Service) QueryDatabase(
	ctx context.Context,
	databaseID string,
	filter *Filter,
	sorts []Sort,
	pagination *Pagination,
) (*PageList, error) {
	type Payload struct {
		Filter      *Filter `json:"filter,omitempty"`
		Sorts       []Sort  `json:"sorts,omitempty"`
		StartCursor *string `json:"start_cursor,omitempty"`
		PageSize    int     `json:"page_size,omitempty"`
	}
	payload := &Payload{
		Filter: filter,
		Sorts:  sorts,
	}
	if pagination != nil {
		if pagination.StartCursor != "" {
			payload.StartCursor = &pagination.StartCursor
		}
		payload.PageSize = pagination.PageSize
	}
	pages := &PageList{}
	apiErr := &Error{}
	if err := s.client.Do(
		ctx,
		http.MethodPost,
		fmt.Sprintf("/databases/%s/query", databaseID),
		nil,
		payload,
		pages,
		apiErr,
	); err != nil {
		return nil, err
	}
	return pages, nil
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
