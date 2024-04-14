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
		ON CONFLICT (id) DO NOTHING
	`

	stmtCreateTags = `
		INSERT INTO tags(id)
		SELECT unnest($1::INTEGER[])
		ON CONFLICT (id) DO NOTHING
	`

	stmtCreateBanner = `
		INSERT INTO banners(content, is_active, feature_id)
		VALUES ($1::JSONB, $2, $3)
		RETURNING id
	`

	stmtCreateBannerRelations = `
		INSERT INTO banner_x_tag(banner_id, tag_id, feature_id)
		SELECT $1 AS banner_id, unnest($2::INTEGER[]), $3
	`
)

func (r *bannerRepository) CreateBanner(ctx context.Context, banner *domains.BannerWithTagIDs) (uint32, error) {
	fn := `bannerRepository.CreateBanner`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = tx.ExecContext(ctx, stmtCreateFeature, banner.FeatureID)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			r.log.Error(fmt.Sprintf("%s: tx error: %s", fn, txErr))
		}
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	var tagIDs pq.Int64Array
	for _, tagID := range banner.TagIDs {
		tagIDs = append(tagIDs, int64(tagID))
	}
	_, err = tx.ExecContext(ctx, stmtCreateTags, tagIDs)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			r.log.Error(fmt.Sprintf("%s: tx error: %s", fn, txErr))
		}
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	var id uint32
	err = tx.QueryRowContext(ctx, stmtCreateBanner, banner.Content, banner.IsActive, banner.FeatureID).
		Scan(&id)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			r.log.Error(fmt.Sprintf("%s: tx error: %s", fn, txErr))
		}
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = tx.ExecContext(ctx, stmtCreateBannerRelations, id, tagIDs, banner.FeatureID)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			r.log.Error(fmt.Sprintf("%s: tx error: %s", fn, txErr))
		}
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == pq.ErrorCode("23505") {
			return 0, fmt.Errorf("%s: %w", fn, ErrAlreadyExists)
		}
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return uint32(id), nil
}
