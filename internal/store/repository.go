package store

import (
	"context"

	"github.com/vlad-marlo/shortener/internal/store/model"
)

// const ...
const (
	// InMemoryStorage ...
	InMemoryStorage string = "in-memory"
	// FileBasedStorage ...
	FileBasedStorage string = "file-based"
	// SQLStore ...
	SQLStore string = "sql-store"
)

// Store ...
type Store interface {
	// Ping ...
	Ping(context.Context) error
	// Close ...
	Close() error
	// Create ...
	Create(context.Context, *model.URL) error
	// GetByID ...
	GetByID(context.Context, string) (*model.URL, error)
	// GetAllUserURLs ...
	GetAllUserURLs(context.Context, string) ([]*model.URL, error)
	// URLsBulkCreate ...
	URLsBulkCreate(context.Context, []*model.URL) ([]*model.BatchCreateURLsResponse, error)
	// URLsBulkDelete ...
	URLsBulkDelete([]string, string) error
}
