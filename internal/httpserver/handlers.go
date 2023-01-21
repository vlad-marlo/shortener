package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

// handleURLGet is redirecting user to base url with id which is provided in url path by chi url params.
func (s *Server) handleURLGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := middleware.GetReqID(ctx)
	fields := []zap.Field{
		zap.String("request_id", reqID),
		zap.String("request_ip", r.Header.Get("X-Real-IP")),
	}

	w.Header().Set("Content-Type", "text/plain")
	id := chi.URLParam(r, "id")

	url, err := s.srv.GetByID(ctx, id)
	switch {
	case errors.Is(err, store.ErrIsDeleted):
		w.WriteHeader(http.StatusGone)
		return
	case errors.Is(err, store.ErrNotFound):
		s.handleErrorOrStatus(w, errors.New("where is no url with that id"), fields, http.StatusNotFound)
		return
	case err != nil:
		s.handleErrorOrStatus(w, fmt.Errorf("internal: %w", err), fields, http.StatusInternalServerError)
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
		zap.String("request_ip", r.Header.Get("X-Real-IP")),
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

	userID := getUserFromRequest(r)

	u, err := s.srv.CreateURL(r.Context(), userID, string(data))
	switch {
	case errors.Is(err, store.ErrAlreadyExists):
		w.WriteHeader(http.StatusConflict)
		_, err = w.Write([]byte(fmt.Sprintf("%s/%s", s.config.BaseURL, u.ID)))

		s.handleErrorOrStatus(w, err, fields, http.StatusInternalServerError)
		return
	case err != nil:
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
		zap.String("request_ip", r.Header.Get("X-Real-IP")),
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

	u := &model.URL{}
	if err = json.Unmarshal(req, u); s.handleErrorOrStatus(w, err, fields, http.StatusInternalServerError) {
		return
	}

	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")
	u, err = s.srv.CreateURL(ctx, userID, u.BaseURL)
	switch {
	case errors.Is(err, store.ErrAlreadyExists):
		w.WriteHeader(http.StatusConflict)
	case s.handleErrorOrStatus(w, err, fields, http.StatusBadRequest):
		return
	default:
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
		zap.String("request_ip", r.Header.Get("X-Real-IP")),
	}
	userID := getUserFromRequest(r)

	ctx := r.Context()

	urls, err := s.srv.GetAllURLsByUser(ctx, userID)
	if s.handleErrorOrStatus(w, err, fields, http.StatusInternalServerError) {
		return
	}
	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	response, err := json.Marshal(urls)
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
		zap.String("request_ip", r.Header.Get("X-Real-IP")),
	}
	ctx := r.Context()

	if err := s.srv.Ping(ctx); err != nil {
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
		zap.String("request_ip", r.Header.Get("X-Real-IP")),
	}
	var (
		data []*model.BulkCreateURLRequest
	)
	body, err := io.ReadAll(r.Body)
	if s.handleErrorOrStatus(w, err, fields, http.StatusInternalServerError) {
		return
	}

	if err = json.Unmarshal(body, &data); s.handleErrorOrStatus(w, err, fields, http.StatusBadRequest) {
		return
	}

	if len(data) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := getUserFromRequest(r)

	var urls []model.URLer
	for _, v := range data {
		urls = append(urls, v)
	}

	resp, err := s.srv.CreateManyURLs(r.Context(), userID, urls)
	if s.handleErrorOrStatus(w, err, fields, http.StatusBadRequest) {
		return
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
		zap.String("request_ip", r.Header.Get("X-Real-IP")),
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
	s.srv.DeleteManyURLs(userID, data)
	w.WriteHeader(http.StatusAccepted)
}

// handleInternalStats give trusted user access to specific stats about data records.
func (s *Server) handleInternalStats(w http.ResponseWriter, r *http.Request) {
	fields := []zap.Field{
		zap.String("request_id", middleware.GetReqID(r.Context())),
		zap.String("request_ip", r.Header.Get("X-Real-IP")),
	}

	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP == "" {
		s.handleErrorOrStatus(w, errors.New("IP was not provided in header X-Real-IP"), fields, http.StatusForbidden)
		return
	}

	stat, err := s.srv.GetInternalStats(r.Context(), xRealIP)
	if err != nil {
		s.handleErrorOrStatus(w, err, fields, http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(stat)
	if err != nil {
		s.handleErrorOrStatus(w, err, fields, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(data); err != nil {
		s.handleErrorOrStatus(w, err, fields, http.StatusInternalServerError)
	}
}
