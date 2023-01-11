package grpc

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/model"
	"github.com/vlad-marlo/shortener/pkg/encryptor"
	pb "github.com/vlad-marlo/shortener/pkg/proto"
)

func (s *server) Ping(ctx context.Context, _ *pb.PingRequest) (*pb.PingResponse, error) {
	var resp pb.PingResponse
	resp.Status = http.StatusOK
	if err := s.store.Ping(ctx); err != nil {
		resp.Status = http.StatusInternalServerError
	}
	return &resp, nil
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
	var resp pb.CreateLinkJSONResponse

	user, err := s.getUser(r.User)
	if err != nil {
		return &resp, status.Error(codes.Unauthenticated, "bad user")
	}

	var u *model.URL
	u, err = model.NewURL(r.Url, user)
	if err != nil {
		resp.Status = http.StatusBadRequest
		return &resp, status.Error(codes.InvalidArgument, "bad request data")
	}

	if err = s.store.Create(ctx, u); err != nil {
		// TODO: add more status codes and err check
		return &resp, status.Error(codes.Internal, "bad request")
	}

	resp.Result = u.ID

	return &resp, nil
}

func (s *server) CreateLink(ctx context.Context, r *pb.CreateLinkRequest) (*pb.CreateLinkResponse, error) {
	var resp pb.CreateLinkResponse

	user, err := s.getUser(r.User)
	if err != nil {
		return &resp, status.Error(codes.Unauthenticated, "bad user")
	}

	var u *model.URL
	u, err = model.NewURL(r.Url, user)
	if err != nil {
		resp.Status = http.StatusBadRequest
		return &resp, status.Error(codes.InvalidArgument, "bad request data")
	}

	if err = s.store.Create(ctx, u); err != nil {
		// TODO: add more status codes and err check
		return &resp, status.Error(codes.Internal, "bad request")
	}

	resp.Result = u.ID

	return &resp, nil
}

func (s *server) GetUser(context.Context, *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	var resp pb.GetUserResponse
	resp.User = encryptor.Get().EncodeUUID(uuid.NewString())
	return &resp, nil
}

func (s *server) getUser(u string) (res string, err error) {
	if err = encryptor.Get().DecodeUUID(u, &res); err != nil {
		return "", fmt.Errorf("decode uuid: %w", err)
	}
	return
}
