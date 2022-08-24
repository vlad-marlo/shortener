package httpserver

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/filebased"
	"github.com/vlad-marlo/shortener/internal/store/inmemory"
)

const (
	InMemoryStorage  = "inmemory"
	FileBasedStorage = "filebased"
)

type Server struct {
	chi.Router

	Store  store.Store
	Config *Config
}

// New ...
func New(config *Config) *Server {
	s := &Server{
		Config: config,
		Router: chi.NewRouter(),
	}

	s.configureRoutes()
	log.Print("routes configured succesfully")

	if err := s.configureStore(); err != nil {
		log.Fatal(err)
	} else {
		log.Print("store configured succesfully")
	}

	return s
}

// NewTestServer ...
func NewTestServer(config *Config) *Server {
	s := &Server{
		Config: config,
		Router: chi.NewMux(),
	}
	if err := s.configureStore(); err != nil {
		log.Fatal(err)
	}
	return s
}

// configureRoutes ...
func (s *Server) configureRoutes() {
	s.Post("/", s.handleURLCreate)
	s.Get("/{id}", s.handleURLGet)
	s.Post("/api/shorten", s.handleURLCreateJSON())
}

// configureStore ...
func (s *Server) configureStore() error {
	switch s.Config.StorageType {
	case InMemoryStorage:
		s.Store = inmemory.New()

	case FileBasedStorage:
		s.Store = filebased.New()

	default:
		return ErrIncorrectStoreType
	}
	return nil
}

// ListenAndServe ...
func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.Config.BindAddr, s.Router)
}

// HandleErrorOr400 return true and handle error if err is not nil
func (s *Server) handleErrorOrStatus(w http.ResponseWriter, err error, status int) bool {
	if err != nil {
		http.Error(w, err.Error(), status)
	}
	return err != nil
}
