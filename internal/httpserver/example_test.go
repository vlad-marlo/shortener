package httpserver_test

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/vlad-marlo/shortener/internal/httpserver"
	"github.com/vlad-marlo/shortener/internal/store/inmemory"
)

func ExampleNew() {
	// parse server config
	config, err := httpserver.NewConfig()
	// always check error
	if err != nil {
	}
	// create some storage; for example inmemory
	storage := inmemory.New()
	defer func() {
		// always close storage
		if err := storage.Close(); err != nil {
			// ...
		}
	}()
	l, _ := zap.NewProduction()
	server := httpserver.New(config, storage, l)
	// always close server
	defer func() {
		if err := server.Close(); err != nil {
			// some err handling
		}
	}()

	go func() {
		// start your http server
		_ = http.ListenAndServe("localhost:8080", server.Router)
	}()
}
