package httpserver

type config struct {
	BindAddr    string
	StorageType string
}

// return pointer to config with default params
func NewConfig() *config {
	return &config{
		BindAddr:    "localhost:8080",
		StorageType: "inmemory",
	}
}