package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"

	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

// handleURLGet is redirecting user to base url with id which is provided in url path by chi url params.
func (s *Server) handleURLGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := middleware.GetReqID(ctx)
	fields := []zap.Field{
		zap.String("request_id", reqID),
		zap.String("handler", "get url by id"),
	}

	w.Header().Set("Content-Type", "text/plain")
	id := chi.URLParam(r, "id")

	url, err := s.store.GetByID(ctx, id)
	if errors.Is(err, store.ErrIsDeleted) {
		w.WriteHeader(http.StatusGone)
		return
	} else if err != nil {
		s.handleErrorOrStatus(w, errors.New("where is no url with that id"), fields, http.StatusNotFound)
		return
	}

	w.Header().Set("Location", url.BaseURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// handleURLCreate is http handler which creates record about url and return
// short link to url in response.
//
// If url was already registered when handler will return old value.
func (s *Server) handleURLCreate(w http.ResponseWriter, r *http.Request) {
	// setting up response meta info
	fields := []zap.Field{
		zap.String("request_id", middleware.GetReqID(r.Context())),
	}
	w.Header().Set("Content-Type", "text/plain")

	data, err := io.ReadAll(r.Body)
	if s.handleErrorOrStatus(w, err, fields, http.StatusInternalServerError) {
		return
	}

	defer func() {
		if err = r.Body.Close(); err != nil {
			s.logger.Error(fmt.Sprintf("request body close: %v", err), fields...)
		}
	}()

	if len(data) == 0 {
		s.handleErrorOrStatus(w, ErrIncorrectRequestBody, fields, http.StatusBadRequest)
		return
	}

	if s.handleErrorOrStatus(w, err, fields, http.StatusBadRequest) {
		return
	}

	userID := getUserFromRequest(r)

	u, err := model.NewURL(string(data), userID)
	if s.handleErrorOrStatus(w, err, fields, http.StatusBadRequest) {
		return
	}

	ctx := r.Context()

	if err = s.store.Create(ctx, u); err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			_, err = w.Write([]byte(fmt.Sprintf("%s/%s", s.config.BaseURL, u.ID)))

			s.handleErrorOrStatus(w, err, fields, http.StatusInternalServerError)
			return
		}

		s.handleErrorOrStatus(w, err, fields, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(fmt.Sprintf("%s/%s", s.config.BaseURL, u.ID)))
	s.handleErrorOrStatus(w, err, fields, http.StatusInternalServerError)
}

