package httpserver

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	BindAddr    string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL     string `env:"BASE_URL"`
	StorageType string
	FilePath    string `env:"FILE_STORAGE_PATH" envDefault:"data.json"`
}

// NewConfig return pointer to config with params. Empty params will be set by default
func NewConfig(StorageType string) *Config {
	c := &Config{
		StorageType: StorageType,
	}
	if err := env.Parse(c); err != nil {
		log.Fatal(err)
	}
	if c.BaseURL == "" {
		c.BaseURL = fmt.Sprintf("http://%s", c.BindAddr)
	}
	return c
}
