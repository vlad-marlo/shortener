package inmemory

import (
	"log"
	"sync"

	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

type Store struct {
	urls map[string]model.URL
	mu   sync.Mutex
	test bool
}

func New() *Store {
	return &Store{
		urls: make(map[string]model.URL),
		test: true,
	}
}

// GetByID Returns BaseURL or URL object by ID
func (s *Store) GetByID(id string) (model.URL, error) {
	if !s.test {
		s.mu.Lock()
		defer s.mu.Unlock()
	}

	if u, ok := s.urls[id]; ok {
		log.Printf("get url %v\n", u)
		return u, nil
	}
	return model.URL{}, store.ErrNotFound
}

// Create Url model to storage
func (s *Store) Create(u model.URL) error {
	if !s.test {
		s.mu.Lock()
		defer s.mu.Unlock()
	}

	if err := u.Validate(); err != nil {
		return err
	}
	if _, ok := s.urls[u.ID]; ok {
		return store.ErrAlreadyExists
	}
	s.urls[u.ID] = u
	log.Printf("created url %v\n", u)
	return nil
}
