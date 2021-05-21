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

func TestService_QueryDatabase(t *testing.T) {
	tests := []struct {
		name           string
		databaseID     string
		filter         *Filter
		sorts          []Sort
		pagination     *Pagination
		respStatusCode int
		respBody       string
		wantPath       string
		wantPayload    string
		wantResult     *PageList
		wantErrMsg     string
	}{
		{
			name:           "should query a database",
			databaseID:     "e65ccf14-e13b-48d1-a6d1-b14cd84c4bed",
			respStatusCode: 200,
			respBody: `{
			  "object": "list",
			  "results": [
				{
				  "object": "page",
				  "id": "ea8229fa-a781-4348-a154-de893e232e27",
				  "created_time": "2021-05-20T09:18:00.000Z",
				  "last_edited_time": "2021-05-20T09:19:00.000Z",
				  "parent": {
					"type": "database_id",
					"database_id": "e65ccf14-e13b-48d1-a6d1-b14cd84c4bed"
				  },
				  "archived": false,
				  "properties": {
					"Date Created": {
					  "id": "'Y6<",
					  "type": "created_time",
					  "created_time": "2021-05-20T09:18:00.000Z"
					},
					"Date Edited": {
					  "id": "M[oR",
					  "type": "last_edited_time",
					  "last_edited_time": "2021-05-20T09:19:00.000Z"
					},
					"Needs ☕️?": {
					  "id": "RRGi",
					  "type": "checkbox",
					  "checkbox": true
					},
					"Tag": {
					  "id": "UHT}",
					  "type": "multi_select",
					  "multi_select": [
						{
						  "id": "0e8b9aa9-b1c5-4964-812d-207d0aec09cf",
						  "name": "go",
						  "color": "brown"
						},
						{
						  "id": "fc51b97d-458a-4bcc-8974-914b54afe2d6",
						  "name": "software-engineering",
						  "color": "default"
						}
					  ]
					},
					"Status": {
					  "id": "^OE@",
					  "type": "select",
					  "select": {
						"id": "1",
						"name": "To Do",
						"color": "red"
					  }
					},
					"Name": {
					  "id": "title",
					  "type": "title",
					  "title": [
						{
						  "type": "text",
						  "text": {
							"content": "Write more integrations tests",
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
						  "plain_text": "Write more integrations tests",
						  "href": null
						}
					  ]
					}
				  }
				}
			  ],
			  "next_cursor": null,
			  "has_more": false
			}`,
			wantPath:    "/v1/databases/e65ccf14-e13b-48d1-a6d1-b14cd84c4bed/query",
			wantPayload: "{}",
			wantResult: &PageList{
				Object: "list",
				Results: []Page{
					{
						Object:         "page",
						ID:             "ea8229fa-a781-4348-a154-de893e232e27",
						CreatedTime:    "2021-05-20T09:18:00.000Z",
						LastEditedTime: "2021-05-20T09:19:00.000Z",
						Parent: Parent{
							Type:       "database_id",
							DatabaseID: "e65ccf14-e13b-48d1-a6d1-b14cd84c4bed",
						},
						Archived: false,
						Properties: map[string]PropertyValue{
							"Date Created": {
								ID:          "'Y6<",
								Type:        "created_time",
								CreatedTime: "2021-05-20T09:18:00.000Z",
							},
							"Date Edited": {
								ID:             "M[oR",
								Type:           "last_edited_time",
								LastEditedTime: "2021-05-20T09:19:00.000Z",
							},
							"Needs ☕️?": {
								ID:       "RRGi",
								Type:     "checkbox",
								Checkbox: true,
							},
							"Tag": {
								ID:   "UHT}",
								Type: "multi_select",
								MultiSelect: []MultiSelectPropertyValue{
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
								},
							},
							"Status": {
								ID:   "^OE@",
								Type: "select",
								Select: &SelectPropertyValue{
									ID:    "1",
									Name:  "To Do",
									Color: "red",
								},
							},
							"Name": {
								ID:   "title",
								Type: "title",
								Title: []RichText{
									{
										Type: "text",
										Text: &Text{
											Content: "Write more integrations tests",
										},
										Annotations: &Annotations{
											Bold:          false,
											Italic:        false,
											Strikethrough: false,
											Underline:     false,
											Code:          false,
											Color:         "default",
										},
										PlainText: "Write more integrations tests",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:       "should pass a payload",
			databaseID: "e65ccf14-e13b-48d1-a6d1-b14cd84c4bed",
			filter: &Filter{
				Property: "Foo",
				Checkbox: &CheckboxFilterCondition{
					Equals: true,
				},
			},
			sorts: []Sort{
				{
					Property:  "Bar",
					Direction: SortAsc,
				},
			},
			pagination: &Pagination{
				StartCursor: "qwerty",
				PageSize:    100,
			},
			respStatusCode: 200,
			respBody: `{
			  "object": "list",
			  "results": [],
			  "next_cursor": null,
			  "has_more": false
			}`,
			wantPath:    "/v1/databases/e65ccf14-e13b-48d1-a6d1-b14cd84c4bed/query",
			wantPayload: `{"filter":{"property":"Foo","checkbox":{"equals":true}},"sorts":[{"property":"Bar","direction":"ascending"}],"start_cursor":"qwerty","page_size":100}`,
			wantResult: &PageList{
				Object:  "list",
				Results: []Page{},
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
			wantPath:    "/v1/databases/not-uuid/query",
			wantPayload: "{}",
			wantErrMsg:  "application error: &{validation_error The provided database ID is not a valid Notion UUID: e65ccf14-e13b-48d1-a6d1-b14cd84c4be.}",
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

			gotDB, gotErr := service.QueryDatabase(context.Background(), tt.databaseID, tt.filter, tt.sorts, tt.pagination)

			gotPath := capturedRequest.URL.Path
			if gotPath != tt.wantPath {
				t.Errorf("path = %v, want %v", gotPath, tt.wantPath)

			}
			payload, _ := ioutil.ReadAll(capturedRequest.Body)
			gotPayload := string(payload)
			if tt.wantPayload != gotPayload {
				t.Errorf("payload = %v, want %v", gotPayload, tt.wantPayload)
			}
			if tt.wantErrMsg != "" {
				if gotErr == nil {
					gotErr = fmt.Errorf("no error")
				}
				if !strings.Contains(gotErr.Error(), tt.wantErrMsg) {
					t.Errorf("QueryDatabase() error = %v, wantErr %v", gotErr, tt.wantErrMsg)
				}
			} else if gotErr != nil {
				t.Errorf("QueryDatabase() error = %v, wantErr <nil>", gotErr)
			}
			if diff := cmp.Diff(tt.wantResult, gotDB); diff != "" {
				t.Errorf("QueryDatabase() mismatch (-want +got):\n%s", diff)
			}

		})
	}
}

func TestService_QueryDatabase_Integration(t *testing.T) {
	token := os.Getenv("NOTION_TOKEN")
	if token == "" {
		t.Skip("set NOTION_TOKEN to run this test")
	}

	s := New(token)

	result, err := s.QueryDatabase(
		context.Background(),
		"e65ccf14-e13b-48d1-a6d1-b14cd84c4bed",
		&Filter{Property: "RRGi", Checkbox: &CheckboxFilterCondition{Equals: true}},
		[]Sort{{Timestamp: "created_time", Direction: SortAsc}},
		nil,
	)
	if err != nil {
		t.Errorf("QueryDatabase error = %v", err)
		return
	}
	if len(result.Results) < 2 {
		t.Errorf("QueryDatabase expected at least two pages, got %v", result.Results)
		return
	}

	// Check if it the first two pages have the results that we expect
	wantTitle1 := "Write more integrations tests"
	gotTitle1 := result.Results[0].Properties["Name"].Title[0].PlainText
	if gotTitle1 != wantTitle1 {
		t.Errorf("Page[0] title = %v, want %v", gotTitle1, wantTitle1)
	}
	wantTitle2 := "Create an integration test workspace"
	gotTitle2 := result.Results[1].Properties["Name"].Title[0].PlainText
	if gotTitle2 != wantTitle2 {
		t.Errorf("Page[1] title = %v, want %v", gotTitle2, wantTitle2)
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
