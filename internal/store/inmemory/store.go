package inmemory

import (
	"fmt"

	"github.com/vlad-marlo/shortener/internal/store/model"
)

type Store struct {
	urls []model.URL
}

func New() *Store {
	return &Store{}
}

// Returns BaseURL or URL object by ID
func (s *Store) GetById(id string) (string, bool) {
	for _, u := range s.urls {
		if u.ID.String() == id {
			return u.BaseURL, true
		}
	}
	return "", false
}

// Create Url model to storage
func (s *Store) Create(u model.URL) error {
	if _, ok := s.GetById(u.ID.String()); ok {
		return fmt.Errorf("URL with id %s already exists", u.ID)
	}
	s.urls = append(s.urls, u)
	return nil
}
