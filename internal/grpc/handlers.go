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

// UserGetter ...
type UserGetter interface {
	// GetUser ...
	GetUser() string
}

// Ping ...
func (s *Server) Ping(ctx context.Context, _ *pb.PingRequest) (*pb.PingResponse, error) {
	var resp pb.PingResponse
	resp.Status = http.StatusOK
	if err := s.store.Ping(ctx); err != nil {
		resp.Status = http.StatusInternalServerError
	}
	return &resp, nil
}

// GetLink ...
func (s *Server) GetLink(ctx context.Context, r *pb.GetLinkRequest) (*pb.GetLinkResponse, error) {
	resp := new(pb.GetLinkResponse)
	url, err := s.store.GetByID(ctx, r.Id)
	resp.Status = http.StatusOK
	if errors.Is(err, store.ErrIsDeleted) {
		resp.Status = http.StatusGone
	} else if errors.Is(err, store.ErrNotFound) {
		return nil, NotFound()
	} else if err != nil {
		resp.Status = http.StatusInternalServerError
		return nil, Internal()
	}
	resp.Location = url.BaseURL
	return resp, nil
}

// CreateLinkJSON ...
func (s *Server) CreateLinkJSON(ctx context.Context, r *pb.CreateLinkJSONRequest) (*pb.CreateLinkJSONResponse, error) {
	var resp pb.CreateLinkJSONResponse

	user, err := s.getUser(r)
	if err != nil {
		return nil, Unauthenticated()
	}

	var u *model.URL
	u, err = model.NewURL(r.Url, user)
	if err != nil {
		resp.Status = http.StatusBadRequest
		return &resp, nil
	}

	if err = s.store.Create(ctx, u); errors.Is(err, store.ErrAlreadyExists) {
		resp.Status = http.StatusConflict
	} else if err != nil {
		resp.Status = http.StatusInternalServerError
		return &resp, nil
	} else {
		resp.Status = http.StatusCreated
	}

	resp.Result = u.ID

	return &resp, nil
}

// CreateManyLinks ...
func (s *Server) CreateManyLinks(ctx context.Context, r *pb.CreateManyRequest) (*pb.CreateManyResponse, error) {
	var resp pb.CreateManyResponse

	user, err := s.getUser(r)
	if err != nil {
		return nil, Unauthenticated()
	}

	var urls []*model.URL
	for _, u := range r.Urls {
		if ctx.Err() != nil {
			return nil, status.Error(codes.Canceled, "canceled")
		}

		var url *model.URL
		url, err = model.NewURL(u.OriginalUrl, user, u.CorrelationId)
		if err != nil {
			return nil, Internal()
		}

		urls = append(urls, url)
	}

	if len(r.Urls) == 0 {
		return nil, BadRequest()
	}

	var res []*model.BatchCreateURLsResponse
	res, err = s.store.URLsBulkCreate(ctx, urls)
	for _, b := range res {
		if ctx.Err() != nil {
			return nil, status.Error(codes.Canceled, "canceled")
		}
		resp.Urls = append(resp.Urls, &pb.CreateManyResponse_URL{
			CorrelationId: b.CorrelationID,
			ShortUrl:      b.ShortURL,
		})
	}

	return &resp, nil
}

// CreateLink ...
func (s *Server) CreateLink(ctx context.Context, r *pb.CreateLinkRequest) (*pb.CreateLinkResponse, error) {
	var resp pb.CreateLinkResponse

	user, err := s.getUser(r)
	if err != nil {
		return nil, Unauthenticated()
	}

	var u *model.URL
	u, err = model.NewURL(r.Url, user)
	if err != nil {
		resp.Status = http.StatusBadRequest
		return &resp, status.Error(codes.InvalidArgument, "bad request data")
	}

	if err = s.store.Create(ctx, u); errors.Is(err, store.ErrAlreadyExists) {
		resp.Status = http.StatusConflict
	} else if err != nil {
		resp.Status = http.StatusInternalServerError
		return &resp, status.Errorf(codes.Internal, "unknown error: %v", err)
	} else {
		resp.Status = http.StatusCreated
	}

	resp.Result = u.ID

	return &resp, nil
}

// GetUser ...
func (s *Server) GetUser(context.Context, *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	var resp pb.GetUserResponse
	resp.User = encryptor.Get().EncodeUUID(uuid.NewString())
	return &resp, nil
}

// GetManyLinks ...
func (s *Server) GetManyLinks(ctx context.Context, r *pb.GetManyLinksRequest) (*pb.GetManyLinksResponse, error) {
	var resp pb.GetManyLinksResponse
	user, err := s.getUser(r)
	if err != nil {
		return nil, Unauthenticated()
	}
	urls, err := s.store.GetAllUserURLs(ctx, user)
	if err != nil {
		resp.Status = http.StatusInternalServerError
		return &resp, nil
	}

	if len(urls) == 0 {
		resp.Status = http.StatusNoContent
		return &resp, nil
	}

	for _, u := range urls {
		resp.Urls = append(resp.Urls, &pb.GetManyLinksResponse_URL{
			OriginalUrl: u.BaseURL,
			ShortUrl:    u.ID,
		})
	}

	return &resp, nil
}

// DeleteMany ...
func (s *Server) DeleteMany(_ context.Context, r *pb.DeleteManyRequest) (*pb.DeleteManyResponse, error) {
	var resp pb.DeleteManyResponse
	user, err := s.getUser(r)
	if err != nil {
		return nil, Unauthenticated()
	}
	s.poller.DeleteURLs(r.Ids, user)
	resp.Status = http.StatusAccepted
	return &resp, nil
}

// getUser ...
func (s *Server) getUser(r UserGetter) (res string, err error) {
	if err = encryptor.Get().DecodeUUID(r.GetUser(), &res); err != nil {
		return "", fmt.Errorf("decode uuid: %w", err)
	}
	return
}
