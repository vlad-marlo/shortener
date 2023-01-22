package grpc

import (
	"context"
	"net"

	grpc_mw "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"

	"github.com/vlad-marlo/shortener/internal/config"
	"github.com/vlad-marlo/shortener/internal/store/model"
	pb "github.com/vlad-marlo/shortener/pkg/proto"
)

type service interface {
	Ping(ctx context.Context) error
	CreateURL(ctx context.Context, user, url string) (*model.URL, error)
	DeleteManyURLs(user string, urls []string)
	GetAllURLsByUser(ctx context.Context, user string) ([]*model.AllUserURLsResponse, error)
	NewURL(url, user string, correlationID ...string) (*model.URL, error)
	CreateManyURLs(ctx context.Context, user string, urls []model.URLer) ([]*model.BatchCreateURLsResponse, error)
	GetByID(ctx context.Context, id string) (*model.URL, error)
	GetInternalStats(ctx context.Context, ip string) (*model.InternalStat, error)
}

// Server is grpc Server
type Server struct {
	pb.UnimplementedShortenerServer
	srv    service
	server *grpc.Server

	listener net.Listener

	logger *zap.Logger
}

// New ...
func New(srv service, l *zap.Logger) (*Server, error) {
	server := &Server{
		UnimplementedShortenerServer: pb.UnimplementedShortenerServer{},
		srv:                          srv,
		logger:                       l,
	}
	listener, err := net.Listen("tcp", config.Get().GRPCAddr)
	if err != nil {
		return nil, err
	}
	server.listener = listener
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_mw.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(l),
			server.CheckAuthInterceptor(),
		)),
	)
	pb.RegisterShortenerServer(grpcServer, server)
	server.server = grpcServer

	return server, nil
}

// Start ...
func (s *Server) Start() error {
	return s.server.Serve(s.listener)
}

// Close ...
func (s *Server) Close() error {
	return s.listener.Close()
}
