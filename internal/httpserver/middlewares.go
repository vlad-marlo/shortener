package httpserver

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Write ...
func (w gzipWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}

// gzipCompression middleware func
func (s *Server) gzipCompression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("content-encoding"), "gzip") {
			reader, err := gzip.NewReader(r.Body)
			defer func() {
				if err = reader.Close(); err != nil {
					log.Print(err)
				}
				if err = r.Body.Close(); err != nil {
					log.Print(err)
				}
			}()
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			r.Body = reader
		}
		if !strings.Contains(r.Header.Get("accept-encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		defer func() {
			if err = gz.Close(); err != nil {
				log.Print(err)
			}
		}()
		w.Header().Set("content-encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
