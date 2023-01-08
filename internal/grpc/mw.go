package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type mw grpc.UnaryServerInterceptor

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md := metadata.New(map[string]string{
		"user": "",
	})
	ctx = metadata.NewIncomingContext(ctx, md)
	resp, err = handler(ctx, req)
	return resp, err
}
