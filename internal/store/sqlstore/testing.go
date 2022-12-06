package sqlstore

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vlad-marlo/logger"

	"github.com/vlad-marlo/shortener/internal/store"
)

func TestStore(t *testing.T) store.Store {
	t.Helper()
	dns := os.Getenv("TEST_DB_URI")
	if dns == "" {
		t.Skip("db uri is not provided in TEST_DB_URI os arg")
	}
	storage, err := New(context.Background(), dns, logger.WithOpts(logger.WithOutput(io.Discard)))
	require.NoError(t, err, fmt.Sprintf("init db storage: %v", err))
	return storage
}
