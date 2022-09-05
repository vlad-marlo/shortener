package sqlstore

import (
	"context"
	"database/sql"
	"log"

	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"

	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

type SQLStore struct {
	dbURL string
}

// New ...
func New(ctx context.Context, connectString string) (*SQLStore, error) {
	s := &SQLStore{dbURL: connectString}
	db, err := pgx.Connect(ctx, connectString)
	if err != nil {
		return nil, err
	}
	defer s.closeDB(ctx, db)

	if err := s.migrate(ctx); err != nil {
		log.Print(err)
		return nil, err
	}

	log.Print("successfully created migrations")
	return s, nil
}

// migrate ...
func (s *SQLStore) migrate(ctx context.Context) error {
	db, err := pgx.Connect(ctx, s.dbURL)
	defer s.closeDB(ctx, db)

	if err != nil {
		return err
	}
	_, err = db.Exec(
		ctx,
		`CREATE TABLE IF NOT EXISTS urls(
			id SERIAL PRIMARY KEY NOT NULL,
			short VARCHAR UNIQUE,
			original_url VARCHAR,
			created_by VARCHAR
		);`,
	)
	return err
}

// closeDB ...
func (s *SQLStore) closeDB(ctx context.Context, conn *pgx.Conn) {
	if err := conn.Close(ctx); err != nil {
		log.Fatal(err)
	}
}

// Create ...
func (s *SQLStore) Create(ctx context.Context, u *model.URL) error {
	db, err := pgx.Connect(ctx, s.dbURL)
	defer s.closeDB(ctx, db)

	if err != nil {
		return err
	}
	_, err = db.Exec(
		ctx,
		`INSERT INTO urls(short, original_url, created_by) VALUES ($1, $2, $3)`,
		u.ID,
		u.BaseURL,
		u.User,
	)
	return err
}

// GetByID ...
func (s *SQLStore) GetByID(ctx context.Context, id string) (*model.URL, error) {
	db, err := pgx.Connect(ctx, s.dbURL)
	defer s.closeDB(ctx, db)
	if err != nil {
		return nil, err
	}

	u := &model.URL{}
	if err := db.QueryRow(
		ctx,
		`SELECT short, original_url, created_by FROM urls WHERE short=$1`,
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
	db, err := pgx.Connect(ctx, s.dbURL)
	defer s.closeDB(ctx, db)
	if err != nil {
		return nil, err
	}

	urls := []*model.URL{}

	r, err := db.Query(
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

	db, err := sql.Open("postgres", s.dbURL)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = db.Close(); err != nil {
			log.Print(err)
		}
	}()

	// start transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	// rollback if somethink went wrong
	defer func() {
		if err = tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Fatalf("update drivers: unable to rollback: %v", err)
		}
	}()

	stmt, err := tx.PrepareContext(
		ctx,
		`INSERT INTO urls(short, original_url, created_by) VALUES ($1, $2, $3)`,
	)

	defer func() {
		if err := stmt.Close(); err != nil && err != sql.ErrTxDone {
			log.Fatalf("update drivers: unable to close stmt: %v", err)
		}
	}()

	for _, v := range urls {
		if _, err := stmt.ExecContext(ctx, v.ID, v.BaseURL, v.User); err != nil {
			return nil, err
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
		log.Fatalf("update drivers: unable to commit: %v", err)
		return nil, err
	}

	return response, err
}

// Ping ...
func (s *SQLStore) Ping(ctx context.Context) error {
	db, _ := pgx.Connect(ctx, s.dbURL)
	defer s.closeDB(ctx, db)
	return db.Ping(ctx)
}
