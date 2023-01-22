package grpc

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	srv "github.com/vlad-marlo/shortener/internal/service"
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
	if err := s.srv.Ping(ctx); err != nil {
		s.logger.Error("ping", zap.Error(err))
		resp.Status = http.StatusInternalServerError
	}
	return &resp, nil
}

// GetLink ...
func (s *Server) GetLink(ctx context.Context, r *pb.GetLinkRequest) (*pb.GetLinkResponse, error) {
	resp := new(pb.GetLinkResponse)
	url, err := s.srv.GetByID(ctx, r.Id)
	resp.Status = http.StatusOK
	switch {
	case errors.Is(err, store.ErrIsDeleted):
		resp.Status = http.StatusGone
	case errors.Is(err, store.ErrNotFound):
		return nil, NotFound()
	case err != nil:
		s.logger.Error("grpc: get link", zap.Error(err))
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
	u, err = s.srv.CreateURL(ctx, user, r.Url)
	if errors.Is(err, store.ErrAlreadyExists) {
		resp.Status = http.StatusConflict
	} else if err != nil {
		return nil, Internal()
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

	if len(r.Urls) == 0 {
		return nil, BadRequest()
	}

	var urls []model.URLer
	for _, u := range r.Urls {
		urls = append(urls, u)
	}

	var res []*model.BatchCreateURLsResponse
	res, err = s.srv.CreateManyURLs(ctx, user, urls)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, Canceled()
		}
		return nil, Internal()
	}
	for _, b := range res {
		if ctx.Err() != nil {
			return nil, Canceled()
		}
		resp.Urls = append(resp.Urls, &pb.CreateManyResponse_URL{
			CorrelationId: b.CorrelationID,
			ShortUrl:      b.ShortURL,
		})
	}

	return &resp, nil
}

// CreateLink xd.
func (s *Server) CreateLink(ctx context.Context, r *pb.CreateLinkRequest) (*pb.CreateLinkResponse, error) {
	var resp pb.CreateLinkResponse

	user, _ := s.getUser(r)

	u, err := s.srv.CreateURL(ctx, user, r.Url)
	if errors.Is(err, store.ErrAlreadyExists) {
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
	user, _ := s.getUser(r)
	urls, err := s.srv.GetAllURLsByUser(ctx, user)
	if err != nil {
		return nil, Internal()
	}

	if len(urls) == 0 {
		resp.Status = http.StatusNoContent
		return &resp, nil
	}

	for _, u := range urls {
		resp.Urls = append(resp.Urls, &pb.GetManyLinksResponse_URL{
			OriginalUrl: u.OriginalURL,
			ShortUrl:    u.ShortURL,
		})
	}

	return &resp, nil
}

// DeleteMany ...
func (s *Server) DeleteMany(_ context.Context, r *pb.DeleteManyRequest) (*pb.DeleteManyResponse, error) {
	var resp pb.DeleteManyResponse
	// we can do not check error because of interceptor which is checking if request have user field than this field is valid
	// else interceptor will return unauthorized error to user.
	u, _ := s.getUser(r)
	s.srv.DeleteManyURLs(u, r.Ids)
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

func (s *Server) GetInternalStats(ctx context.Context, r *pb.GetInternalStatsRequest) (*pb.GetInternalStatsResponse, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, PermissionDenied()
	}
	stat, err := s.srv.GetInternalStats(ctx, p.Addr.String())
	if err != nil {
		if errors.Is(err, srv.ErrForbidden) {
			return nil, PermissionDenied()
		}

		return nil, Internal()
	}
	return &pb.GetInternalStatsResponse{
		Urls:  stat.CountOfURLs,
		Users: stat.CountOfUsers,
	}, nil
}
