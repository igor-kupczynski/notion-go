package notion

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func Test_parseErrorFromResponse(t *testing.T) {
	tests := []struct {
		name string
		resp *http.Response
		want ApiServerError
	}{
		{
			name: "should parse a notion API error",
			resp: response(
				400,
				`{"code":"invalid_request","message":"This request is not supported."}`,
			),
			want: ApiServerError{
				HttpStatusCode: 400,
				Code:           "invalid_request",
				Message:        "This request is not supported.",
			},
		},
		{
			name: "should return only the code of a non-conforming error",
			resp: response(
				503,
				`#yolo`,
			),
			want: ApiServerError{
				HttpStatusCode: 503,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseErrorFromResponse(tt.resp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseErrorFromResponse() = %v, want %v", got, tt.want)
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
