package httpserver

import (
	"log"
	"net/http"

	"github.com/vlad-marlo/shortener/internal/httpserver/middleware"
)

// HandleErrorOr400 return true and handle error if err is not nil
func (s *Server) handleErrorOrStatus(w http.ResponseWriter, err error, status int) bool {
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), status)
	}
	return err != nil
}

func getUserFromRequest(r *http.Request) string {
	user := r.Context().Value(middleware.UserCTXName)
	if user == nil {
		return middleware.UserIDDefaultValue
	}
	return user.(string)
}
