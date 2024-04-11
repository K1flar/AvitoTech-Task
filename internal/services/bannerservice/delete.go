package bannerservice

import (
	"banner_service/internal/repositories/postgres/bannerrepo"
	"context"
	"errors"
	"fmt"
)

func (s *bannerService) DeleteBannerByID(ctx context.Context, id uint32) error {
	fn := `bannerService.DeleteBannerByID`

	err := s.repo.DeleteBannerByID(ctx, id)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err))
		if errors.Is(err, bannerrepo.ErrNotFound) {
			return fmt.Errorf("%s: %w", fn, ErrNotFound)
		}
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
