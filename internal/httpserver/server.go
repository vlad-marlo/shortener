package httpserver

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vlad-marlo/shortener/internal/httpserver/middleware"
	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/filebased"
	"github.com/vlad-marlo/shortener/internal/store/inmemory"
	"github.com/vlad-marlo/shortener/internal/store/sqlstore"
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
	s.configureMiddlewares()
	log.Print("middleware configured successfully")
	s.configureRoutes()
	log.Print("routes configured successfully")

	return s
}

// New return new configured server with params from config object
// need for creating only one connection to db
func Start(config *Config) error {
	s := New(config)

	if err := s.configureStore(); err != nil {
		return err
	}

	defer func() {
		if err := s.Store.Close(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	log.Print("store configured successfully")

	return s.ListenAndServe()
}

// configureRoutes ...
func (s *Server) configureRoutes() {
	s.Post("/", s.handleURLCreate)
	s.Get("/{id}", s.handleURLGet)

	s.Post("/api/shorten", s.handleURLCreateJSON)
	s.Get("/api/user/urls", s.handleGetUserURLs)
	s.Get("/ping", s.handlePingStore)
	s.Post("/api/shorten/batch", s.handleURLBulkCreate)
}

// configureMiddlewares ...
func (s *Server) configureMiddlewares() {
	s.Use(middleware.GzipCompression)
	s.Use(middleware.AuthMiddleware)
	s.Use(middleware.LogResponse)
}

// configureStore ...
func (s *Server) configureStore() (err error) {
	switch s.Config.StorageType {
	case store.InMemoryStorage:
		s.Store, err = inmemory.New(), nil
	case store.FileBasedStorage:
		s.Store, err = filebased.New(s.Config.FilePath)
	case store.SQLStore:
		s.Store, err = sqlstore.New(context.Background(), s.Config.Database)
	default:
		s.Store, err = filebased.New(s.Config.FilePath)
	}
	return
}

// ListenAndServe ...
func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.Config.BindAddr, s.Router)
}
