package inmemory

import (
	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

type Store struct {
	urls []model.URL
}

func New() *Store {
	return &Store{}
}

// GetByID Returns BaseURL or URL object by ID
func (s *Store) GetByID(id string) (model.URL, error) {
	for _, u := range s.urls {
		if u.ID.String() == id {
			return u, nil
		}
	}
	return model.URL{}, store.ErrNotFound
}

func (s *Store) urlExists(url model.URL) bool {
	for _, u := range s.urls {
		if u.ID.String() == url.ID.String() || u.BaseURL == url.BaseURL {
			return true
		}
	}
	return false
}

// Create Url model to storage
func (s *Store) Create(u model.URL) error {
	if s.urlExists(u) {
		return store.ErrAlreadyExists
	}
	s.urls = append(s.urls, u)
	return nil
}
