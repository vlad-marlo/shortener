package httpserver_test

import (
	"github.com/vlad-marlo/shortener/internal/httpserver"
	"net/http"
)

func ExampleStart() {
	server, err := httpserver.Start(nil, nil, nil)
	if err != nil {
		// ...
	}
	go func() {
		// start your http server
		_ = http.ListenAndServe("localhost:8080", server.Router)
	}()
}
