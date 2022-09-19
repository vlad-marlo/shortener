package filebased

import (
	"context"
	"errors"
	"log"

	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

type Store struct {
	Filename string
}

func New(filename string) (store.Store, error) {
	if filename == "" {
		return nil, errors.New("store wasn't configured successfully")
	}

	s := &Store{
		Filename: filename,
	}
	p, err := newProducer(s.Filename)
	defer func() {
		if err = p.Close(); err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		return nil, err
	}
	log.Print("successfully configured file-based store")
	return s, nil
}

func (s *Store) GetByID(_ context.Context, id string) (*model.URL, error) {
	p, err := newProducer(s.Filename)
	defer func() {
		if err = p.Close(); err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		return nil, err
	}
	return p.GetURLByID(id)
}

func (s *Store) Create(_ context.Context, u *model.URL) error {
	p, err := newProducer(s.Filename)
	defer func() {
		if err = p.Close(); err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		return err
	}
	return p.CreateURL(u)
}

func (s *Store) GetAllUserURLs(_ context.Context, user string) ([]*model.URL, error) {
	p, err := newProducer(s.Filename)
	defer func() {
		if err = p.Close(); err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		return nil, err
	}
	return p.GetAllUserURLs(user)
}

func (s *Store) Ping(_ context.Context) error {
	return nil
}

func (s *Store) URLsBulkCreate(_ context.Context, _ []*model.URL) ([]*model.BatchCreateURLsResponse, error) {
	return nil, nil
}

func (s *Store) Close() error {
	return nil
}
