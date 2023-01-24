package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/vlad-marlo/shortener/pkg/encryptor"
)

const (
	// UserIDMDKey ...
	UserIDMDKey = "user_id"
)

// CheckAuthInterceptor checks if grpc request has user field and  authentication data contains in this field or in metadata.
func (s *Server) CheckAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if r, ok := req.(UserGetter); ok {
			if _, err := s.getUser(r); err != nil && s.UserFromCtx(ctx) == "" {
				return nil, Unauthenticated()
			}
		}
		return handler(ctx, req)
	}
}

// UserFromCtx ...
func (s *Server) UserFromCtx(ctx context.Context) (id string) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	ids := md.Get(UserIDMDKey)
	if len(ids) != 1 {
		return ""
	}
	if err := encryptor.Get().DecodeUUID(ids[0], &id); err != nil {
		return ""
	}
	return
}
