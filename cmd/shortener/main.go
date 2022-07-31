package main

import (
	"github.com/vlad-marlo/shortener/internal/httpserver"
)

func main() {
	s := httpserver.New("localhost:8080")
	s.ListenAndServe()
}