// handleURLCreateJSON is http handler which creates record about url and return
// short link to url in response.
//
// If url was already registered when handler will return old value.
func (s *Server) handleURLCreateJSON(w http.ResponseWriter, r *http.Request) {
	fields := []zap.Field{
		zap.String("request_id", middleware.GetReqID(r.Context())),
	}

	req, err := io.ReadAll(r.Body)
	defer func() {
		if err = r.Body.Close(); err != nil {
			s.logger.Warn(fmt.Sprintf("request body close: %v", err), fields...)
		}
	}()
	if s.handleErrorOrStatus(w, err, fields, http.StatusInternalServerError) {
		return
	}

	userID := getUserFromRequest(r)

	u := &model.URL{
		User: userID,
	}
	if err = json.Unmarshal(req, u); s.handleErrorOrStatus(w, err, fields, http.StatusInternalServerError) {
		return
	}

	if err = u.ShortURL(); s.handleErrorOrStatus(w, err, fields, http.StatusBadRequest) {
		return
	}

	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")
	if err = s.store.Create(ctx, u); errors.Is(err, store.ErrAlreadyExists) {
		w.WriteHeader(http.StatusConflict)
	} else if s.handleErrorOrStatus(w, err, fields, http.StatusBadRequest) {
		return
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	resp := model.ResultResponse{
		Result: fmt.Sprintf("%s/%s", s.config.BaseURL, u.ID),
	}
	res, err := json.Marshal(resp)
	if s.handleErrorOrStatus(w, err, fields, http.StatusInternalServerError) {
		return
	}

	_, err = w.Write(res)
	s.handleErrorOrStatus(w, err, fields, http.StatusInternalServerError)
}

// handleGetUserURLs is http handler which return to user all records which was created by him.
func (s *Server) handleGetUserURLs(w http.ResponseWriter, r *http.Request) {
	fields := []zap.Field{
		zap.String("request_id", middleware.GetReqID(r.Context())),
	}
	userID := getUserFromRequest(r)

	ctx := r.Context()

	urls, err := s.store.GetAllUserURLs(ctx, userID)
	if s.handleErrorOrStatus(w, err, fields, http.StatusInternalServerError) {
		return
	}
	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var responseURLs []*model.AllUserURLsResponse
	for _, u := range urls {
		resp := &model.AllUserURLsResponse{
			ShortURL:    fmt.Sprintf("%s/%s", s.config.BaseURL, u.ID),
			OriginalURL: u.BaseURL,
		}
		responseURLs = append(responseURLs, resp)
	}

	response, err := json.Marshal(responseURLs)
	if s.handleErrorOrStatus(w, err, fields, http.StatusInternalServerError) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	s.handleErrorOrStatus(w, err, fields, http.StatusInternalServerError)
}

// handlePingStore is debug handler which gives user access to check db connection.
//
// In order that storage is not available, handler will return http status 500. In other cases 200.
func (s *Server) handlePingStore(w http.ResponseWriter, r *http.Request) {
	fields := []zap.Field{
		zap.String("request_id", middleware.GetReqID(r.Context())),
	}
	ctx := r.Context()

	if err := s.store.Ping(ctx); err != nil {
		s.handleErrorOrStatus(w, fmt.Errorf("handlePingStore: %w", err), fields, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// handleURLBulkCreate is http handler to create many records about many urls in one time.
//
// User gives him json object like
// [{ "correlation_id": "1", "original_url": "https://ya.ru" }]
// after creation in success case response will be like
// [{ "correlation_id": "1", "short_url": "http://<server_addr>/<id>"}].
func (s *Server) handleURLBulkCreate(w http.ResponseWriter, r *http.Request) {
	fields := []zap.Field{
		zap.String("request_id", middleware.GetReqID(r.Context())),
	}
	var (
		data []*model.BulkCreateURLRequest
		urls []*model.URL
	)
	body, err := io.ReadAll(r.Body)
	if s.handleErrorOrStatus(w, err, fields, http.StatusInternalServerError) {
		return
	}

	if err = json.Unmarshal(body, &data); s.handleErrorOrStatus(w, err, fields, http.StatusInternalServerError) {
		return
	}

	if len(data) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := getUserFromRequest(r)

	for _, v := range data {
		var u *model.URL
		u, err = model.NewURL(v.OriginalURL, userID, v.CorrelationID)
		if s.handleErrorOrStatus(w, err, fields, http.StatusInternalServerError) {
			return
		}
		urls = append(
			urls,
			u,
		)
	}

	ctx := r.Context()

	resp, err := s.store.URLsBulkCreate(ctx, urls)
	if s.handleErrorOrStatus(w, err, fields, http.StatusBadRequest) {
		return
	}
	for _, v := range resp {
		id := v.ShortURL
		v.ShortURL = fmt.Sprintf("%s/%s", s.config.BaseURL, id)
	}

	body, err = json.Marshal(resp)
	if s.handleErrorOrStatus(w, err, fields, http.StatusInternalServerError) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err = w.Write(body); err != nil {
		s.logger.Error(fmt.Sprintf("write response: %v", err), fields...)
	}
}

// handleURLBulkDelete is http handler which gives user access to delete all urls
// which was created by him.
//
// Request must be json array of strings where every element is url id.
// Only user which create url have access to deleting urls.
func (s *Server) handleURLBulkDelete(w http.ResponseWriter, r *http.Request) {
	fields := []zap.Field{
		zap.String("request_id", middleware.GetReqID(r.Context())),
	}
	var data []string
	userID := getUserFromRequest(r)

	defer func() {
		if err := r.Body.Close(); err != nil {
			s.logger.Error(fmt.Sprintf("defering request body close: %v ", err), fields...)
		}
	}()

	body, err := io.ReadAll(r.Body)
	if s.handleErrorOrStatus(w, err, fields, http.StatusInternalServerError) {
		return
	}

	if err := json.Unmarshal(body, &data); err != nil {
		s.handleErrorOrStatus(w, fmt.Errorf("handle bulk url delete: json unmarshal data: %w", err), fields, http.StatusBadRequest)
		return
	}
	s.poller.DeleteURLs(data, userID)
	w.WriteHeader(http.StatusAccepted)
}
