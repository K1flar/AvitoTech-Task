package bannerrepo

import (
	"banner_service/internal/domains"
	"context"
	"fmt"

	"github.com/lib/pq"
)

const (
	stmtDeleteOldBanners = `
		DELETE FROM banners
		WHERE (id, updated_dttm) IN (
			SELECT id, updated_dttm 
			FROM banners
			WHERE id=$1
			ORDER BY updated_dttm DESC 
			OFFSET $2
		);
	`

	stmtCreateNewBannerVersion = `
		INSERT INTO banners(id, content, created_dttm, is_active, feature_id)
		SELECT $1, $2::JSONB, created_dttm, $3, $4 
		FROM banners
		WHERE id=$1
		LIMIT 1
	`
)

func (r *BannerRepository) UpdateBannerByID(ctx context.Context, id uint32, banner *domains.BannerWithTagIDs) error {
	fn := `BannerRepository.UpdateBannerByID`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = tx.ExecContext(ctx, stmtDeleteOldBanners, id, r.cfg.NumberOfBannerVersions)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = tx.ExecContext(ctx, stmtCreateFeature, banner.FeatureID)
	if err != nil {
		tx.Rollback()
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == pq.ErrorCode("23514") {
			return fmt.Errorf("%s: %w", fn, ErrInvalidFeatureID)
		}
		return fmt.Errorf("%s: %w", fn, err)
	}

	var tagIDs pq.Int64Array
	for _, tagID := range banner.TagIDs {
		tagIDs = append(tagIDs, int64(tagID))
	}
	_, err = tx.ExecContext(ctx, stmtCreateTags, tagIDs)
	if err != nil {
		tx.Rollback()
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == pq.ErrorCode("23514") {
			return fmt.Errorf("%s: %w", fn, ErrInvalidTagID)
		}
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = tx.ExecContext(ctx, stmtCreateNewBannerVersion, id, banner.Content, banner.IsActive, banner.FeatureID)
	if err != nil {
		tx.Rollback()
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == pq.ErrorCode("22P02") {
			return fmt.Errorf("%s: %w", fn, ErrNotJSON)
		}
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = tx.ExecContext(ctx, stmtCreateBannerRelations, tagIDs, id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", fn, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
