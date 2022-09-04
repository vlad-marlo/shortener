package store

import "github.com/vlad-marlo/shortener/internal/store/model"

const (
	InMemoryStorage  string = "inmemory"
	FileBasedStorage string = "file-based"
	SQLStore         string = "sql-store"
)

type Store interface {
	Create(*model.URL) error
	GetByID(string) (*model.URL, error)
	GetAllUserURLs(string) ([]*model.URL, error)
	Ping() error
}
