package sqlstore

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vlad-marlo/shortener/internal/store/model"
)

func TestTestStore(t *testing.T) {
	const DBEnvKey = "TEST_DB_URI"
	t.Run("check skips", func(t *testing.T) {
		uri := os.Getenv(DBEnvKey)
		err := os.Unsetenv(DBEnvKey)
		require.NoError(t, err)
		defer require.NoError(t, os.Setenv(DBEnvKey, uri))
		TestStore(t)
		if uri == "" {
			t.Fatal("test must be skipped")
		}
	})
}

func TestSQLStore_Ping(t *testing.T) {
	s, td := TestStore(t)
	defer td(t)
	require.NoError(t, s.Ping(context.Background()))
}

func TestSQLStore_GetByID(t *testing.T) {
	s, td := TestStore(t)
	defer td(t)

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
			u    string
			user string
		}
		want struct {
			err assert.ErrorAssertionFunc
		}
	)
	store, teardown := TestStore(t)
	defer teardown(t)
	tt := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive case #1",
			args: args{
				u:    "asdfg",
				user: "1",
			},
			want: want{
				err: assert.NoError,
			},
		},
		{
			name: "negative case #1",
			args: args{
				u:    "asdfg",
				user: "1",
			},
			want: want{
				err: assert.Error,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			u, err := model.NewURL(tc.args.u, tc.args.user)
			require.NoError(t, err, "create new model url")
			err = store.Create(context.Background(), u)
			tc.want.err(t, err)
		})
	}
}
