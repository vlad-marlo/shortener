package store

import (
	"context"

	"github.com/vlad-marlo/shortener/internal/store/model"
)

const (
	InMemoryStorage  string = "in-memory"
	FileBasedStorage string = "file-based"
	SQLStore         string = "sql-store"
)

type Store interface {
	Ping(context.Context) error
	Close() error
	Create(context.Context, *model.URL) error
	GetByID(context.Context, string) (*model.URL, error)
	GetAllUserURLs(context.Context, string) ([]*model.URL, error)
	URLsBulkCreate(context.Context, []*model.URL) ([]*model.BatchCreateURLsResponse, error)
	URLsBulkDelete(context.Context, []string, string) error
}
