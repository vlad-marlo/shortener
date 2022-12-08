package sqlstore

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLStore_Close(t *testing.T) {
	s, teardown := TestStore(t)
	defer teardown(t)

	require.NoError(t, s.Close())
	require.Error(t, s.Close(), "close already closed storage")
}

func TestTestStore(t *testing.T) {
	const DBEnvKey = "TEST_DB_URI"
	t.Run("positive case", func(t *testing.T) {
		s, td := TestStore(t)
		defer td(t)

		require.NoError(t, s.Close())
	})
	t.Run("check skips", func(t *testing.T) {
		uri := os.Getenv(DBEnvKey)
		err := os.Unsetenv(DBEnvKey)
		require.NoError(t, err)
		defer func() {
			err = os.Setenv(DBEnvKey, uri)
			require.NoError(t, err)
		}()
		TestStore(t)
		t.Fatal("test must be skipped")
	})
}

func TestSQLStore_Ping(t *testing.T) {
	s, td := TestStore(t)
	defer td(t)
	require.NoError(t, s.Ping(context.Background()))
	require.NoError(t, s.Close())
}

func TestSQLStore_GetByID(t *testing.T) {
	s, td := TestStore(t)
	defer td(t)
	defer require.NoError(t, s.Close())

	tt := []struct {
		name    string
		args    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "negative",
			args:    uuid.New().String(),
			wantErr: assert.Error,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.GetByID(context.Background(), tc.args)
			tc.wantErr(t, err)
		})
	}
}

func TestSQLStore_Create(t *testing.T) {
	type (
		args struct {
		}
		want struct {
		}
	)
	store, teardown := TestStore(t)
	defer teardown(t)
	defer require.NoError(t, store.Close())
	tt := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive #1",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

		})
	}
}
