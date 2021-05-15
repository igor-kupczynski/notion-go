package notion

import (
	"fmt"
	"io"
	"net/http"
)

const version = "2021-05-13"
const root = "https://api.notion.com/v1/"

// Client is the notion API client
type Client struct {
	// Token to use to connect to notion
	Token string

	httpClient http.Client
}

func (c *Client) request(method string, url string, payload io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, root+url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.Token))
	req.Header.Set("Notion-Version", version)

	return req, nil
}
