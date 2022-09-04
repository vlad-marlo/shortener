package store

import (
	"context"

	"github.com/vlad-marlo/shortener/internal/store/model"
)

const (
	InMemoryStorage  string = "inmemory"
	FileBasedStorage string = "file-based"
	SQLStore         string = "sql-store"
)

type Store interface {
	Create(context.Context, *model.URL) error
	GetByID(context.Context, string) (*model.URL, error)
	GetAllUserURLs(context.Context, string) ([]*model.URL, error)
	Ping(context.Context) error
}
