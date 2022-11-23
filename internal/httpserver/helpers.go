package httpserver

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/vlad-marlo/shortener/internal/httpserver/middleware"
)

// HandleErrorOr400 return true and handle error if err is not nil
func (s *Server) handleErrorOrStatus(w http.ResponseWriter, err error, fields map[string]interface{}, status int) bool {
	if err != nil {
		var lvl logrus.Level
		switch {
		case status >= 500:
			lvl = logrus.ErrorLevel
		case status >= 400:
			lvl = logrus.DebugLevel
		case status >= 100:
			lvl = logrus.TraceLevel
		default:
			lvl = logrus.WarnLevel
		}

		s.logger.WithFields(fields).Log(lvl, fmt.Sprintf("%v", err))
		w.WriteHeader(status)
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
