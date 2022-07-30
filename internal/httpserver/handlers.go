package httpserver

import (
	"net/http"
	"strings"
)

func (s *Server) handleUrlGetCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := strings.TrimPrefix(r.URL.Path, "/")
		url, ok := s.Store.GetById(id)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write([]byte(url))
		w.WriteHeader(http.StatusTemporaryRedirect)
	case http.MethodPost:
		w.WriteHeader(http.StatusCreated)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
