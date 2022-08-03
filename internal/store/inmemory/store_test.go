package inmemory

import (
	"reflect"
	"testing"

	"github.com/vlad-marlo/shortener/internal/store/model"
)

func TestStore_Create(t *testing.T) {
	type fields struct {
		urls []model.URL
	}
	tests := []struct {
		name    string
		fields  fields
		u       model.URL
		wantErr bool
	}{
		{
			name: "with empty urls",
			u: model.URL{
				ID:      "some-text",
				BaseURL: "https://google.com/",
			},
			fields: fields{
				urls: []model.URL{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				urls: tt.fields.urls,
			}
			if err := s.Create(tt.u); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStore_GetByID(t *testing.T) {
	type fields struct {
		urls []model.URL
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
			s := &Store{
				urls: tt.fields.urls,
			}
			got, err := s.GetByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_urlExists(t *testing.T) {
	type fields struct {
		urls []model.URL
	}
	type args struct {
		url model.URL
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				urls: tt.fields.urls,
			}
			if got := s.urlExists(tt.args.url); got != tt.want {
				t.Errorf("urlExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
