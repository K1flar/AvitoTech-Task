package services

import (
	"banner_service/internal/domains"
	"banner_service/internal/services/bannerservice"
	"banner_service/pkg/filters"
	"context"
	"log/slog"
)

type Repository interface {
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

type Service struct {
	bannerservice.BannerService
}

func New(log *slog.Logger, repo Repository, cache Cache) *Service {
	return &Service{
		bannerservice.New(log, repo, cache),
	}
}
