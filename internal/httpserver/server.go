package httpserver

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddlewares "github.com/go-chi/chi/v5/middleware"
	"github.com/vlad-marlo/shortener/internal/httpserver/middleware"
	"github.com/vlad-marlo/shortener/internal/poll"
	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/filebased"
	"github.com/vlad-marlo/shortener/internal/store/inmemory"
	"github.com/vlad-marlo/shortener/internal/store/sqlstore"
)

type Server struct {
	chi.Router

	store  store.Store
	config *Config
	poller *poll.Poll
}

// New ...
func New(config *Config) *Server {
	s := &Server{
		config: config,
		Router: chi.NewRouter(),
	}
	s.configureMiddlewares()
	log.Print("middleware configured successfully")

	s.configureRoutes()
	log.Print("routes configured successfully")

	return s
}

// Start return new configured server with params from config object
// need for creating only one connection to db
func Start(config *Config) error {
	s := New(config)

	if err := s.configureStore(); err != nil {
		return err
	}

	defer func() {
		if err := s.store.Close(); err != nil {
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

	s.Get("/ping", s.handlePingStore)

	s.Route("/api", func(r chi.Router) {
		r.Post("/shorten", s.handleURLCreateJSON)
		r.Post("/shorten/batch", s.handleURLBulkCreate)

		r.Route("/user/urls", func(r chi.Router) {
			r.Get("/", s.handleGetUserURLs)
			r.Delete("/", s.handleURLBulkDelete)
		})
	})
}

// configureMiddlewares ...
func (s *Server) configureMiddlewares() {
	s.Use(
		// my own middlewares
		middleware.GzipCompression,
		middleware.AuthMiddleware,

		// chi middlewares
		chimiddlewares.RealIP,
		chimiddlewares.RequestID,
		chimiddlewares.Logger,
	)
}

// configureStore ...
func (s *Server) configureStore() (err error) {
	switch s.config.StorageType {
	case store.InMemoryStorage:
		s.store, err = inmemory.New(), nil
	case store.FileBasedStorage:
		s.store, err = filebased.New(s.config.FilePath)
	case store.SQLStore:
		s.store, err = sqlstore.New(context.Background(), s.config.Database)
	default:
		s.store, err = filebased.New(s.config.FilePath)
	}
	return
}

// configurePoller ...
func (s *Server) configurePoller() {
	s.poller = poll.New(s.store)
}

// ListenAndServe ...
func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.config.BindAddr, s.Router)
}
