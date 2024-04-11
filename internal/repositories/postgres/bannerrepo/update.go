package bannerrepo

import (
	"banner_service/internal/domains"
	"context"
	"fmt"

	"github.com/lib/pq"
)

const (
	stmtDeleteBannerRelations = `
		DELETE FROM banner_x_tag
		WHERE banner_id=$1
	`

	stmtUpdateBanner = `
		UPDATE banners
		SET (
			content, 
			updated_dttm, 
			is_active, 
			feature_id
		) = ($2::JSONB, CURRENT_TIMESTAMP, $3, $4)
		WHERE id=$1
	`
)

func (r *bannerRepository) UpdateBannerByID(ctx context.Context, id uint32, banner *domains.BannerWithTagIDs) error {
	fn := `bannerRepository.UpdateBannerByID`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = tx.ExecContext(ctx, stmtCreateFeature, banner.FeatureID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", fn, err)
	}

	res, err := tx.ExecContext(ctx, stmtUpdateBanner, id, banner.Content, banner.IsActive, banner.FeatureID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", fn, err)
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", fn, err)
	}

	if rowsAff == 0 {
		tx.Rollback()
		return fmt.Errorf("%s: %w", fn, ErrNotFound)
	}

	var tagIDs pq.Int64Array
	for _, tagID := range banner.TagIDs {
		tagIDs = append(tagIDs, int64(tagID))
	}
	_, err = tx.ExecContext(ctx, stmtCreateTags, tagIDs)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = tx.ExecContext(ctx, stmtDeleteBannerRelations, id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = tx.ExecContext(ctx, stmtCreateBannerRelations, id, tagIDs, banner.FeatureID)
	if err != nil {
		tx.Rollback()
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == pq.ErrorCode("23505") {
			return fmt.Errorf("%s: %w", fn, ErrAlreadyExists)
		}
		return fmt.Errorf("%s: %w", fn, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
