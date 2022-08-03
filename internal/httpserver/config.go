package httpserver

type Config struct {
	BindAddr    string
	StorageType string
}

// NewConfig return pointer to config with default params
func NewConfig() *Config {
	return &Config{
		BindAddr:    "localhost:8080",
		StorageType: "inmemory",
	}
}
