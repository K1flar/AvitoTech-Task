package handlers

import (
	"banner_service/internal/domains"
	"banner_service/internal/handlers/bannerhandler"
	"banner_service/pkg/filters"
	"context"
	"log/slog"
)

type Service interface {
	CreateBanner(ctx context.Context, banner *domains.BannerWithTagIDs) (uint32, error)
	GetBannerByFeatureAndTagID(ctx context.Context, featureID, tagID uint32, useLastRevision, isAdmin bool) (*domains.Banner, error)
	GetBanners(ctx context.Context, filter *filters.BannerFilter) ([]*domains.BannerWithTagIDs, error)
	UpdateBannerByID(ctx context.Context, id uint32, banner *domains.BannerWithTagIDs) error
	DeleteBannerByID(ctx context.Context, id uint32) error
}

type Handler struct {
	bannerhandler.BannerHandler
}

func New(service Service, log *slog.Logger) *Handler {
	return &Handler{
		bannerhandler.New(service, log),
	}
}
