package sqlstore

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLStore_Close(t *testing.T) {
	store := TestStore(t)
	require.NoError(t, store.Close())
}

func TestTestStore(t *testing.T) {
	const DBEnvKey = "TEST_DB_URI"
	t.Run("positive case", func(t *testing.T) {
		require.NoError(t, TestStore(t).Close())
	})
	t.Run("check skips", func(t *testing.T) {
		uri := os.Getenv(DBEnvKey)
		require.NoError(t, os.Unsetenv(DBEnvKey))
		defer require.NoError(t, os.Setenv(DBEnvKey, uri))
		TestStore(t)
		t.Fatal("test must be skipped")
	})
}

func TestSQLStore_Ping(t *testing.T) {
	store := TestStore(t)
	require.NoError(t, store.Ping(context.Background()))
	require.NoError(t, store.Close())
}
