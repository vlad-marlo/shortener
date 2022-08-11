package main

import (
	"log"

	"github.com/vlad-marlo/shortener/internal/httpserver"
)

func main() {
	config := httpserver.NewConfig("localhost:8080", "inmemory")
	s := httpserver.New(config)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
