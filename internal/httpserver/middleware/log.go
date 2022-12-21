package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// codeWriter ...
type codeWriter struct {
	http.ResponseWriter
	code int
}

// WriteHeader ...
func (c *codeWriter) WriteHeader(code int) {
	c.ResponseWriter.WriteHeader(code)
	c.code = code
}

// Write ...
func (c *codeWriter) Write(b []byte) (int, error) {
	return c.ResponseWriter.Write(b)
}

// Logger ...
func Logger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &codeWriter{w, http.StatusOK}

			start := time.Now()
			next.ServeHTTP(rw, r)
			dur := time.Since(start)

			fields := []zap.Field{
				zap.String("duration", dur.String()),
				zap.Int("code", rw.code),
				zap.String("request_id", middleware.GetReqID(r.Context())),
			}

			switch {
			case rw.code >= 500:
				logger.Error(http.StatusText(rw.code), fields...)
			case rw.code >= 400:
				logger.Debug(http.StatusText(rw.code), fields...)
			case rw.code >= 100:
				logger.Debug(http.StatusText(rw.code), fields...)
			default:
				logger.Warn(http.StatusText(rw.code), fields...)
			}
		})
	}
}
