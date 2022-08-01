package httpserver

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/vlad-marlo/shortener/internal/store/model"
)

func (s *Server) handleUrlGetCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	switch r.Method {

	case http.MethodGet:
		// settin up response meta info

		id := strings.TrimPrefix(r.URL.Path, "/")
		if id == "" {
			http.Error(w, "The path argument is missing", http.StatusBadRequest)
			return
		}

		url, err := s.Store.GetById(id)
		if err != nil {
			http.Error(w, "Where is no url with that id!", http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", url.BaseURL)
		w.WriteHeader(http.StatusTemporaryRedirect)

		return

	case http.MethodPost:
		// settin up response meta info
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

		// generate full url alike <base service url>/<url identificator>
		w.Write([]byte(fmt.Sprintf("http://%s/%s", s.Config.BindAddr, u.ID.String())))
		return

	default:
		http.Error(w, "Only POST and GET are allowed!", http.StatusMethodNotAllowed)
		return
	}
}
