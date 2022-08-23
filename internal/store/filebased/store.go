package filebased

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/inmemory"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

type Store struct {
	filename string `env:"FILE_STORAGE_PATH"`
}

type urls struct {
	URLs map[string]*model.URL `json:"urls"`
}

func New() store.Store {
	s := &Store{}
	if err := env.Parse(s); err != nil {
		return inmemory.New()
	}
	if _, err := newProducer(s.filename); err != nil {
		return inmemory.New()
	}
	log.Print("successfully configured filebased store")
	return s
}

func (s *Store) GetByID(id string) (u *model.URL, err error) {
	p, err := newProducer(s.filename)
	defer p.Close()
	if err != nil {
		return nil, err
	}
	u, err = p.ReadURL(id)
	return
}

func (s *Store) Create(u *model.URL) (err error) {
	p, err := newProducer(s.filename)
	defer p.Close()
	if err != nil {
		return
	}
	err = p.WriteURL(u)
	return
}
