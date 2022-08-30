package main

import (
	"log"

	"github.com/vlad-marlo/shortener/internal/httpserver"
)

func main() {
	config := httpserver.NewConfig()
	s := httpserver.New(config)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
