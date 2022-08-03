package httpserver

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/vlad-marlo/shortener/internal/store/model"
)

func (s *Server) handleURLGetCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	switch r.Method {

	case http.MethodGet:
		id := strings.TrimPrefix(r.URL.Path, "/")
		if id == "" {
			http.Error(w, "The path argument is missing", http.StatusBadRequest)
			return
		}

		url, err := s.Store.GetByID(id)
		if err != nil {
			http.Error(w, "Where is no url with that id!", http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", url.BaseURL)
		w.WriteHeader(http.StatusTemporaryRedirect)

		return

	case http.MethodPost:
		// setting up response meta info
		w.WriteHeader(http.StatusCreated)

		data, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if s.HandleErrorOr500(w, err) {
			return
		}

		u, err := model.NewURL(string(data))
		if s.HandleErrorOr500(w, err) {
			return
		}

		if err = s.Store.Create(u); s.HandleErrorOr500(w, err) {
			return
		}

		// generate full url alike <base service url>/<url identificator>
		_, err = w.Write([]byte(fmt.Sprintf("http://%s/%s", s.Config.BindAddr, u.ID)))
		s.HandleErrorOr500(w, err)
		return

	default:
		http.Error(w, "Only POST and GET are allowed!", http.StatusBadRequest)
		return
	}
}
