package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
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

type body struct {
	Body string `json:"body,omitempty"`
}

type success struct {
	Success string `json:"success,omitempty"`
}

type failure struct {
	Failure string `json:"failure,omitempty"`
}

func TestClient_Do(t *testing.T) {
	type args struct {
		method string
		path   string
		query  map[string]string
		body   interface{}
	}
	tests := []struct {
		name              string
		response          RequestToResponse
		opts              Options
		args              args
		wantTargetSuccess success
		wantTargetFailure failure
		wantErrMsg        string
		wantRequest       func(t *testing.T, r *http.Request)
	}{
		{
			name: "should decode successful response into the target",
			response: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"success":"yes"}`)),
				}, nil
			},
			wantTargetSuccess: success{Success: "yes"},
		},
		{
			name: "should decode failed response into the target and return an ApplicationError",
			response: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 500,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"failure":"internal server error"}`)),
				}, nil
			},
			wantTargetFailure: failure{Failure: "internal server error"},
			wantErrMsg:        "application error: &{internal server error}",
		},
		{
			name: "should fail with LocalError when request can't be created",
			args: args{
				method: "🦄",
			},
			wantErrMsg: `local error: failed to create GET request: net/http: invalid method "🦄"`,
		},
		{
			name: "should fail with LocalError when successful response can't be decoded",
			response: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`#yolo`)),
				}, nil
			},
			wantErrMsg: `local error: can't decode successful response: invalid character '#' looking for beginning of value`,
		},
		{
			name: "should fail with LocalError when failure response can't be decoded",
			response: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 500,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`#yolo`)),
				}, nil
			},
			wantErrMsg: `local error: can't decode failure response: invalid character '#' looking for beginning of value`,
		},
		{
			name: "should fail with TransportError when connection fails",
			response: func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("connection error")
			},
			args: args{
				method: http.MethodPost,
				path:   "/foo",
			},
			wantErrMsg: `transport error: request to /foo failed: Post "/foo": connection error`,
		},
		{
			name: "should use all the arguments to build the request",
			response: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"success":"yes"}`)),
				}, nil
			},
			opts: Options{RootURL: "https://api.example.com:9876"},
			args: args{
				method: http.MethodDelete,
				path:   "/foo",
				query:  map[string]string{"userId": "horse"},
				body:   &body{Body: "body"},
			},
			wantTargetSuccess: success{Success: "yes"},
			wantRequest: func(t *testing.T, r *http.Request) {
				wantMethod := http.MethodDelete
				if r.Method != wantMethod {
					t.Errorf("r.Method = %s, want %s", r.Method, wantMethod)
				}
				wantURL := "https://api.example.com:9876/foo?userId=horse"
				if r.URL.String() != wantURL {
					t.Errorf("r.URL = %s, want %s", r.URL.String(), wantURL)
				}
				wantBody := body{Body: "body"}
				var gotBody body
				_ = json.NewDecoder(r.Body).Decode(&gotBody)
				if !reflect.DeepEqual(wantBody, gotBody) {
					t.Errorf("r.Body = %v, want %v", gotBody, wantBody)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient, capturedRequest := RequestCapturingMockHttpClient(tt.response)

			c := &Client{
				httpClient: httpClient,
				opts:       &tt.opts,
			}

			gotTargetSuccess := success{}
			gotTargetFailure := failure{}
			err := c.Do(
				context.Background(),
				tt.args.method,
				tt.args.path,
				tt.args.query,
				tt.args.body,
				&gotTargetSuccess,
				&gotTargetFailure,
			)

			if tt.wantErrMsg != "" {
				if err == nil {
					err = fmt.Errorf("no error")
				}
				if !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErrMsg)
				}
			} else if err != nil {
				t.Errorf("Do() error = %v, wantErr <nil>", err)
			}
			if !reflect.DeepEqual(gotTargetSuccess, tt.wantTargetSuccess) {
				t.Errorf("Do() targetSuccess = %v, want %v", gotTargetSuccess, tt.wantTargetSuccess)
			}
			if !reflect.DeepEqual(gotTargetFailure, tt.wantTargetFailure) {
				t.Errorf("Do() gotTargetFailure = %v, want %v", gotTargetFailure, tt.wantTargetFailure)
			}
			if tt.wantRequest != nil {
				tt.wantRequest(t, capturedRequest)
			}
		})
	}
}
