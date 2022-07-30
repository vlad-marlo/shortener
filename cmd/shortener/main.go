package main

import (
	"github.com/vlad-marlo/shortener/internal/httpserver"
)

func main() {
	s := httpserver.New(":8080")
	s.ListenAndServe()
}
