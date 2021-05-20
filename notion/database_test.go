package notion

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// RequestToResponse is a function which given the request produces a response or an error
type RequestToResponse func(req *http.Request) (*http.Response, error)

// RoundTrip method to implement http.RoundTripper interface
func (f RequestToResponse) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// Create a mock http.Client which instead of making an http call will use provided function to provide a response
func RequestCapturingMockHttpClient(f RequestToResponse) (*http.Client, *http.Request) {
	var capture http.Request
	return &http.Client{
		Transport: RequestToResponse(func(req *http.Request) (*http.Response, error) {
			capture = *req.Clone(req.Context())
			return f(req)
		}),
	}, &capture
}

func TestService_RetrieveDatabase(t *testing.T) {
	tests := []struct {
		name           string
		databaseID     string
		respStatusCode int
		respBody       string
		wantPath       string
		wantDatabase   *Database
		wantErrMsg     string
	}{
		{
			name:           "should retrieve a database",
			databaseID:     "e65ccf14-e13b-48d1-a6d1-b14cd84c4bed",
			respStatusCode: 200,
			respBody: `{
			  "object": "database",
			  "id": "e65ccf14-e13b-48d1-a6d1-b14cd84c4bed",
			  "created_time": "2021-05-15T07:29:53.878Z",
			  "last_edited_time": "2021-05-20T09:19:00.000Z",
			  "title": [
				{
				  "type": "text",
				  "text": {
					"content": "Task List 5132beee",
					"link": null
				  },
				  "annotations": {
					"bold": false,
					"italic": false,
					"strikethrough": false,
					"underline": false,
					"code": false,
					"color": "default"
				  },
				  "plain_text": "Task List 5132beee",
				  "href": null
				}
			  ],
			  "properties": {
				"Date Created": {
				  "id": "'Y6<",
				  "type": "created_time",
				  "created_time": {}
				},
				"Date Edited": {
				  "id": "M[oR",
				  "type": "last_edited_time",
				  "last_edited_time": {}
				},
				"Needs ☕️?": {
				  "id": "RRGi",
				  "type": "checkbox",
				  "checkbox": {}
				},
				"Tag": {
				  "id": "UHT}",
				  "type": "multi_select",
				  "multi_select": {
					"options": [
					  {
						"id": "0e8b9aa9-b1c5-4964-812d-207d0aec09cf",
						"name": "go",
						"color": "brown"
					  },
					  {
						"id": "fc51b97d-458a-4bcc-8974-914b54afe2d6",
						"name": "software-engineering",
						"color": "default"
					  },
					  {
						"id": "c8a7c473-9b5c-4983-badc-238c85637d7e",
						"name": "skiing",
						"color": "pink"
					  },
					  {
						"id": "f9d882cc-85d4-4c73-adf8-0a83e839510d",
						"name": "outdoors",
						"color": "orange"
					  }
					]
				  }
				},
				"Status": {
				  "id": "^OE@",
				  "type": "select",
				  "select": {
					"options": [
					  {
						"id": "1",
						"name": "To Do",
						"color": "red"
					  },
					  {
						"id": "2",
						"name": "Doing",
						"color": "yellow"
					  },
					  {
						"id": "3",
						"name": "Done 🙌",
						"color": "green"
					  }
					]
				  }
				},
				"Name": {
				  "id": "title",
				  "type": "title",
				  "title": {}
				}
			  }
			}`,
			wantPath: "/v1/databases/e65ccf14-e13b-48d1-a6d1-b14cd84c4bed",
			wantDatabase: &Database{
				Object:         "database",
				ID:             "e65ccf14-e13b-48d1-a6d1-b14cd84c4bed",
				CreatedTime:    "2021-05-15T07:29:53.878Z",
				LastEditedTime: "2021-05-20T09:19:00.000Z",
				Title: []RichText{
					{
						Type: "text",
						Text: &Text{Content: "Task List 5132beee"},
						Annotations: &Annotations{
							Color: "default",
						},
						PlainText: "Task List 5132beee",
					},
				},
				Properties: map[string]Property{
					"Date Created": {
						ID:          "'Y6<",
						Type:        "created_time",
						CreatedTime: &CreatedTimeProperty{},
					},
					"Date Edited": {
						ID:             "M[oR",
						Type:           "last_edited_time",
						LastEditedTime: &LastEditedTimeProperty{},
					},
					"Needs ☕️?": {
						ID:       "RRGi",
						Type:     "checkbox",
						Checkbox: &CheckboxProperty{},
					},
					"Tag": {
						ID:   "UHT}",
						Type: "multi_select",
						MultiSelect: &MultiSelectProperty{
							Options: []MultiSelectOption{
								{
									ID:    "0e8b9aa9-b1c5-4964-812d-207d0aec09cf",
									Name:  "go",
									Color: "brown",
								},
								{
									ID:    "fc51b97d-458a-4bcc-8974-914b54afe2d6",
									Name:  "software-engineering",
									Color: "default",
								},
								{
									ID:    "c8a7c473-9b5c-4983-badc-238c85637d7e",
									Name:  "skiing",
									Color: "pink",
								},
								{
									ID:    "f9d882cc-85d4-4c73-adf8-0a83e839510d",
									Name:  "outdoors",
									Color: "orange",
								},
							},
						},
					},
					"Status": {
						ID:   "^OE@",
						Type: "select",
						Select: &SelectProperty{
							Options: []SelectOption{
								{
									ID:    "1",
									Name:  "To Do",
									Color: "red",
								},
								{
									ID:    "2",
									Name:  "Doing",
									Color: "yellow",
								},
								{
									ID:    "3",
									Name:  "Done 🙌",
									Color: "green",
								},
							},
						},
					},
					"Name": {
						ID:    "title",
						Type:  "title",
						Title: &TitleProperty{},
					},
				},
			},
		},
		{
			name:           "should parse an error",
			databaseID:     "not-uuid",
			respStatusCode: 400,
			respBody: `{
			  "object": "error",
			  "status": 400,
			  "code": "validation_error",
			  "message": "The provided database ID is not a valid Notion UUID: e65ccf14-e13b-48d1-a6d1-b14cd84c4be."
			}`,
			wantPath:   "/v1/databases/not-uuid",
			wantErrMsg: "application error: &{validation_error The provided database ID is not a valid Notion UUID: e65ccf14-e13b-48d1-a6d1-b14cd84c4be.}",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			httpClient, capturedRequest := RequestCapturingMockHttpClient(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: tt.respStatusCode,
					Body:       ioutil.NopCloser(bytes.NewBufferString(tt.respBody)),
				}, nil
			})
			service := WithCustomHttpClient("token", httpClient)

			gotDB, gotErr := service.RetrieveDatabase(context.Background(), tt.databaseID)

			gotPath := capturedRequest.URL.Path
			if gotPath != tt.wantPath {
				t.Errorf("path = %v, want %v", gotPath, tt.wantPath)

			}
			if tt.wantErrMsg != "" {
				if gotErr == nil {
					gotErr = fmt.Errorf("no error")
				}
				if !strings.Contains(gotErr.Error(), tt.wantErrMsg) {
					t.Errorf("RetrieveDatabase() error = %v, wantErr %v", gotErr, tt.wantErrMsg)
				}
			} else if gotErr != nil {
				t.Errorf("RetrieveDatabase() error = %v, wantErr <nil>", gotErr)
			}
			if diff := cmp.Diff(tt.wantDatabase, gotDB); diff != "" {
				t.Errorf("RetrieveDatabase() mismatch (-want +got):\n%s", diff)
			}

		})
	}
}

