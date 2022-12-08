package httpserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/vlad-marlo/logger"

	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/inmemory"
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
