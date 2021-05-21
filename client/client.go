// Package client offers a higher-level http client to make it easy to build API clients on top of it
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// LocalError represents a client-side error, i.e. client can't build the request or parse the response
type LocalError struct {
	Reason string
	Inner  error
}

func (e LocalError) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("local error: %s: %v", e.Reason, e.Inner)
	}
	return fmt.Sprintf("local error: %s", e.Reason)
}

func (e LocalError) Unwrap() error {
	return e.Inner
}

// TransportError represents an error while making the request to the server
type TransportError struct {
	URL   string
	Inner error
}

func (e TransportError) Error() string {
	return fmt.Sprintf("transport error: request to %s failed: %v", e.URL, e.Inner)
}

func (e TransportError) Unwrap() error {
	return e.Inner
}

// ApplicationError represents an error on the application layer, i.e. http status code > 2xx
type ApplicationError struct {
	v interface{}
}

func (e ApplicationError) Error() string {
	return fmt.Sprintf("application error: %v", e.v)
}

// Options can customize Client behavior
type Options struct {
	RootURL string
}

// Client is a wrapper over http.Client to make it easier to use from the notion API
type Client struct {
	httpClient *http.Client
	opts       *Options
}

// New creates a Client with provided options
func New(httpClient *http.Client, opts Options) *Client {
	return &Client{
		httpClient: httpClient,
		opts:       &opts,
	}
}

// Do issues a request with given params.
//
// In case of 2xx response decode the response body into targetSuccess.
// In case of >2xx response return ApplicationError and try to decode the body into targetFailure
// May return one of ApplicationError, LocalError, TransportError in case of a failure
func (c *Client) Do(
	ctx context.Context,
	method string,
	path string,
	query map[string]string,
	body interface{},
	targetSuccess interface{},
	targetFailure interface{},
) error {
	req, err := c.newRequest(ctx, method, path, query, body)
	if err != nil {
		return err
	}

	return c.do(req, targetSuccess, targetFailure)
}

func (c *Client) newRequest(
	ctx context.Context,
	method string,
	path string,
	query map[string]string,
	body interface{},
) (*http.Request, error) {
	buf, err := c.encode(body)
	if err != nil {
		return nil, LocalError{Reason: "failed to encode the body", Inner: err}
	}

	req, err := http.NewRequest(method, c.opts.RootURL+path, buf)
	if err != nil {
		return nil, LocalError{Reason: "failed to create GET request", Inner: err}
	}

	if len(query) > 0 {
		q := req.URL.Query()
		for k, v := range query {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	req = req.WithContext(ctx)

	return req, nil
}

func (c *Client) do(r *http.Request, targetSuccess interface{}, targetFailure interface{}) error {
	resp, err := c.httpClient.Do(r)
	if err != nil {
		return TransportError{URL: r.URL.String(), Inner: err}
	}

	defer resp.Body.Close()
	if resp.StatusCode <= 300 {
		if err := c.decode(resp, targetSuccess); err != nil {
			return LocalError{Reason: "can't decode successful response", Inner: err}
		}
		return nil
	}
	if err := c.decode(resp, targetFailure); err != nil {
		return LocalError{Reason: "can't decode failure response", Inner: err}
	}
	return ApplicationError{v: targetFailure}
}

func (c *Client) encode(v interface{}) (io.Reader, error) {
	buf, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(buf), nil
}

func (c *Client) decode(resp *http.Response, v interface{}) error {
	err := json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return err
	}
	return nil
}
