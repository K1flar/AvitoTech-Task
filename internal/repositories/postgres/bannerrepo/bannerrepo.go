package bannerrepo

import (
	"banner_service/internal/config"
	"database/sql"
	"fmt"
)

var (
	ErrNotFound         = fmt.Errorf("banner not found")
	ErrNotJSON          = fmt.Errorf("content is not json")
	ErrInvalidFeatureID = fmt.Errorf("feature id must be positive")
	ErrInvalidTagID     = fmt.Errorf("tag id must be positive")
)

type BannerRepository struct {
	db  *sql.DB
	cfg *config.Database
}

func New(cfg *config.Database, db *sql.DB) *BannerRepository {
	return &BannerRepository{
		db:  db,
		cfg: cfg,
	}
}
