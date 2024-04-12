package bannerservice

import (
	"banner_service/internal/domains"
	"banner_service/internal/repositories/postgres/bannerrepo"
	"context"
	"errors"
	"fmt"
)

func (s *bannerService) UpdateBannerByID(ctx context.Context, id uint32, banner *domains.BannerWithTagIDs) error {
	fn := `bannerService.UpdateBannerByID`

	err := validateBannerWithTagIDs(banner)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err))
		return fmt.Errorf("%s: %w", fn, err)
	}

	err = s.repo.UpdateBannerByID(ctx, id, banner)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err))
		switch {
		case errors.Is(err, bannerrepo.ErrNotFound):
			return fmt.Errorf("%s: %w", fn, ErrNotFound)
		case errors.Is(err, bannerrepo.ErrAlreadyExists):
			return fmt.Errorf("%s: %w", fn, ErrAlreadyExists)
		}
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
