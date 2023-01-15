package grpc

import (
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/vlad-marlo/shortener/internal/config"
	"github.com/vlad-marlo/shortener/internal/poll"
	"github.com/vlad-marlo/shortener/internal/store"
	pb "github.com/vlad-marlo/shortener/pkg/proto"
)

// Server is grpc Server
type Server struct {
	pb.UnimplementedShortenerServer
	store    store.Store
	server   *grpc.Server
	poller   *poll.Poll
	listener net.Listener

	logger *zap.Logger
}

func New(s store.Store, l *zap.Logger) (*Server, error) {
	server := &Server{
		UnimplementedShortenerServer: pb.UnimplementedShortenerServer{},
		store:                        s,
		poller:                       poll.New(s, l),
		logger:                       l,
	}
	listener, err := net.Listen("tcp", config.Get().GRPCAddr)
	if err != nil {
		return nil, err
	}
	server.listener = listener
	grpcServer := grpc.NewServer()
	pb.RegisterShortenerServer(grpcServer, server)
	server.server = grpcServer

	return server, nil
}

func (s *Server) Start() error {
	return s.server.Serve(s.listener)
}
