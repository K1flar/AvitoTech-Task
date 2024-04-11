package bannerservice

import (
	"banner_service/internal/domains"
	"banner_service/internal/repositories/postgres/bannerrepo"
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

func (s *bannerService) CreateBanner(ctx context.Context, banner *domains.BannerWithTagIDs) (uint32, error) {
	fn := `bannerService.CreateBanner`

	err := validateBannerWithTagIDs(banner)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err))
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	id, err := s.repo.CreateBanner(ctx, banner)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err))
		if errors.Is(err, bannerrepo.ErrAlreadyExists) {
			return 0, fmt.Errorf("%s: %w", fn, ErrAlreadyExists)
		}
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return id, nil
}

func validateBannerWithTagIDs(banner *domains.BannerWithTagIDs) error {
	if err := json.Unmarshal([]byte(banner.Content), &map[string]interface{}{}); err != nil {
		return ErrNotJSON
	}

	if banner.FeatureID <= 0 {
		return ErrInvalidFeatureID
	}

	for _, tagID := range banner.TagIDs {
		if tagID <= 0 {
			return ErrInvalidTagID
		}
	}

	return nil
}
