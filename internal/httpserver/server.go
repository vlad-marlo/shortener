package httpserver

import (
	"net/http"

	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/inmemory"
)

type Server struct {
	http.Server
	Store store.Store
}

func New(addr string) *Server {
	s := &Server{
		Store: inmemory.New(),
	}
	s.Addr = addr
	return s
}

func (s *Server) routes() {
	http.HandleFunc("/", s.handleUrlGetCreate)
}
