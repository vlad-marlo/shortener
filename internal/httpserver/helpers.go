package httpserver

import (
	"log"
	"net/http"
)

// HandleErrorOr400 return true and handle error if err is not nil
func (s *Server) handleErrorOrStatus(w http.ResponseWriter, err error, status int) bool {
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), status)
	}
	return err != nil
}
