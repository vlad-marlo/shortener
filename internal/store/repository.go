package store

import "github.com/vlad-marlo/shortener/internal/store/model"

const (
	InMemoryStorage  string = "inmemory"
	FileBasedStorage string = "file-based"
)

type Store interface {
	Create(*model.URL) error
	GetByID(string) (*model.URL, error)
}
