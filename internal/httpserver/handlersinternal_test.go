package httpserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestServer_HandleURLGetAndCreate(t *testing.T) {
	type args struct {
		urlPath    string
		urlToShort string
	}
	type want struct {
		wantInternalServerError bool
	}
	tests := []struct {
		name string

		args args
		want want
	}{
		{
			name: "positive case #1",
			args: args{
				urlPath:    "/",
				urlToShort: "https://google.com",
			},
			want: want{
				wantInternalServerError: false,
			},
		},
		{
			name: "positive case #2",
			args: args{
				urlPath:    "/",
				urlToShort: "https://ya.ru",
			},
			want: want{
				wantInternalServerError: false,
			},
		},
		{
			name: "url already exists",
			args: args{
				urlPath:    "/",
				urlToShort: "https://ya.ru",
			},
			want: want{
				wantInternalServerError: true,
			},
		},
		{
			name: "uncorrect target case",
			args: args{
				urlPath:    "/jkljk/",
				urlToShort: "https://yandex.ru",
			},
			want: want{
				wantInternalServerError: true,
			},
		},
		{
			name: "uncorrect url to short",
			args: args{
				urlPath:    "/",
				urlToShort: "https://hlt v.org",
			},
			want: want{
				wantInternalServerError: true,
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
			},
		},
	}

	s := NewTestServer(NewConfig("", "inmemory"))
	ts := httptest.NewServer(s.Router)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, url := testRequest(t, ts, http.MethodPost, tt.args.urlPath, strings.NewReader(tt.args.urlToShort))

			if tt.want.wantInternalServerError {
				assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
				return
			}
			assert.Equal(t, http.StatusCreated, res.StatusCode)

			require.NotEmpty(t, url, "response body must be not empty")

			res, _ = testRequest(t, ts, http.MethodGet, string(url), nil)

			require.Equal(t, tt.args.urlToShort, res.Header.Get("location"))
			require.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
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
			require.Equal(t, http.StatusInternalServerError, res.StatusCode)
		})
	}
}

// only negative cases, because positive cases are in TestServer_HandleURLGetCreate
func TestServer_HandleURLGet(t *testing.T) {
	tests := []struct {
		name   string
		target string
	}{
		{
			name:   "empty id",
			target: "/",
		},
		{
			name:   "id doesn't exists",
			target: "/" + uuid.New().String(),
		},
	}
	s := NewTestServer(NewConfig("", "inmemory"))
	ts := httptest.NewServer(s.Router)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, _ := testRequest(t, ts, http.MethodGet, tt.target, nil)
			assert.Equal(t, http.StatusNotFound, res.StatusCode)
		})
	}
}
