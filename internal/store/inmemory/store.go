package inmemory

import (
	"context"
	"fmt"
	"sync"

	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

type Store struct {
	mu     sync.Mutex
	closed bool

	urls map[string]*model.URL
}

// New ...
func New() *Store {
	return &Store{
		urls:   make(map[string]*model.URL),
		closed: false,
	}
}

// GetByID returns URL object and error by URL ID
func (s *Store) GetByID(_ context.Context, id string) (u *model.URL, err error) {
	s.mu.Lock()
	u, ok := s.urls[id]
	s.mu.Unlock()

	switch {
	case !ok:
		return nil, store.ErrNotFound
	case u.IsDeleted:
		return nil, store.ErrIsDeleted
	default:
		return
	}
}

// Create URL model to storage
func (s *Store) Create(ctx context.Context, u *model.URL) (err error) {
	if err = u.Validate(); err != nil {
		return fmt.Errorf("validate url: %w", err)
	}

	s.mu.Lock()
	for _, ok := s.urls[u.ID]; ok; _, ok = s.urls[u.ID] {
		if err = ctx.Err(); err != nil {
			return fmt.Errorf("context err: %w", err)
		}

		if err = u.ShortURL(); err != nil {
			return fmt.Errorf("short url: %w", err)
		}
	}

	s.urls[u.ID] = u
	s.mu.Unlock()
	return
}

// GetAllUserURLs ...
func (s *Store) GetAllUserURLs(_ context.Context, user string) (urls []*model.URL, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, u := range s.urls {
		if u.User == user {
			urls = append(urls, u)
		}
	}
	return
}

// URLsBulkCreate ...
func (s *Store) URLsBulkCreate(ctx context.Context, urls []*model.URL) (res []*model.BatchCreateURLsResponse, err error) {
	for _, u := range urls {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("context err: %w", err)
		}

		if err := u.ShortURL(); err != nil {
			return nil, fmt.Errorf("short url: %w", err)
		}

		s.mu.Lock()
		if _, ok := s.urls[u.ID]; ok {
			return nil, fmt.Errorf("already exist: %w", store.ErrAlreadyExists)
		}
		s.urls[u.ID] = u
		s.mu.Unlock()

		res = append(res, &model.BatchCreateURLsResponse{
			ShortURL:      u.ID,
			CorrelationID: u.CorelID,
		})
	}
	return res, nil
}

// URLsBulkDelete ...
func (s *Store) URLsBulkDelete(urls []string, user string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, u := range urls {
		if url := s.urls[u]; url.User == user {
			url.IsDeleted = true
		}
	}
	return nil
}

// Ping returns always
func (s *Store) Ping(_ context.Context) error {
	s.mu.Lock()
	ok := s.closed
	s.mu.Unlock()
	if ok {
		return store.ErrAlreadyClosed
	}
	return nil
}

// Close ...
func (s *Store) Close() error {
	s.mu.Lock()

	ok := s.closed
	s.closed = true

	s.mu.Unlock()

	if ok {
		return store.ErrAlreadyClosed
	}
	return nil
}
