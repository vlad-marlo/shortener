package httpserver

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/vlad-marlo/shortener/internal/store/inmemory"
)

func ExampleNew() {
	// create some storage; for example inmemory
	storage := inmemory.New()
	defer func() {
		// always close storage
		if err := storage.Close(); err != nil {
			// ...
		}
	}()
	l, _ := zap.NewProduction()
	var srv service
	server := New(srv, l)
	// always close server
	defer func() {
		if err := server.Close(); err != nil {
			// handle error
		}
	}()

	go func() {
		// start your http server
		_ = http.ListenAndServe("localhost:8080", server.Router)
	}()
}
