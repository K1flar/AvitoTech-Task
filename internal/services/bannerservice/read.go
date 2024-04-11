package bannerservice

import (
	"banner_service/internal/domains"
	"banner_service/internal/repositories/postgres/bannerrepo"
	"banner_service/pkg/filters"
	"context"
	"errors"
	"fmt"
)

func (s *bannerService) GetBannerByFeatureAndTagID(ctx context.Context, featureID, tagID uint32) (*domains.Banner, error) {
	fn := `bannerService.GetBannerByFeatureAndTagID`

	key := domains.BannerKey{FeatureID: featureID, TagID: tagID}

	if _, ok := s.cache.Get(key); !ok {
		banner, err := s.repo.GetBannerByFeatureAndTagID(ctx, featureID, tagID)
		if err != nil {
			s.log.Error(fmt.Sprintf("%s: %s", fn, err))
			if errors.Is(err, bannerrepo.ErrNotFound) {
				return nil, fmt.Errorf("%s: %w", fn, ErrNotFound)
			}
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		s.cache.Set(key, banner)
	}

	banner, _ := s.cache.Get(key)

	return banner, nil
}

func (s *bannerService) GetBanners(ctx context.Context, filter *filters.BannerFilter) ([]*domains.BannerWithTagIDs, error) {
	fn := `bannerService.GetBanners`

	filter.Validate()

	banners, err := s.repo.GetBanners(ctx, filter)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err))
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return banners, nil
}
