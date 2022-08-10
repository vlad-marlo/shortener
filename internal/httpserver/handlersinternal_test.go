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

var (
	s = New(NewConfig())
)

func TestServer_HandleURLGetCreate(t *testing.T) {
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
	s := NewTestServer(NewConfig())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.args.urlPath, strings.NewReader(tt.args.urlToShort))
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(s.handleURLGetCreate)
			handler.ServeHTTP(w, req)
			res := w.Result()

			if tt.want.wantInternalServerError {
				assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
				return
			}
			assert.Equal(t, http.StatusCreated, res.StatusCode)

			defer res.Body.Close()
			url, err := io.ReadAll(res.Body)

			require.NoError(t, err, "error in response body")
			require.NotEmpty(t, url, "response body must be not empty")

			req = httptest.NewRequest(
				http.MethodGet,
				string(url),
				nil,
			)
			w = httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			res = w.Result()

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
			req := httptest.NewRequest(m, "/", nil)
			w := httptest.NewRecorder()

			h := http.HandlerFunc(s.handleURLGetCreate)
			h.ServeHTTP(w, req)

			res := w.Result()
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
	s := NewTestServer(NewConfig())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(s.handleURLGet)
			req := httptest.NewRequest(http.MethodGet, tt.target, nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
		})
	}
}
