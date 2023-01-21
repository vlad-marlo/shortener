package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/vlad-marlo/shortener/internal/config"
	"github.com/vlad-marlo/shortener/internal/poll"
	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

// Service ...
type Service struct {
	logger *zap.Logger
	poller *poll.Poll
	store  store.Store
	config *config.Config
}

func New(logger *zap.Logger, store store.Store) *Service {
	return &Service{
		logger: logger,
		poller: poll.New(store, logger),
		store:  store,
		config: config.Get(),
	}
}

// CreateURL ...
func (s *Service) CreateURL(ctx context.Context, user, url string) (*model.URL, error) {
	u, err := s.NewURL(url, user)
	if err != nil {
		return nil, fmt.Errorf("model: new url: %w", err)
	}
	if err = s.store.Create(ctx, u); err != nil {
		return nil, fmt.Errorf("store: create url: %w", err)
	}
	return u, nil
}

// DeleteManyURLs ...
func (s *Service) DeleteManyURLs(user string, urls []string) {
	s.poller.DeleteURLs(urls, user)
}

// GetAllURLsByUser ...
func (s *Service) GetAllURLsByUser(ctx context.Context, user string) ([]*model.AllUserURLsResponse, error) {
	urls, err := s.store.GetAllUserURLs(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("store: get all user urls: %w", err)
	}
	var responseURLs []*model.AllUserURLsResponse
	for _, u := range urls {
		if err = ctx.Err(); err != nil {
			return nil, err
		}
		resp := &model.AllUserURLsResponse{
			ShortURL:    fmt.Sprintf("%s/%s", s.config.BaseURL, u.ID),
			OriginalURL: u.BaseURL,
		}
		responseURLs = append(responseURLs, resp)
	}
	return responseURLs, nil
}

// Ping ...
func (s *Service) Ping(ctx context.Context) error {
	if err := s.store.Ping(ctx); err != nil {
		return fmt.Errorf("store: ping: %w", err)
	}
	return nil
}

// CreateManyURLs ...
func (s *Service) CreateManyURLs(ctx context.Context, user string, urls []model.URLer) ([]*model.BatchCreateURLsResponse, error) {
	var u []*model.URL
	for _, i := range urls {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("context: %w", err)
		}
		url, err := s.NewURL(i.GetOriginalUrl(), user, i.GetCorrelationId())
		if err != nil {
			return nil, fmt.Errorf("model: url: %w", err)
		}
		u = append(u, url)
	}

	resp, err := s.store.URLsBulkCreate(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("store: urls: bulk create: %w", err)
	}

	for _, v := range resp {
		id := v.ShortURL
		v.ShortURL = fmt.Sprintf("%s/%s", s.config.BaseURL, id)
	}
	return resp, nil
}

// NewURL ...
func (s *Service) NewURL(url, user string, correlationID ...string) (*model.URL, error) {
	return model.NewURL(url, user, correlationID...)
}

// GetByID ...
func (s *Service) GetByID(ctx context.Context, id string) (*model.URL, error) {
	return s.store.GetByID(ctx, id)
}

// GetInternalStats ...
func (s *Service) GetInternalStats(ctx context.Context, _ string) (*model.InternalStat, error) {
	// TODO add checks that ip is in trusted ip subnetwork.
	return s.store.GetData(ctx)
}

// Close ...
func (s *Service) Close() error {
	s.poller.Close()
	return s.store.Close()
}
