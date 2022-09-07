package main

import (
	"log"

	"github.com/vlad-marlo/shortener/internal/httpserver"
)

func main() {
	config := httpserver.NewConfig()
	if err := httpserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
