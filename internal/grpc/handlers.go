package grpc

import (
	"context"
	"errors"
	"github.com/vlad-marlo/shortener/internal/store"
	pb "github.com/vlad-marlo/shortener/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/http"
)

func (s *server) Ping(ctx context.Context, _ *pb.PingRequest) (*pb.PingResponse, error) {
	resp := new(pb.PingResponse)
	resp.Status = http.StatusOK
	if err := s.store.Ping(ctx); err != nil {
		resp.Status = http.StatusInternalServerError
	}
	return resp, nil
}

func (s *server) GetLink(ctx context.Context, r *pb.GetLinkRequest) (*pb.GetLinkResponse, error) {
	resp := new(pb.GetLinkResponse)
	url, err := s.store.GetByID(ctx, r.Id)
	if errors.Is(err, store.ErrIsDeleted) {
		resp.Status = http.StatusGone
	} else if errors.Is(err, store.ErrNotFound) {
		resp.Status = http.StatusNotFound
	} else if err != nil {
		return nil, status.Errorf(codes.Unknown, "got unexpected error: %s", err.Error())
	}
	resp.Location = url.BaseURL
	return resp, nil
}

func (s *server) CreateLinkJSON(ctx context.Context, r *pb.CreateLinkJSONRequest) (*pb.CreateLinkJSONResponse, error) {
	resp := new(pb.CreateLinkJSONResponse)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md.Set("user")
	}
	return nil, nil
}
