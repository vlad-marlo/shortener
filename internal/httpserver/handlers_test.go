package httpserver

import (
	"net/http"
	"testing"

	"github.com/vlad-marlo/shortener/internal/store"
)

func TestServer_handleURLGetCreate(t *testing.T) {
	type fields struct {
		Server http.Server
		Store  store.Store
		Config *Config
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				srv:    tt.fields.Server,
				Store:  tt.fields.Store,
				Config: tt.fields.Config,
			}
			s.handleURLGetCreate(tt.args.w, tt.args.r)
		})
	}
}
