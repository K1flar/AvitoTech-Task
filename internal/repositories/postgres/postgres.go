package postgres

import (
	"banner_service/internal/config"
	"database/sql"

	_ "github.com/lib/pq"
)

type Repository struct{}

func New(cfg config.Database) (*Repository, error) {
	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)

	return &Repository{}, nil
}
