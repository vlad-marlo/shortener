package httpserver

import (
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	BindAddr    string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL     string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	StorageType string
	FilePath    string `env:"FILE_STORAGE_PATH" envDefault:"data.json"`
}

// NewConfig return pointer to config with params. Empty params will be set by default
func NewConfig(StorageType, serverAddr, baseURL, fileStoragePath string) *Config {
	c := &Config{
		StorageType: StorageType,
	}
	if err := env.Parse(c); err != nil {
		log.Fatal(err)
	}
	if baseURL != "http://localhost:8080" {
		c.BaseURL = baseURL
	}
	if serverAddr != "localhost:8080" {
		c.BindAddr = serverAddr
	}
	if fileStoragePath != "data.json" {
		c.FilePath = fileStoragePath
	}
	return c
}
