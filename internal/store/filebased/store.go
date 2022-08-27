package filebased

import (
	"log"

	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/inmemory"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

type Store struct {
	Filename string
}

func New(filename string) store.Store {
	s := &Store{
		Filename: filename,
	}
	p, err := newProducer(s.Filename)
	defer func() {
		err = p.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		log.Print(err)
		return inmemory.New()
	}
	log.Print("successfully configured file-based store")
	return s
}

func (s *Store) GetByID(id string) (u *model.URL, err error) {
	p, err := newProducer(s.Filename)
	defer func() {
		err = p.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		return nil, err
	}
	u, err = p.GetURLByID(id)
	return
}

func (s *Store) Create(u *model.URL) error {
	p, err := newProducer(s.Filename)
	defer func() {
		err = p.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		return err
	}
	return p.CreateURL(u)
}
