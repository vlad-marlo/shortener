package httpserver

import (
	"flag"

	"github.com/caarlos0/env/v6"

	"github.com/vlad-marlo/shortener/internal/store"
)

// Config ...
type Config struct {
	BindAddr    string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL     string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FilePath    string `env:"FILE_STORAGE_PATH"`
	Database    string `env:"DATABASE_DSN"`
	HTTPS       bool   `env:"ENABLE_HTTPS"`
	StorageType string
}

// NewConfig return pointer to config with params. Empty params will be set by default
func NewConfig() (*Config, error) {
	c := &Config{}
	if err := env.Parse(c); err != nil {
		return nil, err
	}
	flag.StringVar(&c.BindAddr, "a", c.BindAddr, "server will be started with this url")
	flag.StringVar(&c.BaseURL, "b", c.BaseURL, "url will be used in generation of shorten url")
	flag.StringVar(&c.FilePath, "f", c.FilePath, "path to storage path")
	flag.StringVar(&c.Database, "d", c.Database, "db dns")
	flag.BoolVar(&c.HTTPS, "s", c.HTTPS, "if true, server will start with https protocol")
	flag.Parse()

	if c.Database != "" {
		c.StorageType = store.SQLStore
	} else if c.FilePath != "" {
		c.StorageType = store.FileBasedStorage
	} else {
		c.StorageType = store.InMemoryStorage
	}
	return c, nil
}
