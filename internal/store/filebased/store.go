package filebased

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/inmemory"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

type Store struct {
	Filename string `env:"FILE_STORAGE_PATH"`
}

type URLs struct {
	URLs map[string]*model.URL `json:"urls"`
}

func New() store.Store {
	s := &Store{}
	if err := env.Parse(s); err != nil {
		log.Print(err)
		return inmemory.New()
	}
	log.Print(s.Filename)
	p, err := newProducer(s.Filename)
	defer p.Close()
	if err != nil {
		log.Print(err)
		return inmemory.New()
	}
	log.Print("successfully configured file-based store")
	return s
}

func (s *Store) GetByID(id string) (u *model.URL, err error) {
	p, err := newProducer(s.Filename)
	defer p.Close()
	if err != nil {
		return nil, err
	}
	u, err = p.GetURLByID(id)
	return
}

func (s *Store) Create(u *model.URL) error {
	p, err := newProducer(s.Filename)
	defer p.Close()
	if err != nil {
		return err
	}
	return p.CreateURL(u)
}
