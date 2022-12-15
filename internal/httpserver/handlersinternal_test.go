package httpserver

//goland:noinspection ALL
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/vlad-marlo/shortener/internal/httpserver/middleware"
	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/inmemory"
	mock_store "github.com/vlad-marlo/shortener/internal/store/mock"
	"github.com/vlad-marlo/shortener/internal/store/model"
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
	l, _ := zap.NewProduction()
	s := New(&Config{
		BaseURL:     "http://localhost:8080",
		BindAddr:    "localhost:8080",
		StorageType: store.InMemoryStorage,
	}, storage, l)

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
	l, _ := zap.NewProduction()
	s := New(&Config{
		BaseURL:     "http://localhost:8080",
		BindAddr:    "localhost:8080",
		StorageType: store.InMemoryStorage,
	}, storage, l)

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
	l, _ := zap.NewProduction()
	s := New(&Config{
		BaseURL:     "http://localhost:8080",
		BindAddr:    "localhost:8080",
		StorageType: store.InMemoryStorage,
	}, storage, l)

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

// TestServer_handleURLBulkDelete_Positive ...
func TestServer_handleURLBulkDelete_Positive(t *testing.T) {
	ctrl := gomock.NewController(t)

	storage := mock_store.NewMockStore(ctrl)

	server, td := TestServer(t, storage)
	defer require.NoError(t, td())

	data := `["1", "2", "3"]`

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/shorten/batch", strings.NewReader(data))
	defer assert.NoError(t, w.Result().Body.Close())
	defer assert.NoError(t, r.Body.Close())

	storage.
		EXPECT().
		URLsBulkDelete(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	server.handleURLBulkDelete(w, r)
	assert.Equal(t, http.StatusAccepted, w.Code)
}

// TestServer_handleURLBulkDelete_Negative ...
func TestServer_handleURLBulkDelete_Negative(t *testing.T) {
	ctrl := gomock.NewController(t)
	storage := mock_store.NewMockStore(ctrl)

	server, td := TestServer(t, storage)
	defer require.NoError(t, td())

	data := `["1",`

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", strings.NewReader(data))
	defer assert.NoError(t, w.Result().Body.Close())
	defer assert.NoError(t, r.Body.Close())

	storage.
		EXPECT().
		URLsBulkDelete(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	server.handleURLBulkDelete(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestServer_handlePingStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	storage := mock_store.NewMockStore(ctrl)
	server, td := TestServer(t, storage)
	defer require.NoError(t, td())

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/shorten/batch", nil)
	defer assert.NoError(t, w.Result().Body.Close())
	defer assert.NoError(t, r.Body.Close())

	storage.
		EXPECT().
		Ping(gomock.Any()).
		Return(nil).
		AnyTimes()

	server.handlePingStore(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

// func TestServer_handleURLBulkCreate_Positive(t *testing.T) {
// 	type args struct {
// 		urls []*model.BatchCreateURLsResponse
// 		err  error
// 	}
// 	type want struct {
// 		statusCode int
// 		data       []string
// 	}
// 	tt := []struct {
// 		name string
// 		data string
// 		args args
// 		want want
// 	}{
// 		// {
// 		// 	name: "positive case #1",
// 		// 	data: `[{ "correlation_id": "a", "original_url": "https://ya.ru" }, { "correlation_id": "b", "original_url": "https://yandex.ru" }]`,
// 		// 	args: args{
// 		// 		urls: []*model.BatchCreateURLsResponse{
// 		// 			{
// 		// 				ShortURL:      "a",
// 		// 				CorrelationID: "a",
// 		// 			},
// 		// 			{
// 		// 				ShortURL:      "b",
// 		// 				CorrelationID: "b",
// 		// 			},
// 		// 		},
// 		// 		err: nil,
// 		// 	},
// 		// 	want: want{
// 		// 		statusCode: http.StatusCreated,
// 		// 		data:       []string{"a", "b"},
// 		// 	},
// 		// },
// 		// {
// 		// 	name: "negative case #1",
// 		// 	data: `[]`,
// 		// 	args: args{
// 		// 		urls: nil,
// 		// 		err:  store.ErrAlreadyExists,
// 		// 	},
// 		// 	want: want{
// 		// 		data:       nil,
// 		// 		statusCode: http.StatusBadRequest,
// 		// 	},
// 		// },
// 		// {
// 		// 	name: "negative case #2",
// 		// 	data: `[]`,
// 		// 	args: args{
// 		// 		urls: nil,
// 		// 		err:  store.ErrAlreadyExists,
// 		// 	},
// 		// 	want: want{
// 		// 		data:       nil,
// 		// 		statusCode: http.StatusBadRequest,
// 		// 	},
// 		// },
// 	}
// 	for _, tc := range tt {
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			storage := mock_store.NewMockStore(ctrl)
//
// 			s, td := TestServer(t, storage)
// 			defer require.NoError(t, td())
//
// 			storage.
// 				EXPECT().
// 				URLsBulkCreate(gomock.Any(), gomock.Any()).
// 				Return(tc.args.urls, tc.args.err).
// 				AnyTimes()
//
// 			w := httptest.NewRecorder()
// 			defer assert.NoError(t, w.Result().Body.Close())
//
// 			r := httptest.NewRequest("POST", "/", strings.NewReader(tc.data))
// 			defer assert.NoError(t, r.Body.Close())
//
// 			s.handleURLBulkCreate(w, r)
// 			res := w.Result()
// 			defer assert.NoError(t, res.Body.Close())
// 			assert.Equal(t, tc.want.statusCode, res.StatusCode)
// 			if tc.want.statusCode != http.StatusCreated {
// 				return
// 			}
//
// 			var resp []*model.BatchCreateURLsResponse
// 			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
//
// 			require.Contains(t, "application/json", res.Header.Get("content-type"))
// 			for _, m := range resp {
// 				assert.Contains(t, tc.want.data, m.CorrelationID, "xdddddd", tc.want.data, resp)
// 			}
// 		})
// 	}
// }

func TestServer_handleURLGetAllByUser_Positive(t *testing.T) {
	data := map[string][]*model.URL{
		"1": {
			&model.URL{
				BaseURL: "first",
				User:    "1",
				ID:      "1",
			},
			&model.URL{
				BaseURL: "second",
				User:    "1",
				ID:      "2",
			},
		},
		"2": {
			&model.URL{
				BaseURL: "third",
				User:    "2",
				ID:      "3",
			},
			&model.URL{
				BaseURL: "fourth",
				User:    "2",
				ID:      "4",
			},
		},
	}

	for u, urls := range data {
		t.Run(fmt.Sprintf("test user: %s", u), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			storage := mock_store.NewMockStore(ctrl)

			server, td := TestServer(t, storage)
			defer require.NoError(t, td())

			storage.
				EXPECT().
				GetAllUserURLs(gomock.Any(), u).
				Return(urls, nil).
				AnyTimes()

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/sdf", nil)
			r = r.WithContext(context.WithValue(r.Context(), middleware.UserCTXName, u))
			defer assert.NoError(t, r.Body.Close())
			server.handleGetUserURLs(w, r)

			// t.Logf("status cod: %d %s", w.Result().StatusCode, w.Body.String())

			var responseURLs []*model.AllUserURLsResponse
			for _, u := range urls {
				resp := &model.AllUserURLsResponse{
					ShortURL:    fmt.Sprintf("%s/%s", server.config.BaseURL, u.ID),
					OriginalURL: u.BaseURL,
				}
				responseURLs = append(responseURLs, resp)
			}

			expected, err := json.Marshal(responseURLs)
			require.NoError(t, err, fmt.Sprintf("json marshal: %v", err))
			assert.JSONEq(t, string(expected), w.Body.String())
		})

	}
}
