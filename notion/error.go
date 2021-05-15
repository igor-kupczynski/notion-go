package notion

import (
	"fmt"
)

// ApplicationError represents an error returned by the Notion API server.
//
// It is an application layer http error.
// See https://developers.notion.com/reference/errors
type ApplicationError struct {
	HttpStatusCode int
	Code           string `json:"code,omitempty"`
	Message        string `json:"message,omitempty"`
}

func (e ApplicationError) Error() string {
	return fmt.Sprintf("%d %s [%s]", e.HttpStatusCode, e.Code, e.Message)
}

// TransportError represents an error while making the buildRequest to the API server
//
// It's on transport (or lower) layer error.
type TransportError struct {
	URL   string
	Inner error
}

func (e TransportError) Error() string {
	return fmt.Sprintf("buildRequest to %s failed: %v", e.URL, e.Inner)
}

func (e TransportError) Unwrap() error {
	return e.Inner
}

// ClientError represents a client-side error, i.e. before the client made the buildRequest or after it received the response
type ClientError struct {
	Reason string
	Inner  error
}

func (e ClientError) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("client error: %s: %v", e.Reason, e.Inner)
	}
	return fmt.Sprintf("client error: %s", e.Reason)
}

func (e ClientError) Unwrap() error {
	return e.Inner
}
