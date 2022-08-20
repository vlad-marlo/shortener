package main

import (
	"log"

	"github.com/vlad-marlo/shortener/internal/httpserver"
)

func main() {
	config, err := httpserver.NewConfig("localhost:8080", "inmemory")
	if err != nil {
		log.Fatal(err)
	}
	s := httpserver.New(config)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
