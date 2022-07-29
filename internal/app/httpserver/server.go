package httpserver

import "net/http"

type Server struct {
	*http.Server
}

func New() *Server {
	return &Server{}
}
