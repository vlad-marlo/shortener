package main

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vlad-marlo/logger"

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
			storage, err := initStorage(tc.cfg, logrus.NewEntry(logger.WithOpts(logger.WithOutput(io.Discard))))
			require.NoError(t, err)

			if tc.cfg.Database == "" {
				require.NoError(t, err, fmt.Sprintf("init storage: %v", err))
				assert.IsType(t, storage, tc.wantType)
			}

			if tc.cfg.StorageType == store.FileBasedStorage && tc.cfg.FilePath != "" {
				require.NoError(t, os.Remove(tc.cfg.FilePath))
			}
		})
	}
}
