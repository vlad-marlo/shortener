package httpserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	s = New(NewConfig())
)

func TestServer_handleURLCreate(t *testing.T) {
	type args struct {
		url  string
		body io.Reader
	}
	type want struct {
		statusCode              int
		wantInternalServerError bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "correct post request",
			args: args{
				url:  "/",
				body: strings.NewReader("howdy.ho"),
			},
			want: want{
				statusCode:              http.StatusCreated,
				wantInternalServerError: false,
			},
		},
		{
			name: "uncorrect url",
			args: args{
				url:  "/sdf",
				body: strings.NewReader("howdy.ho"),
			},
			want: want{
				statusCode:              http.StatusInternalServerError,
				wantInternalServerError: true,
			},
		},
		{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.url == "" {
				tt.args.url = "/"
			}
			req := httptest.NewRequest(http.MethodPost, tt.args.url, tt.args.body)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(s.handleURLCreate)
			h.ServeHTTP(w, req)

			res := w.Result()

			if tt.want.wantInternalServerError {
				assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
			} else {
				assert.Equal(t, tt.want.statusCode, res.StatusCode)
				assert.NotNil(t, res.Header.Get("Locate"), "locate header must be not null")
			}

			// if tt.args.bodyResponse {
			// 	resBody, err := io.ReadAll(res.Body)
			// 	defer res.Body.Close()
			//
			// 	assert.NoError(t, err)
			//
			// 	// checking response answer
			// 	assert.Equal(t, string(resBody), tt.want.body)
			// }
		})
	}
}

func TestServer_handlerURLGetCreate_UnsupportedMethods(t *testing.T) {
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
