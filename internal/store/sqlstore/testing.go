package sqlstore

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestStore(t *testing.T) (*SQLStore, func(t *testing.T)) {
	t.Helper()
	dns := os.Getenv("TEST_DB_URI")
	if dns == "" {
		t.Skip("db uri is not provided in TEST_DB_URI os arg")
	}

	db, err := sql.Open("postgres", dns)
	require.NoError(t, err, "connect to db: %w")

	if err = db.Ping(); err != nil {
		t.Skipf("db is not accessible: %v", err)
	}

	l, _ := zap.NewProduction()
	storage, err := New(context.Background(), dns, l, db)
	require.NoError(t, err, fmt.Sprintf("init db storage: %v", err))

	return storage, func(t *testing.T) {
		_, err = storage.DB.Exec("TRUNCATE urls CASCADE;")
		assert.NoError(t, err, fmt.Sprintf("truncate db: %v", err))
		require.NoError(t, storage.Close(), "close storage")
	}
}
