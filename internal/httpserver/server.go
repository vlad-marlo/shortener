package httpserver

import (
	"log"
	"net/http"

	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/inmemory"
)

type Server struct {
	http.Server

	Store  store.Store
	Config *config
}

func New(config *config) *Server {
	s := &Server{
		Config: config,
	}
	s.Addr = s.Config.BindAddr
	s.routes()
	if err := s.configureStore(); err != nil {
		log.Fatal(err)
	}
	return s
}

func (s *Server) routes() {
	http.HandleFunc("/", s.handleURLGetCreate)
}

func (s *Server) configureStore() error {
	switch s.Config.StorageType {
	case "inmemory":
		s.Store = inmemory.New()

	default:
		return ErrIncorrectStoreType
	}

	return nil
}
