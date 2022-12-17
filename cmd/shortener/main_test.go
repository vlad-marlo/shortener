package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/vlad-marlo/shortener/internal/httpserver"
	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/filebased"
	"github.com/vlad-marlo/shortener/internal/store/inmemory"
	"github.com/vlad-marlo/shortener/internal/store/sqlstore"
)

func TestInitStorage(t *testing.T) {
	tt := []struct {
		name     string
		cfg      *httpserver.Config
		wantType store.Store
	}{
		{
			name: "inmemory storage #1",
			cfg: &httpserver.Config{
				StorageType: store.InMemoryStorage,
			},
			wantType: inmemory.New(),
		},
		{
			name:     "inmemory storage #2",
			cfg:      &httpserver.Config{},
			wantType: inmemory.New(),
		},
		{
			name: "sql storage #1",
			cfg: &httpserver.Config{
				Database:    "postgresql://postgres:postgres@localhost:5432/shortner_test?sslmode=disable",
				StorageType: store.SQLStore,
			},
			wantType: &sqlstore.SQLStore{},
		},
		{
			name: "file storage #1",
			cfg: &httpserver.Config{
				FilePath:    "xd",
				StorageType: store.FileBasedStorage,
			},
			wantType: &filebased.Store{},
		},
	}

	defer func() {
		_ = os.Remove("xd")
	}()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			l, err := zap.NewProduction()
			require.NoError(t, err)
			storage, err := initStorage(tc.cfg, l)

			if tc.cfg.StorageType == store.SQLStore {
				if tc.cfg.Database == "" || errors.Is(err, store.ErrNotAccessible) {
					return
				}
				assert.IsType(t, storage, tc.wantType)
			}

			if tc.cfg.StorageType == store.FileBasedStorage && tc.cfg.FilePath != "" {
				require.NoError(t, os.Remove(tc.cfg.FilePath))
			}
		})
	}
}
