package grpc

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/vlad-marlo/shortener/pkg/encryptor"
)

const (
	// UserIDMDKey ...
	UserIDMDKey = "user_id"
)

// AuthInterceptor ...
func (s *Server) AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if r, ok := req.(UserGetter); ok {
			if _, err := s.getUser(r); err != nil {
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

func (s *Server) GZIPInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}
		if strings.Contains(strings.Join(md.Get("content-encoding"), " "), "gzip") {
			r, ok := req.(io.Reader)
			if !ok {
				return handler(ctx, req)
			}
			req, err := gzip.NewReader(r)
			if err != nil {
				return handler(ctx, req)
			}
			defer func() {
				if err = req.Close(); err != nil {
					s.logger.Error(fmt.Sprintf("gzip: reader: close: %v", err))
				}
			}()
		}

		return nil, nil
	}
}
