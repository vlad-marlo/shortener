package sqlstore

import (
	"context"

	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4"

	"github.com/vlad-marlo/shortener/internal/store/model"
)

type SQLStore struct {
	db *pgx.Conn
}

func New(connectString string) (*SQLStore, error) {
	db, err := pgx.Connect(context.Background(), connectString)
	if err != nil {
		return nil, err
	}
	return &SQLStore{db: db}, nil
}

func (s *SQLStore) Create(u *model.URL) error {
	return nil
}

func (s *SQLStore) GetByID(id string) (*model.URL, error) {
	return nil, nil
}

func (s *SQLStore) GetAllUserURLs(userID string) ([]*model.URL, error) {
	return nil, nil
}

func (s *SQLStore) Ping() error {
	return s.db.Ping()
}
