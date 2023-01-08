package grpc

import (
	"github.com/vlad-marlo/shortener/internal/store"
	pb "github.com/vlad-marlo/shortener/pkg/proto"
)

// server is grpc server
type server struct {
	pb.UnimplementedShortenerServer
	store store.Store
}
