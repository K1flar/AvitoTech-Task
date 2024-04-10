package bannerrepo

import (
	"banner_service/internal/domains"
	"context"
	"fmt"

	"github.com/lib/pq"
)

const (
	stmtCreateFeature = `
		INSERT INTO features(id)
		VALUES ($1)
		ON CONFLICT (id) DO NOTHING;
	`

	stmtCreateTags = `
		INSERT INTO tags(id)
		SELECT unnest($1::INTEGER[])
		ON CONFLICT (id) DO NOTHING;
	`

	stmtCreateBanner = `
		INSERT INTO banners(content, is_active, feature_id)
		VALUES ($1::JSONB, $2, $3)
		RETURNING id;
	`

	stmtCreateBannerRelations = `
		INSERT INTO banners_x_tags(tag_id, banner_id, banner_updated_dttm)
		SELECT unnest($1::INTEGER[]), id, MAX(updated_dttm)
		FROM banners WHERE id=$2
		GROUP BY id;
	`
)

func (r *BannerRepository) CreateBanner(ctx context.Context, banner *domains.BannerWithTagIDs) (uint32, error) {
	fn := `BannerRepository.CreateBanner`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = tx.ExecContext(ctx, stmtCreateFeature, banner.FeatureID)
	if err != nil {
		tx.Rollback()
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == pq.ErrorCode("23514") {
			return 0, fmt.Errorf("%s: %w", fn, ErrInvalidFeatureID)
		}
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	var tagIDs pq.Int64Array
	for _, tagID := range banner.TagIDs {
		tagIDs = append(tagIDs, int64(tagID))
	}
	_, err = tx.ExecContext(ctx, stmtCreateTags, tagIDs)
	if err != nil {
		tx.Rollback()
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == pq.ErrorCode("23514") {
			return 0, fmt.Errorf("%s: %w", fn, ErrInvalidTagID)
		}
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	var id uint32
	err = tx.QueryRowContext(ctx, stmtCreateBanner, banner.Content, banner.IsActive, banner.FeatureID).
		Scan(&id)
	if err != nil {
		tx.Rollback()
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == pq.ErrorCode("22P02") {
			return 0, fmt.Errorf("%s: %w", fn, ErrNotJSON)
		}
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = tx.ExecContext(ctx, stmtCreateBannerRelations, tagIDs, id)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return uint32(id), nil
}
