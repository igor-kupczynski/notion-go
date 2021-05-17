package notion

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"notion-go/client"
)

const version = "2021-05-13"
const root = "https://api.notion.com/v1"

// Client is the notion API client
type Client struct {
	// Token to use to connect to notion
	Token string

	httpClient http.Client
}

func (c *Client) makeRequest( // TODO: test
	method string,
	url string,
	query map[string]string,
	payload io.Reader,
	target interface{},
) error {
	req, err := buildRequest(method, url, query, payload)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.Token))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return TransportError{URL: req.URL.String(), Inner: err}
	}

	return decodeResponse(resp, target)
}

func buildRequest(method string, url string, query map[string]string, payload io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, root+url, payload)
	if err != nil {
		return nil, ClientError{Reason: "can't create a request", Inner: err}
	}

	if len(query) > 0 {
		q := req.URL.Query()
		for k, v := range query {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	req.Header.Set("Notion-Version", version)

	return req, nil
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

// Service is the facade for the notion API
type Service struct {
	client *client.Client
	token  string
}

// New creates a Service
func New(token string) *Service {
	// TODO: allow to customize http client
	t := transport{
		token: token,
	}
	return &Service{
		client: client.New(
			&http.Client{Transport: &t},
			client.Options{RootURL: root},
		),
	}
}

type transport struct {
	token string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	r := req.Clone(req.Context())
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %v", t.token))
	r.Header.Add("Notion-Version", version)
	return http.DefaultTransport.RoundTrip(r)
}
