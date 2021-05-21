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
func New(token string, trace bool) *Service {
	return WithCustomHttpClient(token, http.DefaultClient, trace)
}

// WithCustomHttpClient creates a Service using the custom http.Client
func WithCustomHttpClient(token string, httpClient *http.Client, trace bool) *Service {
	return &Service{
		client: client.New(
			httpClient,
			client.Options{
				AddHeaders: map[string]string{
					"Authorization":  fmt.Sprintf("Bearer %v", token),
					"Notion-Version": version,
				},
				RootURL: root,
				Trace:   trace,
			},
		),
	}
}
