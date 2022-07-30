package httpserver

import (
	"net/http"
	"strings"

	"github.com/vlad-marlo/shortener/internal/store/model"
)

func handleUrlGetCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := strings.TrimPrefix(r.URL.Path, "/")
		url := model.GetUrlById(id)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case http.MethodPost:
		w.WriteHeader(http.StatusCreated)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
