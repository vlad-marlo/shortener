package httpserver

import (
	"github.com/go-chi/chi/v5"
	chimiddlewares "github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"github.com/vlad-marlo/shortener/internal/httpserver/middleware"
	"github.com/vlad-marlo/shortener/internal/poll"
	"github.com/vlad-marlo/shortener/internal/store"
	"net/http"
	"net/http/pprof"
)

type Server struct {
	chi.Router

	store  store.Store
	config *Config
	poller *poll.Poll
	logger *logrus.Logger
}

// New return new configured server with params from config object
// need for creating only one connection to db
func New(config *Config, storage store.Store, l *logrus.Logger) (*Server, error) {
	s := &Server{
		config: config,
		Router: chi.NewRouter(),
		logger: l,
		store:  storage,
	}
	s.configureMiddlewares()
	l.Info("middleware configured successfully")

	s.configureRoutes()
	l.Info("routes configured successfully")

	s.configurePoller()
	defer s.poller.Close()

	l.Info("store configured successfully")

	return s, nil
}

// configureRoutes ...
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

		r.Route("/user/urls", func(rc chi.Router) {
			rc.Get("/", s.handleGetUserURLs)
			rc.Delete("/", s.handleURLBulkDelete)
		})
	})
}

// configureMiddlewares ...
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

// configurePoller ...
func (s *Server) configurePoller() {
	s.poller = poll.New(s.store)
}

// ListenAndServe ...
func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.config.BindAddr, s.Router)
}
