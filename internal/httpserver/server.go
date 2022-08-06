package httpserver

import (
	"log"
	"net/http"

	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/inmemory"
)

type Server struct {
	srv http.Server

	Store  store.Store
	Config *Config
}

func New(config *Config) *Server {
	s := &Server{
		Config: config,
		srv:    http.Server{Addr: config.BindAddr},
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

func (s *Server) configureRoutes() {
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

func (s *Server) ListenAndServe() error {
	return s.srv.ListenAndServe()
}

// return true if err is not nil
func (s *Server) HandleErrorOr500(w http.ResponseWriter, err error) bool {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return err != nil
}
