package main

import (
	"github.com/vlad-marlo/shortener/internal/httpserver"
)

func main() {
	addr := "http://localhost:8080"
	s := httpserver.New(addr)
	s.ListenAndServe()
}
