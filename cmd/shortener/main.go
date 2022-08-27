package main

import (
	"flag"
	"log"

	"github.com/vlad-marlo/shortener/internal/httpserver"
	"github.com/vlad-marlo/shortener/internal/store"
)

var (
	fileStoragePath string
	serverAddr      string
	baseURL         string
)

// init setting up string path arguments
func init() {
	flag.StringVar(&serverAddr, "a", "localhost:8080", "server will be started with this url")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "url will be used in generation of shorten url")
	flag.StringVar(&fileStoragePath, "f", "data.json", "path to storage path")
}

func main() {
	flag.Parse()
	config := httpserver.NewConfig(store.FileBasedStorage, serverAddr, baseURL, fileStoragePath)
	s := httpserver.New(config)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
