package bannerrepo

import (
	"banner_service/internal/domains"
	"banner_service/pkg/filters"
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

const (
	stmtGetBannerByFeatureAndTagID = `
		SELECT 
			b.id,
			b.content,
			b.created_dttm,
			b.updated_dttm,
			b.is_active,
			b.feature_id
		FROM 
			banners AS b 
			JOIN banner_x_tag AS bxt 
			ON b.id=bxt.banner_id
		WHERE b.feature_id = $1 AND bxt.tag_id = $2
	`
)

func (r *bannerRepository) GetBannerByFeatureAndTagID(ctx context.Context, featureID, tagID uint32) (*domains.Banner, error) {
	fn := `bannerRepository.GetBannerByFeatureAndTagID`

	var banner domains.Banner
	err := r.db.QueryRowContext(ctx, stmtGetBannerByFeatureAndTagID, featureID, tagID).
		Scan(&banner.ID, &banner.Content, &banner.CreatedAt, &banner.UpdatedAt, &banner.IsActive, &banner.FeatureID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", fn, ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &banner, nil
}

func (r *bannerRepository) GetBanners(ctx context.Context, filter *filters.BannerFilter) ([]*domains.BannerWithTagIDs, error) {
	fn := `bannerRepository.GetBanners`
	var stmt strings.Builder
	stmt.WriteString(`
		SELECT 
			b.id, 
			b.content, 
			b.created_dttm, 
			b.updated_dttm, 
			b.is_active, 
			b.feature_id, 
			array_agg(tag_id) AS tag_ids 
		FROM 
			banners AS b 
			JOIN banner_x_tag AS bxt 
			ON b.id=bxt.banner_id
		WHERE 1=1 
	`)

	args := []any{}

	if filter.MustContainsFeatureID {
		args = append(args, filter.ByFeatureID)
		stmt.WriteString(fmt.Sprintf(" AND b.feature_id=$%d", len(args)))
	}

	if filter.MustContainsTagID {
		args = append(args, filter.ByTagID)
		stmt.WriteString(fmt.Sprintf(` AND b.id IN (SELECT banner_id FROM banner_x_tag WHERE tag_id=$%d)`, len(args)))
	}

	stmt.WriteString(" GROUP BY b.id,b.content,b.created_dttm,b.updated_dttm")
	args = append(args, filter.Pagination.Limit, filter.Pagination.Offset)
	stmt.WriteString(fmt.Sprintf(`
		LIMIT $%d
		OFFSET $%d
	`, len(args)-1, len(args)))

	banners := []*domains.BannerWithTagIDs{}
	rows, err := r.db.QueryContext(ctx, stmt.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	for rows.Next() {
		banner := &domains.BannerWithTagIDs{}
		var tagIDs pq.Int64Array
		err := rows.Scan(&banner.ID, &banner.Content, &banner.CreatedAt, &banner.UpdatedAt, &banner.IsActive, &banner.FeatureID, &tagIDs)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		for _, tagID := range tagIDs {
			banner.TagIDs = append(banner.TagIDs, uint32(tagID))
		}
		banners = append(banners, banner)
	}

	return banners, nil
}
