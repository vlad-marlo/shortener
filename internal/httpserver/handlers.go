package httpserver

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

func (s *Server) handleURLGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		s.HandleErrorOrStatus(w, errors.New("the path argument is missing"), http.StatusBadRequest)
		return
	}

	url, err := s.Store.GetByID(id)
	if err != nil {
		s.HandleErrorOrStatus(w, errors.New("where is no url with that id"), http.StatusNotFound)
		return
	}

	w.Header().Set("Location", url.BaseURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *Server) handleURLCreate(w http.ResponseWriter, r *http.Request) {
	// setting up response meta info
	w.Header().Set("Content-Type", "text/plain")

	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if s.HandleErrorOrStatus(w, err, http.StatusInternalServerError) {
		return
	}
	if len(data) == 0 {
		s.HandleErrorOrStatus(w, ErrIncorrectRequestBody, http.StatusBadRequest)
		return
	}

	u, err := model.NewURL(string(data))
	if s.HandleErrorOrStatus(w, err, http.StatusBadRequest) {
		return
	}

	if err = s.Store.Create(u); s.HandleErrorOrStatus(w, err, http.StatusBadRequest) {
		return
	}

	// generate full url like <base service url>/<url identificator>
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(fmt.Sprintf("http://%s/%s", s.Config.BindAddr, u.ID)))
	s.HandleErrorOrStatus(w, err, http.StatusInternalServerError)
}

func (s *Server) handleURLGetCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		s.handleURLGet(w, r)

	case http.MethodPost:
		s.handleURLCreate(w, r)

	default:
		s.HandleErrorOrStatus(w, errors.New("only POST and GET are allowed"), http.StatusMethodNotAllowed)
	}
}
