package httpserver

//goland:noinspection ALL
import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

var errUnknownErr = errors.New("")

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
			URL string `json:"result"`
		}
		args struct {
			request      request
			prepareStore func(s *mock_store.MockStore)
			id           string
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
				request: request{
					URL: "https://www.google.com",
				},
				prepareStore: func(s *mock_store.MockStore) {
					s.EXPECT().
						Create(gomock.Any(), gomock.Any()).
						DoAndReturn(func(arg0 context.Context, arg1 *model.URL) error {
							arg1.ID = "a"
							return nil
						}).
						AnyTimes()
				},
				id: "a",
			},

			want: want{
				wantInternalServerError: false,
				status:                  http.StatusCreated,
			},
		},
		// {
		// 	name: "positive case #2",
		// 	args: args{
		// 		urlPath: "/api/shorten",
		// 		request: request{
		// 			URL: "https://ya.ru",
		// 		},
		// 	},
		// 	want: want{
		// 		wantInternalServerError: false,
		// 		status:                  http.StatusCreated,
		// 	},
		// },
		// {
		// 	name: "incorrect url to short",
		// 	args: args{
		// 		urlPath: "/api/shorten",
		// 		request: request{
		// 			URL: "https://hlt v.org",
		// 		},
		// 	},
		// 	want: want{
		// 		wantInternalServerError: true,
		// 		status:                  http.StatusBadRequest,
		// 	},
		// },
		// {
		// 	name: "empty data",
		// 	args: args{
		// 		urlPath: "/api/shorten",
		// 		request: request{
		// 			URL: "",
		// 		},
		// 	},
		// 	want: want{
		// 		wantInternalServerError: true,
		// 		status:                  http.StatusBadRequest,
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			storage := mock_store.NewMockStore(ctrl)
			tt.args.prepareStore(storage)

			s, td := TestServer(t, storage)
			defer require.NoError(t, td())

			data, err := json.Marshal(tt.args.request)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))

			s.handleURLCreateJSON(w, r)
			res := w.Result()
			defer require.NoError(t, res.Body.Close())
			defer require.NoError(t, r.Body.Close())

			assert.Equal(t, tt.want.status, res.StatusCode)
			if tt.want.wantInternalServerError {
				return
			}

			require.NotEmpty(t, w.Body.String(), "response body must be not empty")
			// t.Log(w.Body.String())
			var result response
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			// t.Logf("%+v", result)

			id := strings.TrimPrefix(result.URL, "http://localhost:8080/")
			require.Equal(t, tt.args.id, id)
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
	tt := []struct {
		name string
		err  error
		code int
	}{
		{
			name: "positive case #1",
			err:  nil,
			code: http.StatusOK,
		},
		{
			name: "negative case",
			err:  store.ErrNotAccessible,
			code: http.StatusInternalServerError,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			storage := mock_store.NewMockStore(ctrl)
			server, td := TestServer(t, storage)
			defer require.NoError(t, td())

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/", nil)
			defer assert.NoError(t, w.Result().Body.Close())
			defer assert.NoError(t, r.Body.Close())

			storage.
				EXPECT().
				Ping(gomock.Any()).
				Return(tc.err).
				AnyTimes()

			server.handlePingStore(w, r)
			assert.Equal(t, tc.code, w.Code)
		})
	}

}

func TestServer_handleURLBulkCreate_mock(t *testing.T) {
	type mockData struct {
		err  error
		urls []*model.BatchCreateURLsResponse
	}
	type want struct {
		code int
	}
	tt := []struct {
		name string
		data string
		mock mockData
		want want
	}{
		{
			name: "empty data",
			data: "",
			mock: mockData{
				err:  nil,
				urls: []*model.BatchCreateURLsResponse{},
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "positive case #1",
			data: `[{"correlation_id": "a", "original_url": "https://xd.com"}]`,
			mock: mockData{
				err: nil,
				urls: []*model.BatchCreateURLsResponse{
					{
						ShortURL:      "abcd",
						CorrelationID: "a",
					},
				},
			},
			want: want{
				code: http.StatusCreated,
			},
		},
		{
			name: "positive case #2",
			data: `[
	{"correlation_id": "a", "original_url": "https://xd.com"},
	{"correlation_id": "b", "original_url": "https://ya.rmarka"}
]`,
			mock: mockData{
				err: nil,
				urls: []*model.BatchCreateURLsResponse{
					{
						ShortURL:      "xd",
						CorrelationID: "a",
					},
					{
						ShortURL:      "ya",
						CorrelationID: "b",
					},
				},
			},
			want: want{
				code: http.StatusCreated,
			},
		},
		{
			name: "negative case #1",
			data: `[
	{"correlation_id": "a", "original_url": "https://xd.com"},
	{"correlation_id": "b", "original_url": "https://ya.rmarka"}
]`,
			mock: mockData{
				err:  store.ErrAlreadyExists,
				urls: nil,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			storage := mock_store.NewMockStore(ctrl)
			storage.
				EXPECT().
				URLsBulkCreate(gomock.Any(), gomock.Any()).
				Return(tc.mock.urls, tc.mock.err).
				AnyTimes()
			s, td := TestServer(t, storage)
			defer require.NoError(t, td())

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/", strings.NewReader(tc.data))

			s.handleURLBulkCreate(w, r)

			res := w.Result()
			defer assert.NoError(t, res.Body.Close())
			defer assert.NoError(t, r.Body.Close())

			assert.Equal(t, tc.want.code, res.StatusCode)
			if tc.want.code != http.StatusCreated {
				return
			}

			data, err := json.Marshal(tc.mock.urls)
			require.NoError(t, err)
			assert.JSONEq(t, string(data), w.Body.String())
		})
	}

}

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

func TestServer_handleURLGet(t *testing.T) {
	type args struct {
		id  string
		url string
	}
	type mock struct {
		u     *model.URL
		error error
	}
	tt := []struct {
		name string
		args args
		mock mock
		code int
	}{
		{
			name: "positive case",
			args: args{
				id:  "a",
				url: "https://ya.ru",
			},
			mock: mock{
				u: &model.URL{
					BaseURL:   "https://ya.ru",
					ID:        "a",
					IsDeleted: false,
				},
				error: nil,
			},
			code: http.StatusTemporaryRedirect,
		},
		{
			name: "is deleted",
			args: args{
				id:  "a",
				url: "https://ya.ru",
			},
			mock: mock{
				u:     nil,
				error: store.ErrIsDeleted,
			},
			code: http.StatusGone,
		},
		{
			name: "internal error",
			args: args{
				id:  "a",
				url: "https://ya.ru",
			},
			mock: mock{
				u:     nil,
				error: errUnknownErr,
			},
			code: http.StatusNotFound,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			storage := mock_store.NewMockStore(ctrl)
			s, td := TestServer(t, storage)
			defer require.NoError(t, td())

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/a", strings.NewReader(""))

			storage.
				EXPECT().
				GetByID(gomock.Any(), gomock.Any()).
				Return(tc.mock.u, tc.mock.error).
				AnyTimes()

			s.handleURLGet(w, r)

			res := w.Result()
			defer require.NoError(t, res.Body.Close())
			defer require.NoError(t, r.Body.Close())

			assert.Equal(t, tc.code, res.StatusCode)
			if tc.code != http.StatusTemporaryRedirect {
				return
			}
			assert.Contains(t, tc.args.url, res.Header.Get("location"))
		})
	}
}
