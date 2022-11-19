package middleware

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type codeWriter struct {
	http.ResponseWriter
	code int
}

func (c *codeWriter) WriteHeader(code int) {
	c.ResponseWriter.WriteHeader(code)
	c.code = code
}

func (c *codeWriter) Write(b []byte) (int, error) {
	return c.ResponseWriter.Write(b)
}

func Logger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &codeWriter{w, http.StatusOK}

			start := time.Now()
			next.ServeHTTP(rw, r)
			dur := time.Since(start)

			var lvl logrus.Level
			switch {
			case rw.code >= 500:
				lvl = logrus.ErrorLevel
			case rw.code >= 400:
				lvl = logrus.DebugLevel
			case rw.code >= 100:
				lvl = logrus.TraceLevel
			default:
				lvl = logrus.WarnLevel
			}

			logger.WithFields(map[string]interface{}{
				"duration":   dur,
				"code":       rw.code,
				"request_id": middleware.GetReqID(r.Context()),
			}).Log(lvl, http.StatusText(rw.code))
		})
	}
}
