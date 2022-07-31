package httpserver

import (
	"io"
	"net/http"
	"strings"

	"github.com/vlad-marlo/shortener/internal/store/model"
)

func (s *Server) handleUrlGetCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		w.WriteHeader(http.StatusTemporaryRedirect)
		id := strings.TrimPrefix(r.URL.Path, "/")
		if id == "" {
			http.Error(w, "The path argument is missing", http.StatusBadRequest)
			return
		}

		url, ok := s.Store.GetById(id)
		if !ok {
			http.Error(w, "Where is no url with that id!", http.StatusNotFound)
			return
		}

		w.Write([]byte(url))
		return

	case http.MethodPost:
		d, err := io.ReadAll(r.Body)
		w.WriteHeader(http.StatusCreated)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		u := model.NewUrl(string(d))

		if err := s.Store.Create(u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(u.ID.String()))
		return

	default:
		http.Error(w, "Only POST and GET are allowed!", http.StatusMethodNotAllowed)
		return
	}
}
