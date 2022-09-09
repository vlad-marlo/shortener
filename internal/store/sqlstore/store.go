package sqlstore

import (
	"context"
	"database/sql"
	"log"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"

	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

type SQLStore struct {
	DB *sql.DB
}

// New ...
func New(ctx context.Context, connectString string) (*SQLStore, error) {
	db, err := sql.Open("postgres", connectString)
	if err != nil {
		return nil, err
	}
	s := &SQLStore{DB: db}

	if err := s.migrate(ctx); err != nil {
		log.Print(err)
		return nil, err
	}

	log.Print("successfully created migrations")
	return s, nil
}

// migrate run migrations instead of go-migrate package
func (s *SQLStore) migrate(ctx context.Context) error {
	_, err := s.DB.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS urls(
			id SERIAL UNIQUE PRIMARY KEY NOT NULL,
			short VARCHAR,
			original_url VARCHAR UNIQUE,
			created_by VARCHAR
		);`,
	)
	return err
}

// Create ...
func (s *SQLStore) Create(ctx context.Context, u *model.URL) error {
	_, err := s.DB.ExecContext(
		ctx,
		`INSERT INTO urls(short, original_url, created_by) VALUES ($1, $2, $3);`,
		u.ID,
		u.BaseURL,
		u.User,
	)

	log.Print(err)
	if err != nil {
		pgErr := err.(*pq.Error)
		if pgErr.Code == pgerrcode.UniqueViolation {
			if err := s.GetByOriginalURL(ctx, u); err != nil {
				return err
			}
			return store.ErrAlreadyExists
		}
		return err
	}
	return nil
}

// GetByOriginalURL ...
func (s *SQLStore) GetByOriginalURL(ctx context.Context, u *model.URL) error {
	if err := s.DB.QueryRowContext(
		ctx,
		`SELECT short FROM urls WHERE original_url = $1;`,
		u.BaseURL,
	).Scan(&u.ID); err != nil {
		return err
	}
	return nil
}

// GetByID ...
func (s *SQLStore) GetByID(ctx context.Context, id string) (*model.URL, error) {

	u := &model.URL{}
	if err := s.DB.QueryRowContext(
		ctx,
		`SELECT short, original_url, created_by FROM urls WHERE short=$1;`,
		id,
	).Scan(
		&u.ID,
		&u.BaseURL,
		&u.User,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, store.ErrNotFound
		}
		return nil, err
	}
	return u, nil
}

// GetAllUserURLs ...
func (s *SQLStore) GetAllUserURLs(ctx context.Context, userID string) ([]*model.URL, error) {
	urls := []*model.URL{}

	r, err := s.DB.QueryContext(
		ctx,
		`SELECT short, original_url, created_by FROM urls WHERE created_by=$1`,
		userID,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return urls, nil
		}
		return nil, err
	}
	if err := r.Err(); err != nil {
		return nil, err
	}
	defer r.Close()

	for r.Next() {
		u := new(model.URL)
		if err := r.Scan(&u.ID, &u.BaseURL, &u.User); err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}

	return urls, nil
}

// URLsBulkCreate ...
func (s *SQLStore) URLsBulkCreate(ctx context.Context, urls []*model.URL) ([]*model.BatchCreateURLsResponse, error) {
	if len(urls) == 0 {
		return nil, store.ErrNoContent
	}

	response := []*model.BatchCreateURLsResponse{}

	// start transaction
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}

	// rollback if somethink went wrong
	defer func() {
		if err = tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("update drivers: unable to rollback: %v", err)
		}
	}()

	stmt, err := tx.PrepareContext(
		ctx,
		`INSERT INTO urls(short, original_url, created_by) VALUES ($1, $2, $3);`,
	)

	defer func() {
		if err := stmt.Close(); err != nil && err != sql.ErrTxDone {
			log.Printf("update drivers: unable to close stmt: %v", err)
		}
	}()

	for _, v := range urls {
		if _, err := stmt.ExecContext(ctx, v.ID, v.BaseURL, v.User); err != nil {
			pgERR := err.(*pq.Error)
			if pgERR.Code != pgerrcode.UniqueViolation {
				return nil, err
			}
			if err := tx.QueryRowContext(
				ctx,
				`SELECT short FROM urls WHERE original_url = $1`,
				v.BaseURL,
			).Scan(&v.ID); err != nil {
				return nil, err
			}
		}
		response = append(
			response,
			&model.BatchCreateURLsResponse{
				ShortURL:      v.ID,
				CorrelationID: v.CorelID,
			},
		)
	}
	if err := tx.Commit(); err != nil {
		log.Printf("update drivers: unable to commit: %v", err)
		return nil, err
	}

	return response, err
}

// Ping ...
func (s *SQLStore) Ping(ctx context.Context) error {
	return s.DB.PingContext(ctx)
}

// Close ...
func (s *SQLStore) Close(ctx context.Context) error {
	return s.DB.Close()
}
