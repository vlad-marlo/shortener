package httpserver

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/vlad-marlo/shortener/internal/store/model"
)

func (s *Server) handleUrlGetCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		// settin up response meta info
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Header().Set("Content-Type", "text/plain")

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
		w.WriteHeader(http.StatusCreated)

		d, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		u := model.NewUrl(string(d))

		if err := s.Store.Create(u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// generate full url like <base service url>/<url identificator>
		w.Write([]byte(fmt.Sprintf("%s/%s", s.Addr, u.ID.String())))
		return

	default:
		http.Error(w, "Only POST and GET are allowed!", http.StatusMethodNotAllowed)
		return
	}
}
