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
	return WithCustomHttpClient(token, http.DefaultClient)
}

// WithCustomHttpClient creates a Service using the custom http.Client
func WithCustomHttpClient(token string, httpClient *http.Client) *Service {
	rt := httpClient.Transport
	if rt == nil {
		rt = http.DefaultTransport
	}
	httpClient.Transport = &transport{
		token:    token,
		delegate: rt,
	}
	return &Service{
		client: client.New(
			httpClient,
			client.Options{RootURL: root},
		),
	}
}

type transport struct {
	token    string
	delegate http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	r := req.Clone(req.Context())
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %v", t.token))
	r.Header.Add("Notion-Version", version)
	return t.delegate.RoundTrip(r)
}
