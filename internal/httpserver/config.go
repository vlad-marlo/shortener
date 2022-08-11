package httpserver

type Config struct {
	BindAddr    string
	StorageType string
}

// NewConfig return pointer to config with params. Empty params will be set by default
func NewConfig(BindAddr string, StorageType string) *Config {
	c := &Config{
		BindAddr:    "localhost:8080",
		StorageType: "inmemory",
	}
	if BindAddr != "" {
		c.BindAddr = BindAddr
	}
	if StorageType != "" {
		c.StorageType = StorageType
	}
	return c
}
