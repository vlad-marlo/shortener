package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"sync"

	"github.com/caarlos0/env/v6"

	"github.com/vlad-marlo/shortener/internal/store"
)

// Config ...
type Config struct {
	ConfigFile string
	BindAddr   string `env:"SERVER_ADDRESS" json:"server_address"`
	BaseURL    string `env:"BASE_URL" json:"base_url"`
	FilePath   string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	Database   string `env:"DATABASE_DSN" json:"database_dsn"`
	HTTPS      bool   `env:"ENABLE_HTTPS" json:"enable_https"`
	GRPC       bool   `env:"ENABLE_GRPC" json:"enable_grpc"`
	GRPCAddr   string `env:"GRPC_ADDRESS" json:"grpc_address"`
	TrustedIP  string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`

	StorageType string
	IP          net.IP
}

// defaultBindAddr ...
const defaultBindAddr = "localhost:8080"

// defaultBaseURL ...
const defaultBaseURL = "http://localhost:8080"

var c *Config
var once sync.Once

// Get return pointer to config with params. Empty params will be set by default
func Get() *Config {
	once.Do(func() {
		c := &Config{}
		if err := env.Parse(c); err != nil {
			log.Fatalf("parse env: %v", err)
		}
		defer c.setDefaultValues()
		flag.StringVar(&c.BindAddr, "a", c.BindAddr, "server will be started with this url")
		flag.StringVar(&c.BaseURL, "b", c.BaseURL, "url will be used in generation of shorten url")
		flag.StringVar(&c.FilePath, "f", c.FilePath, "path to storage path")
		flag.StringVar(&c.Database, "d", c.Database, "db dns")
		flag.BoolVar(&c.HTTPS, "s", c.HTTPS, "if true, server will start with https protocol")
		flag.BoolVar(&c.HTTPS, "g", c.GRPC, "if true, server will start with grpc")
		flag.StringVar(&c.ConfigFile, "c", c.ConfigFile, "server will use this settings")
		flag.StringVar(&c.TrustedIP, "t", c.TrustedIP, "trusted ip in CIDR presentation")
		flag.Parse()

		if c.Database != "" {
			c.StorageType = store.SQLStore
		} else if c.FilePath != "" {
			c.StorageType = store.FileBasedStorage
		} else {
			c.StorageType = store.InMemoryStorage
		}

		if err := c.parseFile(); err != nil {
			log.Fatalf("parse file: %v", err)
		}
	})
	return c
}

// parseFile ...
func (c *Config) parseFile() error {
	if c.ConfigFile == "" {
		return nil
	}

	ex, err := os.Executable()
	if err != nil {
		return fmt.Errorf("os executable: %w", err)
	}

	var data []byte
	data, err = os.ReadFile(path.Join(ex, c.FilePath))
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil
		}
		return fmt.Errorf("os readfile: %w", err)
	}
	newConfig := &Config{}

	err = json.Unmarshal(data, newConfig)
	if err != nil {
		return fmt.Errorf("unmarshal json: %w", err)
	}

	// setting all values which are not provided
	if c.Database == "" {
		c.Database = newConfig.Database
	}
	if c.BaseURL == "" {
		c.BaseURL = newConfig.BaseURL
	}
	if c.FilePath == "" {
		c.FilePath = newConfig.FilePath
	}
	if c.BindAddr == "" {
		c.BindAddr = newConfig.BindAddr
	}
	if c.GRPCAddr == "" {
		c.GRPCAddr = newConfig.GRPCAddr
	}
	if !c.HTTPS {
		c.HTTPS = newConfig.HTTPS
	}
	if !c.GRPC {
		c.GRPC = newConfig.GRPC
	}
	if c.TrustedIP == "" {
		c.TrustedIP = newConfig.TrustedIP
	}

	return nil
}

// setDefaultValues ...
func (c *Config) setDefaultValues() {
	switch {
	case c.Database != "":
		c.StorageType = store.SQLStore
	case c.FilePath != "":
		c.StorageType = store.FileBasedStorage
	default:
		c.StorageType = store.InMemoryStorage
	}

	if c.BindAddr == "" {
		c.BindAddr = defaultBindAddr
	}
	if c.BaseURL == "" {
		c.BaseURL = defaultBaseURL
	}
}
