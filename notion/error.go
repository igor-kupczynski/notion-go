package notion

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// ApiServerError represents an error returned by the Notion API server
//
// See https://developers.notion.com/reference/errors
type ApiServerError struct {
	HttpStatusCode int
	Code           string `json:"code,omitempty"`
	Message        string `json:"message,omitempty"`
}

func (n ApiServerError) Error() string {
	return fmt.Sprintf("%d %s [%s]", n.HttpStatusCode, n.Code, n.Message)
}

func parseErrorFromResponse(resp *http.Response) ApiServerError {
	var apiErr ApiServerError
	if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
		log.Printf("Can't decode the response: %v", err)
	}
	apiErr.HttpStatusCode = resp.StatusCode
	return apiErr
}
