package bannerrepo

import (
	"banner_service/internal/config"
	"banner_service/internal/domains"
	"banner_service/pkg/filters"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

var (
	ErrNotFound      = fmt.Errorf("banner not found")
	ErrAlreadyExists = fmt.Errorf("there is already a pair that uniquely identifies the banner")
)

type BannerRepository interface {
	CreateBanner(ctx context.Context, banner *domains.BannerWithTagIDs) (uint32, error)
	GetBannerByFeatureAndTagID(ctx context.Context, featureID, tagID uint32) (*domains.Banner, error)
	GetBanners(ctx context.Context, filter *filters.BannerFilter) ([]*domains.BannerWithTagIDs, error)
	UpdateBannerByID(ctx context.Context, id uint32, banner *domains.BannerWithTagIDs) error
	DeleteBannerByID(ctx context.Context, id uint32) error
}

type bannerRepository struct {
	cfg *config.Database
	db  *sql.DB
	log *slog.Logger
}

func New(cfg *config.Database, db *sql.DB, log *slog.Logger) BannerRepository {
	return &bannerRepository{cfg, db, log}
}