func TestService_ListDatabases_Integration(t *testing.T) {
	token := os.Getenv("NOTION_TOKEN")
	if token == "" {
		t.Skip("set NOTION_TOKEN to run this test")
	}

	wantTitle := "Task List 5132beee"

	s := New(token)

	// Get the list of the databases, list them one-by-one to exercise the pagination code path
	var got []Database

	page := Pagination{PageSize: 1}
	for {
		result, err := s.ListDatabases(context.Background(), page)
		if err != nil {
			t.Errorf("ListDatabases(%v) error = %v", page, err)
			return
		}
		for _, db := range result.Results {
			got = append(got, db)
		}
		if len(result.Results) > page.PageSize {
			t.Errorf("ListDatabases(%v) got too many items [%s]", page, renderTitles(result.Results))
			return
		}
		if !result.HasMore {
			break
		}
		page.StartCursor = result.NextCursor
	}

	// Check if there are at least two databases
	if len(got) < 2 {
		t.Errorf("Expected at least 2 databases, got [%s]", renderTitles(got))
		return
	}

	// Check if it contains the one we know it has to contain
	var taskList *Database
FindDB:
	for _, db := range got {
		for _, titlet := range db.Title {
			if strings.Contains(titlet.PlainText, wantTitle) {
				taskList = &db
				break FindDB
			}
		}
	}

	// If not then lets print what we have to make it easier to debug
	if taskList == nil {
		t.Errorf("Test DB [%s] not found. Got: [%s]", wantTitle, renderTitles(got))
	}
}

func renderTitles(got []Database) string {
	allTitles := []string{}
	for _, db := range got {
		title := []string{}
		for _, titlet := range db.Title {
			title = append(title, titlet.PlainText)
		}
		allTitles = append(allTitles, strings.Join(title, ""))
	}
	return strings.Join(allTitles, ", ")
}
