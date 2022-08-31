package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Print(err)
		}
	}()

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

func (s *Server) handleURLCreateJSON(w http.ResponseWriter, r *http.Request) {
	req, err := io.ReadAll(r.Body)
	defer func() {
		if err = r.Body.Close(); err != nil {
			log.Print(err)
		}
	}()

	if s.handleErrorOrStatus(w, err, http.StatusInternalServerError) {
		return
	}
	u := &model.URL{}
	if err = json.Unmarshal(req, u); s.handleErrorOrStatus(w, err, http.StatusInternalServerError) {
		return
	}

	if err = u.ShortURL(); s.handleErrorOrStatus(w, err, http.StatusBadRequest) {
		return
	}
	if err = s.Store.Create(u); s.handleErrorOrStatus(w, err, http.StatusBadRequest) {
		return
	}

	resp := model.ResultResponse{
		Result: fmt.Sprintf("%s/%s", s.Config.BaseURL, u.ID),
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
