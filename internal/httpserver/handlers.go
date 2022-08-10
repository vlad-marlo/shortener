package httpserver

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/vlad-marlo/shortener/internal/store/model"
)

func (s *Server) handleURLGet(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		s.HandleErrorOr400(w, errors.New("The path argument is missing"))
		return
	}

	url, err := s.Store.GetByID(id)
	if err != nil {
		s.HandleErrorOr400(w, errors.New("Where is no url with that id!"))
		return
	}

	w.Header().Set("Location", url.BaseURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *Server) handleURLCreate(w http.ResponseWriter, r *http.Request) {
	// setting up response meta info
	w.Header().Set("Content-Type", "text/plain")

	if r.URL.Path != "/" {
		s.HandleErrorOr400(w, ErrIncorrectUrlPath)
		return
	}

	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if s.HandleErrorOr400(w, err) {
		return
	}
	if len(data) == 0 {
		s.HandleErrorOr400(w, ErrIncorrectRequestBody)
		return
	}

	u, err := model.NewURL(string(data))
	if s.HandleErrorOr400(w, err) {
		return
	}

	if err = s.Store.Create(u); s.HandleErrorOr400(w, err) {
		return
	}

	// generate full url like <base service url>/<url identificator>
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(fmt.Sprintf("http://%s/%s", s.Config.BindAddr, u.ID)))
	s.HandleErrorOr400(w, err)
}

func (s *Server) handleURLGetCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		s.handleURLGet(w, r)

	case http.MethodPost:
		s.handleURLCreate(w, r)

	default:
		s.HandleErrorOr400(w, errors.New("Only POST and GET are allowed!"))
	}
}
