package store

import "github.com/vlad-marlo/shortener/internal/store/model"

type Store interface {
	Create(model.URL) error
	GetById(int) (model.URL, error)
}
