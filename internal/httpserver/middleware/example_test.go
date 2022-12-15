package middleware_test

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/vlad-marlo/shortener/internal/httpserver/middleware"
)

func ExampleLogger() {
	// init your handler
	handler := &http.ServeMux{}
	// init your logger
	logger, _ := zap.NewProduction()
	go func() {
		_ = http.ListenAndServe(
			"localhost:8080",
			// start server wrapping your handler with middleware
			middleware.Logger(logger)(handler),
		)
	}()
}

func ExampleGzipCompression() {
	// init your handler which you like
	handler := &http.ServeMux{}
	go func() {
		// start server
		_ = http.ListenAndServe(
			"localhost:8080",
			// use middleware to compress responses
			middleware.GzipCompression(handler),
		)
	}()
}
