package filebased

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

type Store struct {
	Filename string
	closed   bool
	mu       sync.Mutex
}

func New(filename string) (*Store, error) {
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

func (s *Store) URLsBulkDelete(_ []string, _ string) error {
	return nil
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return store.ErrAlreadyClosed
	}
	s.closed = true
	return nil
}
