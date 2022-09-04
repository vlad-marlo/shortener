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

// New return new configured server with params from config object
func New(config *Config) *Server {
	s := &Server{
		Config: config,
		Router: chi.NewRouter(),
	}

	s.configureMiddlewares()
	s.configureRoutes()
	log.Print("routes configured successfully")

	s.configureStore()
	log.Print("store configured successfully")

	return s
}

// configureRoutes ...
func (s *Server) configureRoutes() {
	s.Post("/", s.handleURLCreate)
	s.Get("/{id}", s.handleURLGet)

	s.Post("/api/shorten", s.handleURLCreateJSON)
	s.Get("/api/user/urls", s.handleGetUserURLs)
	s.Get("/ping", s.handlePingStore)
}

// configureMiddlewares ...
func (s *Server) configureMiddlewares() {
	s.Use(middleware.GzipCompression)
	s.Use(middleware.AuthMiddleware)
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
