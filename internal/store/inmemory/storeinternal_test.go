package inmemory

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

var (
	googleURL = model.URL{
		ID:      "goo",
		BaseURL: "https://google.com/",
	}
	yandexURL = model.URL{
		ID:      "ya",
		BaseURL: "yandex.ru",
	}
)

func TestStore_Create(t *testing.T) {
	type fields struct {
		urls map[string]model.URL
	}
	tests := []struct {
		name    string
		fields  fields
		u       model.URL
		wantErr bool
	}{
		{
			name: "with empty urls",
			u:    googleURL,
			fields: fields{
				urls: map[string]model.URL{},
			},
			wantErr: false,
		},
		{
			name: "with duplicates",
			u:    googleURL,
			fields: fields{
				urls: map[string]model.URL{
					googleURL.ID: googleURL,
					yandexURL.ID: yandexURL,
				},
			},
			wantErr: true,
		},
		{
			name: "without duplicates",
			u:    googleURL,
			fields: fields{
				urls: map[string]model.URL{
					yandexURL.ID: yandexURL,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &Store{
				urls: tt.fields.urls,
			}
			err := s.Create(tt.u)
			if tt.wantErr {
				require.Error(t, err, "Create() error is nil")
				if err != store.ErrAlreadyExists {
					assert.NotContains(t, s.urls, tt.u)
				}
			} else {
				assert.NoError(t, err, "Create() error = %v", err)
				assert.Contains(t, s.urls, tt.u.ID)
			}
		})
	}
}

func TestStore_GetByID(t *testing.T) {
	type fields struct {
		urls map[string]model.URL
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.URL
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &Store{
				urls: tt.fields.urls,
			}
			got, err := s.GetByID(tt.args.id)
			if tt.wantErr {
				assert.Error(t, err, "GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else {
				assert.NoError(t, err, "There must be error got %v", err)
			}
			assert.True(t, assert.ObjectsAreEqual(got, tt.want), "GetByID() got = %v, want %v", got, tt.want)
		})
	}
}
