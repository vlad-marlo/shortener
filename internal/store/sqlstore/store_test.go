package sqlstore

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vlad-marlo/shortener/internal/store"
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
	var (
		u1, u2, u3 *model.URL
		err        error
	)
	u1, err = model.NewURL("https://xd.com", "marlo")
	assert.NoError(t, err)
	u2, err = model.NewURL("https://ya.com", "marlo")
	assert.NoError(t, err)
	u3, err = model.NewURL("https://gooogle.com", "marlo")
	assert.NoError(t, err)
	create := []*model.URL{
		u1,
		u2,
		u3,
	}
	for _, u := range create {
		require.NoError(t, s.Create(context.Background(), u))
	}
	require.NoError(t, s.URLsBulkDelete([]string{u3.ID}, "marlo"))

	tt := []struct {
		name     string
		args     string
		wantErr  assert.ErrorAssertionFunc
		exactErr error
	}{
		{
			name:     "negative",
			args:     uuid.New().String(),
			wantErr:  assert.Error,
			exactErr: nil,
		},
		{
			name:    "positive #1",
			args:    u1.ID,
			wantErr: assert.NoError,
		},
		{
			name:    "positive #2",
			args:    u2.ID,
			wantErr: assert.NoError,
		},
		{
			name:     "negative is deleted",
			args:     u3.ID,
			wantErr:  assert.Error,
			exactErr: store.ErrIsDeleted,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.GetByID(context.Background(), tc.args)
			tc.wantErr(t, err)
			if tc.exactErr != nil {
				assert.ErrorIs(t, err, tc.exactErr)
			}
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
	storage, teardown := TestStore(t)
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
			err = storage.Create(context.Background(), u)
			tc.want.err(t, err)
		})
	}
}
