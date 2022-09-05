package inmemory

import (
	"context"
	"log"
	"sync"

	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

type Store struct {
	mu sync.Mutex

	urls            map[string]*model.URL
	useMutexLocking bool
}

// New ...
func New() *Store {
	log.Print("successfully configured inmemory storage")
	return &Store{
		urls:            make(map[string]*model.URL),
		useMutexLocking: true,
	}
}

// GetByID returns URL object and error by URL ID
func (s *Store) GetByID(_ context.Context, id string) (u *model.URL, err error) {
	if s.useMutexLocking {
		defer s.mu.Unlock()
		s.mu.Lock()
	}

	err = nil
	u, ok := s.urls[id]
	if !ok {
		err = store.ErrAlreadyExists
	}
	return
}

// Create URL model to storage
func (s *Store) Create(_ context.Context, u *model.URL) (err error) {
	if s.useMutexLocking {
		defer s.mu.Unlock()
		s.mu.Lock()
	}

	if err = u.Validate(); err != nil {
		return
	}
	for _, ok := s.urls[u.ID]; ok; _, ok = s.urls[u.ID] {
		u.ShortURL()
	}
	s.urls[u.ID] = u
	return
}

func (s *Store) GetAllUserURLs(_ context.Context, user string) (urls []*model.URL, err error) {
	if s.useMutexLocking {
		defer s.mu.Unlock()
		s.mu.Lock()
	}

	for _, u := range s.urls {
		if u.User == user {
			urls = append(urls, u)
		}
	}

	return
}

func (s *Store) URLsBulkCreate(_ context.Context, _ []*model.URL) ([]*model.BatchCreateURLsResponse, error) {
	return nil, nil
}

func (s *Store) Ping(_ context.Context) error {
	return nil
}
