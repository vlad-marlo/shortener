package main

import (
	"github.com/vlad-marlo/shortener/internal/httpserver"
	"log"
)

func main() {
	if err := httpserver.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
