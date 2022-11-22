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
type Fields interface {
	map[string]interface{} | logrus.Fields
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

			switch {
			case rw.code >= 500:
				logger.WithFields(logrus.Fields{
					"duration":   dur.String(),
					"code":       rw.code,
					"request_id": middleware.GetReqID(r.Context()),
				}).Error(http.StatusText(rw.code))
			case rw.code >= 400:
				logger.WithFields(logrus.Fields{
					"duration":   dur.String(),
					"code":       rw.code,
					"request_id": middleware.GetReqID(r.Context()),
				}).Debug(http.StatusText(rw.code))
			case rw.code >= 100:
				logger.WithFields(logrus.Fields{
					"duration":   dur.String(),
					"code":       rw.code,
					"request_id": middleware.GetReqID(r.Context()),
				}).Trace(http.StatusText(rw.code))
			default:
				logger.WithFields(logrus.Fields{
					"duration":   dur.String(),
					"code":       rw.code,
					"request_id": middleware.GetReqID(r.Context()),
				}).Warn(http.StatusText(rw.code))
			}
		})
	}
}
