package httpserver

import (
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	chimiddlewares "github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"

	"github.com/vlad-marlo/shortener/internal/httpserver/middleware"
	"github.com/vlad-marlo/shortener/internal/poll"
	"github.com/vlad-marlo/shortener/internal/store"
)

type Server struct {
	chi.Router

	store  store.Store
	config *Config
	poller *poll.Poll
	logger *logrus.Entry
}

// New return new configured server with params from config object
// need for creating only one connection to db
func New(config *Config, storage store.Store, l *logrus.Entry) *Server {
	s := &Server{
		config: config,
		Router: chi.NewRouter(),
		logger: l,
		store:  storage,
		poller: poll.New(storage),
	}
	s.configureMiddlewares()
	l.Info("middleware configured successfully")

	s.configureRoutes()
	l.Info("routes configured successfully")

	l.Info("store configured successfully")

	return s
}

func (s *Server) Close() {
	s.poller.Close()
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

// configurePoller creates new poller
func (s *Server) configurePoller() {
	s.poller = poll.New(s.store)
}

// ListenAndServe is starting http server on correct address
func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.config.BindAddr, s.Router)
}
