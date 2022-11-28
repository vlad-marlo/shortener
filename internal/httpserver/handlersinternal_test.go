package httpserver

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/vlad-marlo/logger"

	"github.com/vlad-marlo/shortener/internal/store/inmemory"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vlad-marlo/shortener/internal/store"
)

// testRequest ...
func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, []byte) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, respBody
}

// TestServer_HandleURLGetAndCreate ...
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
				urlToShort: "https://google.com",
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
				urlToShort: "https://ya.ru",
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
				urlToShort: "https://yandex.ru",
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
				urlToShort: "https://hl tv.org",
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

	storage := inmemory.New()
	s := New(&Config{
		BaseURL:     "http://localhost:8080",
		BindAddr:    "localhost:8080",
		StorageType: store.InMemoryStorage,
	}, storage, logrus.New())

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
			defer require.NoError(t, res.Body.Close())

			assert.Equal(t, tt.want.status, res.StatusCode)
			if tt.want.wantInternalServerError {
				return
			}
			require.NotEmpty(t, string(url), "response body must be not empty")

			id := strings.TrimPrefix(string(url), "http://localhost:8080")
			res, _ = testRequest(t, ts, http.MethodGet, id, nil)
			defer require.NoError(t, res.Body.Close())

			require.Contains(t, res.Request.URL.String(), strings.TrimPrefix("https://", tt.args.urlToShort))
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
			defer require.NoError(t, res.Body.Close())
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

	storage := inmemory.New()
	s := New(&Config{
		BaseURL:     "http://localhost:8080",
		BindAddr:    "localhost:8080",
		StorageType: store.InMemoryStorage,
	}, storage, logrus.New())

	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, _ := testRequest(t, ts, http.MethodGet, tt.target, nil)
			defer require.NoError(t, res.Body.Close())
			assert.Equal(t, tt.status, res.StatusCode)
		})
	}
}

// TestServer_HandleURLGetAndCreateJSON ...
func TestServer_HandleURLGetAndCreateJSON(t *testing.T) {
	type (
		request struct {
			URL string `json:"url"`
		}
		response struct {
			Result string `json:"result"`
		}
		args struct {
			urlPath string
			request request
		}
		want struct {
			wantInternalServerError bool
			status                  int
		}
	)
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
					URL: "https://hlt v.org",
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

	storage := inmemory.New()
	s := New(&Config{
		BaseURL:     "http://localhost:8080",
		BindAddr:    "localhost:8080",
		StorageType: store.InMemoryStorage,
	}, storage, logrus.New())

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
			defer require.NoError(t, res.Body.Close())
			_ = json.Unmarshal(url, &resp)

			assert.Equal(t, tt.want.status, res.StatusCode)
			if tt.want.wantInternalServerError {
				return
			}

			require.NotEmpty(t, resp.Result, "response body must be not empty")

			id := strings.TrimPrefix(resp.Result, "http://localhost:8080")
			res, _ = testRequest(t, ts, http.MethodGet, id, nil)
			defer require.NoError(t, res.Body.Close())

			require.Contains(t, res.Request.URL.String(), tt.args.request.URL)
		})
	}
}

func BenchmarkServer_handleURLGet(b *testing.B) {
	storage := inmemory.New()
	s := New(
		&Config{
			BaseURL:     "http://localhost:8080",
			BindAddr:    "localhost:8080",
			StorageType: store.InMemoryStorage,
		},
		storage,
		logger.WithOpts(
			logger.WithOutput(io.Discard),
		),
	)

	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		req, err := http.NewRequest(http.MethodGet, ts.URL+"/xd", nil)
		if err != nil {
			b.Fatalf("new request: %v", err)
		}

		b.StartTimer()

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			b.Fatalf("doing http request: %v", err)
		}

		if err := resp.Body.Close(); err != nil {
			b.Fatalf("close body: %v", err)
		}
	}
}

func BenchmarkServer_handleURLPost(b *testing.B) {
	storage := inmemory.New()
	s := New(
		&Config{
			BaseURL:     "http://localhost:8080",
			BindAddr:    "localhost:8080",
			StorageType: store.InMemoryStorage,
		},
		storage,
		logger.WithOpts(
			logger.WithOutput(io.Discard),
		),
	)

	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		req, err := http.NewRequest(http.MethodPost, ts.URL+"/", strings.NewReader("https://ya.ru/"))
		if err != nil {
			b.Fatalf("new request: %v", err)
		}

		b.StartTimer()

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			b.Fatalf("doing http request: %v", err)
		}

		if err := resp.Body.Close(); err != nil {
			b.Fatalf("close body: %v", err)
		}
	}
}

func BenchmarkServer_handleURLPostJSON(b *testing.B) {
	storage := inmemory.New()
	s := New(
		&Config{
			BaseURL:     "http://localhost:8080",
			BindAddr:    "localhost:8080",
			StorageType: store.InMemoryStorage,
		},
		storage,
		logger.WithOpts(
			logger.WithOutput(io.Discard),
		),
	)

	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/shorten", strings.NewReader(`{"url": "ya.ru"}`))
		if err != nil {
			b.Fatalf("new request: %v", err)
		}

		b.StartTimer()

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			b.Fatalf("doing http request: %v", err)
		}

		if err := resp.Body.Close(); err != nil {
			b.Fatalf("close body: %v", err)
		}
	}
}

func BenchmarkServer_handleURLBatchCreate(b *testing.B) {
	storage := inmemory.New()
	s := New(
		&Config{
			BaseURL:     "http://localhost:8080",
			BindAddr:    "localhost:8080",
			StorageType: store.InMemoryStorage,
		},
		storage,
		logger.WithOpts(
			logger.WithOutput(io.Discard),
		),
	)

	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		data := `
		[
			{"original_url": "ya.ru/a", "correlation_id": "a"},
			{"original_url": "ya.ru/b", "correlation_id": "b"},
			{"original_url": "ya.ru/c", "correlation_id": "c"}
		]
		`
		req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/shorten/batch", strings.NewReader(data))
		if err != nil {
			b.Fatalf("new request: %v", err)
		}

		b.StartTimer()

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			b.Fatalf("doing http request: %v", err)
		}

		if err := resp.Body.Close(); err != nil {
			b.Fatalf("close body: %v", err)
		}
	}
}
