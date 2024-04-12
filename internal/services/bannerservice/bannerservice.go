package bannerservice

import (
	"banner_service/internal/domains"
	"banner_service/pkg/filters"
	"context"
	"fmt"
	"log/slog"
)

var (
	ErrNotFound         = fmt.Errorf("banner not found")
	ErrAlreadyExists    = fmt.Errorf("there is already a pair that uniquely identifies the banner")
	ErrNotJSON          = fmt.Errorf("content is not json")
	ErrInvalidFeatureID = fmt.Errorf("feature id must be positive")
	ErrInvalidTagID     = fmt.Errorf("tag id must be positive")
	ErrNoTagIDs         = fmt.Errorf("no tag ids")
)

type BannerRepository interface {
	CreateBanner(ctx context.Context, banner *domains.BannerWithTagIDs) (uint32, error)
	GetBannerByFeatureAndTagID(ctx context.Context, featureID, tagID uint32) (*domains.Banner, error)
	GetBanners(ctx context.Context, filter *filters.BannerFilter) ([]*domains.BannerWithTagIDs, error)
	UpdateBannerByID(ctx context.Context, id uint32, banner *domains.BannerWithTagIDs) error
	DeleteBannerByID(ctx context.Context, id uint32) error
}

type Cache interface {
	Get(domains.BannerKey) (*domains.Banner, bool)
	Set(domains.BannerKey, *domains.Banner)
	GetFrequency(domains.BannerKey) int
}

type BannerService interface {
	CreateBanner(ctx context.Context, banner *domains.BannerWithTagIDs) (uint32, error)
	GetBannerByFeatureAndTagID(ctx context.Context, featureID, tagID uint32, useLastRevision, isAdmin bool) (*domains.Banner, error)
	GetBanners(ctx context.Context, filter *filters.BannerFilter) ([]*domains.BannerWithTagIDs, error)
	UpdateBannerByID(ctx context.Context, id uint32, banner *domains.BannerWithTagIDs) error
	DeleteBannerByID(ctx context.Context, id uint32) error
}

type bannerService struct {
	log   *slog.Logger
	repo  BannerRepository
	cache Cache
}

func New(log *slog.Logger, repo BannerRepository, cache Cache) BannerService {
	return &bannerService{log, repo, cache}
}
