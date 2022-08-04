package httpserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_handleURLGetCreate(t *testing.T) {
	s := New(NewConfig())
	type args struct {
		method          string
		url             string
		body            io.Reader
		bodyResponse    bool
		headersResponse bool
	}
	type want struct {
		statusCode              int
		body                    string
		wantInternalServerError bool
		contentType             string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "correct post request",
			args: args{
				method:          http.MethodPost,
				url:             "/",
				body:            strings.NewReader("howdy.ho"),
				headersResponse: true,
				bodyResponse:    false,
			},
			want: want{
				statusCode:              http.StatusCreated,
				body:                    "",
				wantInternalServerError: false,
				contentType:             "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.url == "" && tt.args.method == http.MethodPost {
				tt.args.url = "/"
			}
			req := httptest.NewRequest(tt.args.method, tt.args.url, tt.args.body)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(s.handleURLGetCreate)
			h.ServeHTTP(w, req)

			res := w.Result()

			if tt.want.wantInternalServerError && tt.want.statusCode == 0 {
				tt.want.statusCode = http.StatusInternalServerError
			}
			assert.Equal(t, res.StatusCode, tt.want.statusCode)

			if tt.args.method == http.MethodPost && tt.args.headersResponse {
				assert.NotNil(t, res.Header.Get("Locate"), "locate header must be not null")
			}

			if tt.args.bodyResponse {
				defer res.Body.Close()
				resBody, err := io.ReadAll(res.Body)
				assert.NoError(t, err)

				// checking response answer
				assert.Equal(t, string(resBody), tt.want.body)
			}
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
	for _, method := range unsupportedMethods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/", nil)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(s.handleURLGetCreate)
			h.ServeHTTP(w, req)

			res := w.Result()
			assert.Equal(t, res.StatusCode, http.StatusInternalServerError)
		})
	}
}
