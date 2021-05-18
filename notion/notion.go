package notion

import (
	"fmt"
	"net/http"

	"notion-go/client"
)

const version = "2021-05-13"
const root = "https://api.notion.com/v1"

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
