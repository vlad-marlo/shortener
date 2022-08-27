package httpserver

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, []byte) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, respBody
}

func TestServer_HandleURLGetAndCreate(t *testing.T) {
	type args struct {
		urlPath    string
		urlToShort string
	}
	type want struct {
		wantInternalServerError bool
		status                  int
	}
	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		name string

		args args
		want want
	}{
		{
			name: "positive case #1",
			args: args{
				urlPath:    "/",
				urlToShort: "google.com",
			},
			want: want{
				wantInternalServerError: false,
				status:                  http.StatusCreated,
			},
		},
		{
			name: "positive case #2",
			args: args{
				urlPath:    "/",
				urlToShort: "ya.ru",
			},
			want: want{
				wantInternalServerError: false,
				status:                  http.StatusCreated,
			},
		},
		{
			name: "incorrect target case",
			args: args{
				urlPath:    "/jkljk/",
				urlToShort: "yandex.ru",
			},
			want: want{
				wantInternalServerError: true,
				status:                  http.StatusNotFound,
			},
		},
		{
			name: "uncorrect url to short",
			args: args{
				urlPath:    "/",
				urlToShort: "hlt v.org",
			},
			want: want{
				wantInternalServerError: true,
				status:                  http.StatusBadRequest,
			},
		},
		{
			name: "empty data",
			args: args{
				urlPath:    "/",
				urlToShort: "",
			},
			want: want{
				wantInternalServerError: true,
				status:                  http.StatusBadRequest,
			},
		},
	}

	s := New(NewConfig("", "inmemory"))
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, url := testRequest(
				t,
				ts,
				http.MethodPost,
				tt.args.urlPath,
				strings.NewReader(tt.args.urlToShort),
			)
			defer res.Body.Close()

			assert.Equal(t, tt.want.status, res.StatusCode)
			if tt.want.wantInternalServerError {
				return
			}
			require.NotEmpty(t, string(url), "response body must be not empty")

			id := strings.TrimPrefix(string(url), "http://localhost:8080")
			res, _ = testRequest(t, ts, http.MethodGet, id, nil)
			defer res.Body.Close()

			require.Contains(t, res.Request.URL.String(), tt.args.urlToShort)
		})
	}

	unsupportedMethods := []string{
		http.MethodConnect,
		http.MethodOptions,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodTrace,
		http.MethodHead,
		http.MethodPut,
	}
	for _, m := range unsupportedMethods {
		t.Run(m, func(t *testing.T) {
			res, _ := testRequest(t, ts, m, "/", nil)
			defer res.Body.Close()
			require.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
		})
	}
}

// only negative cases, because positive cases are in TestServer_HandleURLGetCreate
func TestServer_HandleURLGet(t *testing.T) {
	tests := []struct {
		name   string
		target string
		status int
	}{
		{
			name:   "empty id",
			target: "/",
			status: http.StatusMethodNotAllowed,
		},
		{
			name:   "id doesn't exists",
			target: "/" + uuid.New().String(),
			status: http.StatusNotFound,
		},
	}
	s := New(NewConfig("", "inmemory"))
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, _ := testRequest(t, ts, http.MethodGet, tt.target, nil)
			defer res.Body.Close()
			assert.Equal(t, tt.status, res.StatusCode)
		})
	}
}

func TestServer_HandleURLGetAndCreateJSON(t *testing.T) {
	type request struct {
		URL string `json:"url"`
	}
	type response struct {
		Result string `json:"result"`
	}
	type args struct {
		urlPath string
		request request
	}
	type want struct {
		wantInternalServerError bool
		status                  int
	}
	tests := []struct {
		name string

		args args
		want want
	}{
		{
			name: "positive case #1",
			args: args{
				urlPath: "/api/shorten",
				request: request{
					URL: "https://www.google.com",
				},
			},
			want: want{
				wantInternalServerError: false,
				status:                  http.StatusCreated,
			},
		},
		{
			name: "positive case #2",
			args: args{
				urlPath: "/api/shorten",
				request: request{
					URL: "https://ya.ru",
				},
			},
			want: want{
				wantInternalServerError: false,
				status:                  http.StatusCreated,
			},
		},
		{
			name: "incorrect url to short",
			args: args{
				urlPath: "/api/shorten",
				request: request{
					URL: "hlt v.org",
				},
			},
			want: want{
				wantInternalServerError: true,
				status:                  http.StatusBadRequest,
			},
		},
		{
			name: "empty data",
			args: args{
				urlPath: "/api/shorten",
				request: request{
					URL: "",
				},
			},
			want: want{
				wantInternalServerError: true,
				status:                  http.StatusBadRequest,
			},
		},
	}

	s := New(NewConfig("", "inmemory"))
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp response
			data, err := json.Marshal(tt.args.request)
			require.NoError(t, err)
			body := bytes.NewReader(data)
			res, url := testRequest(
				t,
				ts,
				http.MethodPost,
				tt.args.urlPath,
				body,
			)
			defer res.Body.Close()
			json.Unmarshal(url, &resp)

			assert.Equal(t, tt.want.status, res.StatusCode)
			if tt.want.wantInternalServerError {
				return
			}

			require.NotEmpty(t, resp.Result, "response body must be not empty")

			id := strings.TrimPrefix(resp.Result, "http://localhost:8080")
			res, _ = testRequest(t, ts, http.MethodGet, id, nil)
			defer res.Body.Close()

			require.Contains(t, res.Request.URL.String(), tt.args.request.URL)
		})
	}
}
