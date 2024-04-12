package bannerhandler

import (
	"banner_service/internal/domains"
	"banner_service/pkg/filters"
	"context"
	"log/slog"
	"net/http"
)

type BannerService interface {
	CreateBanner(ctx context.Context, banner *domains.BannerWithTagIDs) (uint32, error)
	GetBannerByFeatureAndTagID(ctx context.Context, featureID, tagID uint32, useLastRevision, isAdmin bool) (*domains.Banner, error)
	GetBanners(ctx context.Context, filter *filters.BannerFilter) ([]*domains.BannerWithTagIDs, error)
	UpdateBannerByID(ctx context.Context, id uint32, banner *domains.BannerWithTagIDs) error
	DeleteBannerByID(ctx context.Context, id uint32) error
}

type BannerHandler interface {
	PostBanner(w http.ResponseWriter, r *http.Request)
	GetBanner(w http.ResponseWriter, r *http.Request)
	GetUserBanner(w http.ResponseWriter, r *http.Request)
	PatchBannerId(w http.ResponseWriter, r *http.Request)
	DeleteBannerId(w http.ResponseWriter, r *http.Request)
}

type bannerHandler struct {
	service BannerService
	log     *slog.Logger
}

func New(service BannerService, log *slog.Logger) BannerHandler {
	return &bannerHandler{service, log}
}
