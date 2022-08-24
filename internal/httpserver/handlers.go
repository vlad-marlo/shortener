package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

func (s *Server) handleURLGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	id := chi.URLParam(r, "id")

	url, err := s.Store.GetByID(id)
	if err != nil {
		s.handleErrorOrStatus(w, errors.New("where is no url with that id"), http.StatusNotFound)
		return
	}

	w.Header().Set("Location", url.BaseURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *Server) handleURLCreate(w http.ResponseWriter, r *http.Request) {
	// setting up response meta info
	w.Header().Set("Content-Type", "text/plain")

	data, err := io.ReadAll(r.Body)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	if s.handleErrorOrStatus(w, err, http.StatusInternalServerError) {
		return
	}
	if len(data) == 0 {
		s.handleErrorOrStatus(w, ErrIncorrectRequestBody, http.StatusBadRequest)
		return
	}

	u, err := model.NewURL(string(data))
	if s.handleErrorOrStatus(w, err, http.StatusBadRequest) {
		return
	}

	if err = s.Store.Create(u); s.handleErrorOrStatus(w, err, http.StatusBadRequest) {
		return
	}

	// generate full url like <base service url>/<url identificator>
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(fmt.Sprintf("%s/%s", s.Config.BaseURL, u.ID)))
	s.handleErrorOrStatus(w, err, http.StatusInternalServerError)
}

func (s *Server) handleURLCreateJSON() http.HandlerFunc {
	type response struct {
		ResultURL string `json:"result"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		u := &model.URL{}
		req, err := io.ReadAll(r.Body)
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(r.Body)

		if s.handleErrorOrStatus(w, err, http.StatusInternalServerError) {
			return
		}
		err = json.Unmarshal(req, u)
		if s.handleErrorOrStatus(w, err, http.StatusInternalServerError) {
			return
		}

		if err = u.ShortURL(); s.handleErrorOrStatus(w, err, http.StatusBadRequest) {
			return
		}
		if err = s.Store.Create(u); s.handleErrorOrStatus(w, err, http.StatusBadRequest) {
			return
		}

		resp := response{
			ResultURL: fmt.Sprintf("%s/%s", s.Config.BaseURL, u.ID),
		}
		res, err := json.Marshal(resp)
		if s.handleErrorOrStatus(w, err, http.StatusInternalServerError) {
			return
		}
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(res)
		s.handleErrorOrStatus(w, err, http.StatusInternalServerError)
	}
}
