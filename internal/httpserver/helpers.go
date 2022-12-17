package httpserver

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/vlad-marlo/shortener/internal/httpserver/middleware"
)

// HandleErrorOr400 return true and handle error if err is not nil
func (s *Server) handleErrorOrStatus(w http.ResponseWriter, err error, fields []zap.Field, status int) bool {
	if err != nil {
		var lvl zapcore.Level
		switch {
		case status >= 500:
			lvl = zapcore.ErrorLevel
		case status >= 400:
			lvl = zapcore.DebugLevel
		case status >= 100:
			lvl = zapcore.DebugLevel
		default:
			lvl = zapcore.WarnLevel
		}

		s.logger.Log(lvl, fmt.Sprintf("%v", err), fields...)
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
