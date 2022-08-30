package httpserver

import "net/http"

// HandleErrorOr400 return true and handle error if err is not nil
func (s *Server) handleErrorOrStatus(w http.ResponseWriter, err error, status int) bool {
	if err != nil {
		http.Error(w, err.Error(), status)
	}
	return err != nil
}
