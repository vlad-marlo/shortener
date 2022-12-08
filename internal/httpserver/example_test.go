package httpserver_test

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/vlad-marlo/shortener/internal/httpserver"
	"github.com/vlad-marlo/shortener/internal/store/inmemory"
)

func ExampleNew() {
	// parse server config
	config := httpserver.NewConfig()
	// create some storage; for example inmemory
	storage := inmemory.New()
	defer func() {
		// always close storage
		if err := storage.Close(); err != nil {
			// ...
		}
	}()
	server := httpserver.New(config, storage, logrus.NewEntry(logrus.New()))
	// always close server
	defer server.Close()

	go func() {
		// start your http server
		_ = http.ListenAndServe("localhost:8080", server.Router)
	}()
}
