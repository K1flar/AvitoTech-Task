package postgres

import (
	"banner_service/internal/config"
	"banner_service/internal/repositories/postgres/bannerrepo"
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"
)

type Repository struct {
	bannerrepo.BannerRepository
}

func New(cfg *config.Database, log *slog.Logger) (*Repository, error) {
	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)

	return &Repository{
		bannerrepo.New(cfg, db, log),
	}, nil
}
