package inmemory

import (
	"sync"

	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

type Store struct {
	mu sync.Mutex

	urls map[string]*model.URL
	test bool
}

// New ...
func New() *Store {
	return &Store{
		urls: make(map[string]*model.URL),
		test: false,
	}
}

// GetByID returns URL object and error by URL ID
func (s *Store) GetByID(id string) (u *model.URL, err error) {
	if !s.test {
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
func (s *Store) Create(u *model.URL) (err error) {
	if !s.test {
		defer s.mu.Unlock()
		s.mu.Lock()
	}

	if err = u.Validate(); err != nil {
		return
	}
	if _, ok := s.urls[u.ID]; ok {
		err = store.ErrAlreadyExists
		return
	}
	s.urls[u.ID] = u
	return
}
