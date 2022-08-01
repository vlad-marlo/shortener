package main

import (
	"github.com/vlad-marlo/shortener/internal/httpserver"
)

func main() {
	config := httpserver.NewConfig()
	s := httpserver.New(config)
	s.ListenAndServe()
}
