package httpserver

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/vlad-marlo/shortener/internal/store"
)

type Config struct {
	BindAddr    string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL     string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FilePath    string `env:"FILE_STORAGE_PATH"`
	StorageType string
}

// NewConfig return pointer to config with params. Empty params will be set by default
func NewConfig() *Config {
	c := &Config{}
	if err := env.Parse(c); err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&c.BindAddr, "a", c.BindAddr, "server will be started with this url")
	flag.StringVar(&c.BaseURL, "b", c.BaseURL, "url will be used in generation of shorten url")
	flag.StringVar(&c.FilePath, "f", c.FilePath, "path to storage path")
	flag.Parse()

	if c.FilePath == "" {
		c.StorageType = store.InMemoryStorage
	}
	return c
}
