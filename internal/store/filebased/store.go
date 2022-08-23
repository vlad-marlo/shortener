package filebased

import (
	"sync"

	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/inmemory"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

type (
	URLStore struct {
		URLs map[string]*model.URL `json:"urls"`
	}
	Store struct {
		mu       sync.Mutex
		filename string
	}
)

func New(filename string) store.Store {
	if _, err := newProducer(filename); err != nil {
		return inmemory.New()
	}
	if _, err := newCustomer(filename); err != nil {
		return inmemory.New()
	}
	return &Store{
		filename: filename,
	}
}

// Create ...
func (s *Store) Create(u *model.URL) (err error) {
	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
	}()
	return
}

// GetByID ...
func (s *Store) GetByID(id string) (u *model.URL, err error) {
	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
	}()
	return
}
