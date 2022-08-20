package httpserver

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	BindAddr    string `env:"SERVER_ADDRESS"`
	BaseURL     string `env:"BASE_URL"`
	StorageType string
}

// NewConfig return pointer to config with params. Empty params will be set by default
func NewConfig(BindAddr string, StorageType string) (*Config, error) {
	c := &Config{
		StorageType: "inmemory",
	}
	if err := env.Parse(&c); err != nil {
		return c, err
	}
	if c.BindAddr == "" {
		if BindAddr != "" {
			c.BindAddr = BindAddr
		} else {
			c.BindAddr = "localhost:8080"
		}
	}
	if c.BaseURL == "" {
		c.BaseURL = fmt.Sprintf("http://%s", c.BindAddr)
	}
	if c.StorageType != "" {
		c.StorageType = StorageType
	}
	return c, nil
}
