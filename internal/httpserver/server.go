package httpserver

import (
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddlewares "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"

	"github.com/vlad-marlo/shortener/internal/config"
	"github.com/vlad-marlo/shortener/internal/httpserver/middleware"
	"github.com/vlad-marlo/shortener/internal/poll"
	"github.com/vlad-marlo/shortener/internal/store"
)

// vars
var (
	// timeOut is server timeout
	timeOut = 10 * time.Minute
)

// Server ...
type Server struct {
	chi.Router

	srv    *http.Server
	config *config.Config
	dev    bool

	store  store.Store
	poller *poll.Poll
	logger *zap.Logger
}

// New return new configured server with params from config object
// need for creating only one connection to db
func New(storage store.Store, l *zap.Logger) *Server {
	s := &Server{
		dev:    true,
		Router: chi.NewRouter(),
		logger: l,
		store:  storage,
		config: config.Get(),
		poller: poll.New(storage, l),
	}

	s.srv = &http.Server{
		Addr:         config.Get().BindAddr,
		Handler:      s.Router,
		ReadTimeout:  timeOut,
		WriteTimeout: timeOut,
		IdleTimeout:  timeOut,
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
	if s.dev {
		return s.srv.Close()
	}
	return nil
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
	if config.Get().HTTPS {
		manager := &autocert.Manager{
			Cache:      autocert.DirCache("cache-dir"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(config.Get().BindAddr),
		}
		s.srv.TLSConfig = manager.TLSConfig()

		return s.srv.ListenAndServeTLS("", "")
	}
	return s.srv.ListenAndServe()
}
