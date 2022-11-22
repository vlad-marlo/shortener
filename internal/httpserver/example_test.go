package httpserver_test

import (
	"github.com/sirupsen/logrus"
	"github.com/vlad-marlo/shortener/internal/httpserver"
	"github.com/vlad-marlo/shortener/internal/store/inmemory"
	"net/http"
)

func ExampleNew() {
	// parse server config
	config := httpserver.NewConfig()
	// create some storage; for example inmemory
	storage := inmemory.New()
	server, err := httpserver.New(config, storage, logrus.New())
	if err != nil {
		// ...
	}
	go func() {
		// start your http server
		_ = http.ListenAndServe("localhost:8080", server.Router)
	}()
}
