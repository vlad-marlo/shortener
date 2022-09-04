package sqlstore

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"

	"github.com/vlad-marlo/shortener/internal/store/model"
)

type SQLStore struct {
	db *pgx.Conn
}

func New(ctx context.Context, connectString string) (*SQLStore, error) {
	db, err := pgx.Connect(ctx, connectString)
	defer db.Close(ctx)
	if err != nil {
		return nil, err
	}

	store := &SQLStore{db: db}

	if err := store.migrate(ctx); err != nil {
		log.Print(err)
		return nil, err
	}

	log.Print("successfully created migrations")
	return store, nil
}

func (s *SQLStore) migrate(ctx context.Context) error {
	_, err := s.db.Exec(
		ctx,
		`CREATE TABLE IF NOT EXISTS urls(
			id SERIAL PRIMARY KEY NOT NULL,
			short CHAR(255) UNIQUE,
			original_url CHAR(512),
			created_by CHAR(255)
		);`,
	)
	return err
}

func (s *SQLStore) Create(ctx context.Context, u *model.URL) error {

	return nil
}

func (s *SQLStore) GetByID(ctx context.Context, id string) (*model.URL, error) {
	return nil, nil
}

func (s *SQLStore) GetAllUserURLs(ctx context.Context, userID string) ([]*model.URL, error) {
	return nil, nil
}

func (s *SQLStore) Ping(ctx context.Context) error {
	return s.db.Ping(ctx)
}
