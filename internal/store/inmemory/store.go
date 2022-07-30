package inmemory

import (
	"errors"
	"fmt"

	"github.com/vlad-marlo/shortener/internal/store/model"
)

type Store struct {
	Urls []model.URL
}

// Returns BaseURL or URL object by ID
func (s *Store) GetUrlById(id string) (string, bool) {
	for _, u := range s.Urls {
		if u.ID.String() == id {
			return u.BaseURL, true
		}
	}
	return "", false
}

func (s *Store) Create(u model.URL) error {
	if _, ok := s.GetUrlById(u.ID.String()); ok {
		return fmt.Errorf("URL with id %s already exists", u.ID)
	}
	s.Urls = append(s.Urls, u)
	return nil
}

func (s *Store) Delete(u model.URL) (bool, error) {
	if _, ok := s.GetUrlById(u.ID.String()); !ok {
		return false, errors.New("")
	}
	return true, nil
}
