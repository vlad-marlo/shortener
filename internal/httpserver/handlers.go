package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/vlad-marlo/shortener/internal/httpserver/middleware"
	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

const (
	cancelCoolDown time.Duration = 30 * time.Millisecond
)

// handleURLGet ...
func (s *Server) handleURLGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	id := chi.URLParam(r, "id")

	ctx, cancel := context.WithTimeout(r.Context(), cancelCoolDown)
	defer cancel()

	url, err := s.Store.GetByID(ctx, id)
	if errors.Is(err, store.ErrIsDeleted) {
		w.WriteHeader(http.StatusGone)
		return
	} else if err != nil {
		s.handleErrorOrStatus(w, errors.New("where is no url with that id"), http.StatusNotFound)
		return
	}

	w.Header().Set("Location", url.BaseURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// handleURLCreate ...
func (s *Server) handleURLCreate(w http.ResponseWriter, r *http.Request) {
	// setting up response meta info
	w.Header().Set("Content-Type", "text/plain")

	data, err := io.ReadAll(r.Body)
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("r body close: %v", err)
		}
	}()

	if s.handleErrorOrStatus(w, err, http.StatusInternalServerError) {
		return
	}
	if len(data) == 0 {
		s.handleErrorOrStatus(w, ErrIncorrectRequestBody, http.StatusBadRequest)
		return
	}

	if s.handleErrorOrStatus(w, err, http.StatusBadRequest) {
		return
	}

	user := r.Context().Value(middleware.UserCTXName)
	var userID string
	if user == nil {
		userID = middleware.UserIDDefaultValue
	} else {
		userID = user.(string)
	}

	u, err := model.NewURL(string(data), userID)
	if s.handleErrorOrStatus(w, err, http.StatusBadRequest) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), cancelCoolDown)
	defer cancel()

	if err = s.Store.Create(ctx, u); err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			_, err = w.Write([]byte(fmt.Sprintf("%s/%s", s.Config.BaseURL, u.ID)))
			s.handleErrorOrStatus(w, err, http.StatusInternalServerError)
			return
		}
		s.handleErrorOrStatus(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	// generate full url like <base service url>/<url identificator>
	_, err = w.Write([]byte(fmt.Sprintf("%s/%s", s.Config.BaseURL, u.ID)))
	s.handleErrorOrStatus(w, err, http.StatusInternalServerError)
}

// handleURLCreateJSON ...
func (s *Server) handleURLCreateJSON(w http.ResponseWriter, r *http.Request) {
	req, err := io.ReadAll(r.Body)
	defer func() {
		if err = r.Body.Close(); err != nil {
			log.Printf("r body close: %v", err)
		}
	}()
	if s.handleErrorOrStatus(w, err, http.StatusInternalServerError) {
		return
	}

	var userID string
	if user := r.Context().Value(middleware.UserCTXName); user != nil {
		userID = user.(string)
	} else {
		userID = middleware.UserIDDefaultValue
	}

	u := &model.URL{
		User: userID,
	}
	if err = json.Unmarshal(req, u); s.handleErrorOrStatus(w, err, http.StatusInternalServerError) {
		return
	}

	if err = u.ShortURL(); s.handleErrorOrStatus(w, err, http.StatusBadRequest) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), cancelCoolDown)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")
	if err = s.Store.Create(ctx, u); errors.Is(err, store.ErrAlreadyExists) {
		w.WriteHeader(http.StatusConflict)
	} else if s.handleErrorOrStatus(w, err, http.StatusBadRequest) {
		return
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	resp := model.ResultResponse{
		Result: fmt.Sprintf("%s/%s", s.Config.BaseURL, u.ID),
	}
	res, err := json.Marshal(resp)
	if s.handleErrorOrStatus(w, err, http.StatusInternalServerError) {
		return
	}

	_, err = w.Write(res)
	s.handleErrorOrStatus(w, err, http.StatusInternalServerError)
}

// handleGetUserURLs ...
func (s *Server) handleGetUserURLs(w http.ResponseWriter, r *http.Request) {
	var userID string
	if user := r.Context().Value(middleware.UserCTXName); user != nil {
		userID = user.(string)
	} else {
		userID = middleware.UserIDDefaultValue
	}

	ctx, cancel := context.WithTimeout(r.Context(), cancelCoolDown)
	defer cancel()

	urls, err := s.Store.GetAllUserURLs(ctx, userID)
	if s.handleErrorOrStatus(w, err, http.StatusInternalServerError) {
		return
	}
	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	responseURLs := []*model.AllUserURLsResponse{}
	for _, u := range urls {
		resp := &model.AllUserURLsResponse{
			ShortURL:    fmt.Sprintf("%s/%s", s.Config.BaseURL, u.ID),
			OriginalURL: u.BaseURL,
		}
		responseURLs = append(responseURLs, resp)
	}

	response, err := json.Marshal(responseURLs)
	if s.handleErrorOrStatus(w, err, http.StatusInternalServerError) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	s.handleErrorOrStatus(w, err, http.StatusInternalServerError)
}

// handlePingStore ...
func (s *Server) handlePingStore(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cancelCoolDown)
	defer cancel()

	if err := s.Store.Ping(ctx); err != nil {
		s.handleErrorOrStatus(w, fmt.Errorf("handlePingStore: %w", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// handleURLBulkCreate ...
func (s *Server) handleURLBulkCreate(w http.ResponseWriter, r *http.Request) {
	var (
		data = []*model.BulkCreateURLRequest{}
		urls = []*model.URL{}
	)
	body, err := io.ReadAll(r.Body)
	if s.handleErrorOrStatus(w, err, http.StatusInternalServerError) {
		return
	}

	if err := json.Unmarshal(body, &data); s.handleErrorOrStatus(w, err, http.StatusInternalServerError) {
		return
	}

	var userID string
	if user := r.Context().Value(middleware.UserCTXName); user != nil {
		userID = user.(string)
	} else {
		userID = middleware.UserIDDefaultValue
	}

	for _, v := range data {
		u, err := model.NewURL(v.OriginalURL, userID, v.CorrelationID)
		if s.handleErrorOrStatus(w, err, http.StatusInternalServerError) {
			return
		}
		urls = append(
			urls,
			u,
		)
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(len(urls))*cancelCoolDown)
	defer cancel()

	resp, err := s.Store.URLsBulkCreate(ctx, urls)
	for _, v := range resp {
		id := v.ShortURL
		v.ShortURL = fmt.Sprintf("%s/%s", s.Config.BaseURL, id)
	}

	if s.handleErrorOrStatus(w, err, http.StatusBadRequest) {
		return
	}

	body, err = json.Marshal(resp)
	if s.handleErrorOrStatus(w, err, http.StatusInternalServerError) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err = w.Write(body); err != nil {
		log.Fatal(err)
	}
}
