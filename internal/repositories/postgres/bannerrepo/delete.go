package bannerrepo

import (
	"context"
	"fmt"
)

const (
	stmtDeleteBanner = `
		DELETE FROM banners
		WHERE id=$1
	`
)

func (r *bannerRepository) DeleteBannerByID(ctx context.Context, id uint32) error {
	fn := `bannerRepository.DeleteByID`

	res, err := r.db.ExecContext(ctx, stmtDeleteBanner, id)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	if rowsAff == 0 {
		return fmt.Errorf("%s: %w", fn, ErrNotFound)
	}

	return nil
}
