package notion

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func Test_buildRequest(t *testing.T) {
	type args struct {
		method  string
		url     string
		query   map[string]string
		payload io.Reader
	}
	tests := []struct {
		name       string
		args       args
		assertions []func(req *http.Request) string
		wantErr    bool
	}{
		{
			name: "should build the request",
			args: args{
				method: "GET",
				url:    "/pages",
			},
			assertions: []func(req *http.Request) string{
				func(req *http.Request) string {
					if req.Method != "GET" {
						return "method = GET"
					}
					return ""
				},
				func(req *http.Request) string {
					if req.URL.String() != "https://api.notion.com/v1/pages" {
						return "URL = https://api.notion.com/v1/pages"
					}
					return ""
				},
				func(req *http.Request) string {
					if req.Header.Get("Notion-Version") != version {
						return fmt.Sprintf("Header Notion-Version = %s", version)
					}
					return ""
				},
			},
		},
		{
			name: "should include the query params",
			args: args{
				query: map[string]string{
					"foo": "bar",
					"abc": "123",
				},
			},
			assertions: []func(req *http.Request) string{
				func(req *http.Request) string {
					url := req.URL.String()
					if !strings.Contains(url, "foo=bar") || !strings.Contains(url, "abc=123") {
						return "foo=bar & abc=123 in query params"
					}
					return ""
				},
			},
		},
		{
			name: "should error on illegal parameters",
			args: args{
				method: "🦄",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := buildRequest(tt.args.method, tt.args.url, tt.args.query, tt.args.payload)
			for _, assertion := range tt.assertions {
				if msg := assertion(got); msg != "" {
					t.Errorf("buildRequest() = %v, want %s", got, msg)
				}
			}
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("buildRequest() gotErr=%v, wantErr=%t", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_decodeResponse(t *testing.T) {
	tests := []struct {
		name         string
		resp         *http.Response
		decodeTarget interface{}
		want         interface{}
		wantErr      error
	}{
		{
			name: "should decide the response into target",
			resp: response(
				200,
				`{
					"results": [
						{
							"object": "database",
							"id": "668d797c-76fa-4934-9b05-ad288df2d136",
							"title": [
								{
									"type": "text",
									"text": {
										"content": "Grocery list",
										"link": null
									},
									"annotations": {
										"bold": false,
										"italic": false,
										"strikethrough": false,
										"underline": false,
										"code": false,
										"color": ""
									},
									"plain_text": "Grocery list",
									"href": null
								}
							]
						}
					],
					"next_cursor": "MTY3NDE4NGYtZTdiYy00NzFlLWE0NjctODcxOTIyYWU3ZmM3",
					"has_more": false
				}`,
			),
			decodeTarget: &DatabaseList{},
			want: &DatabaseList{
				HasMore:    false,
				NextCursor: "MTY3NDE4NGYtZTdiYy00NzFlLWE0NjctODcxOTIyYWU3ZmM3",
				Results: []Database{
					{
						Object: "database",
						ID:     "668d797c-76fa-4934-9b05-ad288df2d136",
						Title: []RichText{
							{
								PlainText: "Grocery list",
								Type:      "text",
							},
						},
					},
				},
			},
		},
		{
			name: "should return a ClientError when the json can't be decoded",
			resp: response(
				200,
				`#yolo`,
			),
			decodeTarget: &DatabaseList{},
			wantErr: ClientError{
				Reason: "can't parse the response: invalid character '#' looking for beginning of value",
			},
		},
		{
			name: "should parse a notion API error",
			resp: response(
				400,
				`{"code":"invalid_request","message":"This buildRequest is not supported."}`,
			),
			wantErr: ApplicationError{
				HttpStatusCode: 400,
				Code:           "invalid_request",
				Message:        "This buildRequest is not supported.",
			},
		},
		{
			name: "should return the http status code of error with unexpected body",
			resp: response(
				503,
				`#yolo`,
			),
			wantErr: ApplicationError{
				HttpStatusCode: 503,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErrStr := ""
			gotErr := decodeResponse(tt.resp, tt.decodeTarget)
			if gotErr != nil {
				gotErrStr = gotErr.Error()
			}
			wantErrStr := ""
			if tt.wantErr != nil {
				wantErrStr = tt.wantErr.Error()
			}
			if gotErrStr != wantErrStr {
				t.Errorf("parseErrorFromResponse() = %v, wantErr %v", gotErr, tt.wantErr)
			}
			if gotErr == nil && !reflect.DeepEqual(tt.decodeTarget, tt.want) {
				t.Errorf("parseErrorFromResponse() decoded = %+v, want %+v", tt.decodeTarget, tt.want)
			}
		})
	}
}

func response(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body: io.NopCloser(
			bytes.NewReader(
				[]byte(body),
			),
		),
	}
}
