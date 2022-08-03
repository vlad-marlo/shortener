package httpserver

import (
	"io"
	"net/http"
	"net/http/httptest"
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
		statusCode int
		headers    http.Header
		body       string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		// TODO: add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.method, tt.args.url, tt.args.body)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(s.handleURLGetCreate)
			h.ServeHTTP(w, request)

			res := w.Result()

			assert.Equal(t, res.StatusCode, tt.want.statusCode)

			if tt.args.headersResponse {
				for k, v := range tt.want.headers {
					assert.Equal(t, res.Header.Get(k), v, "Header[%v] want %v got %v", k, v, res.Header.Get(k))
				}
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
}
