package httpserver

import (
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	chimiddlewares "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/vlad-marlo/shortener/internal/httpserver/middleware"
	"github.com/vlad-marlo/shortener/internal/poll"
	"github.com/vlad-marlo/shortener/internal/store"
)

// Server ...
type Server struct {
	chi.Router

	store  store.Store
	config *Config
	poller *poll.Poll
	logger *zap.Logger
}

// New return new configured server with params from config object
// need for creating only one connection to db
func New(config *Config, storage store.Store, l *zap.Logger) *Server {
	s := &Server{
		config: config,
		Router: chi.NewRouter(),
		logger: l,
		store:  storage,
		poller: poll.New(storage, l),
	}
	s.configureMiddlewares()
	l.Info("middleware configured successfully")

	s.configureRoutes()
	l.Info("routes configured successfully")

	l.Info("store configured successfully")

	return s
}

// Close closes poller and storage connection.
func (s *Server) Close() error {
	s.poller.Close()
	return s.store.Close()
}

// configureRoutes initialize all endpoints of server
func (s *Server) configureRoutes() {
	s.Route("/debug", func(r chi.Router) {
		r.HandleFunc("/pprof/", pprof.Index)
		r.HandleFunc("/pprof/allocs", pprof.Index)
		r.HandleFunc("/pprof/heap", pprof.Index)
		r.HandleFunc("/pprof/mutex", pprof.Index)
		r.HandleFunc("/pprof/block", pprof.Index)
		r.HandleFunc("/pprof/threadcreate", pprof.Index)
		r.HandleFunc("/pprof/goroutine", pprof.Index)
		r.HandleFunc("/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/pprof/profile", pprof.Profile)
		r.HandleFunc("/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/pprof/trace", pprof.Trace)
	})

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

// configureMiddlewares is adding middlewares to server
func (s *Server) configureMiddlewares() {
	s.Use(
		chimiddlewares.RequestID,
		// my own middlewares
		middleware.GzipCompression,
		middleware.AuthMiddleware,

		// chi middlewares
		middleware.Logger(s.logger),
	)
}

// ListenAndServe is starting http server on correct address
func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.config.BindAddr, s.Router)
}
