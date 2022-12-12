package httpserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/vlad-marlo/logger"

	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/inmemory"
	mock_store "github.com/vlad-marlo/shortener/internal/store/mock"
)

func BenchmarkServer_handleURLGet(b *testing.B) {
	storage := inmemory.New()
	s := New(
		&Config{
			BaseURL:     "http://localhost:8080",
			BindAddr:    "localhost:8080",
			StorageType: store.InMemoryStorage,
		},
		storage,
		logrus.NewEntry(logger.WithOpts(
			logger.WithOutput(io.Discard),
		)),
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
		logrus.NewEntry(logger.WithOpts(
			logger.WithOutput(io.Discard),
		)),
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
		logrus.NewEntry(logger.WithOpts(
			logger.WithOutput(io.Discard),
		)),
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
	ctrl := gomock.NewController(b)
	storage := mock_store.NewMockStore(ctrl)
	storage.
		EXPECT().
		URLsBulkCreate(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		AnyTimes()
	s, td := TestServer(b, storage)
	defer func() {
		if err := td(); err != nil {
			b.Fatalf("close server: %v", err)
		}
	}()

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
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/", strings.NewReader(data))
		s.handleURLBulkCreate(w, r)
		if err := r.Body.Close(); err != nil {
			b.Fatalf("close body: %v", err)
		}
	}
}
