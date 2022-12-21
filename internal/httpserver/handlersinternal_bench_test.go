package httpserver

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vlad-marlo/shortener/internal/store"
	mock_store "github.com/vlad-marlo/shortener/internal/store/mock"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

func BenchmarkServer_handleURLGet(b *testing.B) {
	tt := map[string]error{
		"no error":              nil,
		"with is-deleted error": store.ErrIsDeleted,
		"with not found":        store.ErrNotFound,
		"unknown error":         errors.New(""),
	}
	for name, err := range tt {
		b.Run(name, func(b *testing.B) {
			ctrl := gomock.NewController(b)
			storage := mock_store.NewMockStore(ctrl)

			s, td := TestServer(b, storage)
			defer require.NoError(b, td())

			u := &model.URL{
				BaseURL:   "xdsd",
				ID:        "xd",
				IsDeleted: false,
			}

			storage.
				EXPECT().
				GetByID(
					gomock.Any(),
					gomock.Any(),
				).
				Return(u, err).
				AnyTimes()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodGet, "/xd", nil)
				b.StartTimer()

				s.handleURLGet(w, r)

				b.StopTimer()
				res := w.Result()
				assert.NoError(b, r.Body.Close(), fmt.Sprintf("iteration: %d", i))
				assert.NoError(b, res.Body.Close(), fmt.Sprintf("iteration: %d", i))
				b.StartTimer()

			}
		})
	}
}

func BenchmarkServer_handleURLPost(b *testing.B) {
	tt := map[string]error{
		"no error":     nil,
		"unknown err":  errors.New(""),
		"exists error": store.ErrAlreadyExists,
	}

	for name, err := range tt {
		b.Run(name, func(b *testing.B) {

			ctrl := gomock.NewController(b)
			storage := mock_store.NewMockStore(ctrl)

			s, td := TestServer(b, storage)
			defer require.NoError(b, td())
			storage.
				EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Return(err).
				AnyTimes()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://ya.ru/"))
				w := httptest.NewRecorder()
				b.StartTimer()

				s.handleURLCreate(w, r)

				b.StopTimer()
				res := w.Result()
				assert.NoError(b, res.Body.Close())
				assert.NoError(b, r.Body.Close())
				b.StartTimer()
			}

		})
	}
}

func BenchmarkServer_handleURLPostJSON(b *testing.B) {
	tt := map[string]error{
		"no error":     nil,
		"unknown err":  errors.New(""),
		"exists error": store.ErrAlreadyExists,
	}

	for name, err := range tt {
		b.Run(name, func(b *testing.B) {
			// prepare mock storage
			ctrl := gomock.NewController(b)
			storage := mock_store.NewMockStore(ctrl)
			storage.
				EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Return(err).
				AnyTimes()
			s, td := TestServer(b, storage)
			defer require.NoError(b, td())

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				// prepare test recorder and test request
				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url": "ya.ru"}`))
				b.StartTimer()

				s.handleURLCreateJSON(w, r)

				b.StopTimer()
				res := w.Result()
				assert.NoError(b, res.Body.Close())
				assert.NoError(b, r.Body.Close())
				b.StartTimer()
			}
		})
	}
}

func BenchmarkServer_handleURLBatchCreate(b *testing.B) {
	data := `
	[
		{"original_url": "ya.ru/a", "correlation_id": "a"},
		{"original_url": "ya.ru/b", "correlation_id": "b"},
		{"original_url": "ya.ru/c", "correlation_id": "c"}
	]`
	tt := map[string]error{
		"no error": nil,
		"error":    store.ErrAlreadyExists,
	}
	for name, err := range tt {
		b.Run(name, func(b *testing.B) {
			ctrl := gomock.NewController(b)
			storage := mock_store.NewMockStore(ctrl)
			storage.
				EXPECT().
				URLsBulkCreate(gomock.Any(), gomock.Any()).
				Return(nil, err).
				AnyTimes()
			s, td := TestServer(b, storage)
			defer require.NoError(b, td())

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()

				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/", strings.NewReader(data))
				b.StartTimer()

				s.handleURLBulkCreate(w, r)

				b.StopTimer()
				res := w.Result()
				assert.NoError(b, res.Body.Close())
				assert.NoError(b, r.Body.Close())
				b.StartTimer()
				b.ReportMetric(2, "B/op")
			}
		})
	}

}
