package bannerrepo

import (
	"banner_service/internal/domains"
	"banner_service/pkg/filters"
	"context"
	"database/sql"
	"fmt"

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
			JOIN banners_x_tags AS bt 
			ON bt.banner_id=b.id
		WHERE (b.feature_id, bt.tag_id) IN (($1, $2)) 
			AND b.is_active 
			AND b.updated_dttm=(SELECT MAX(updated_dttm) FROM banners WHERE id=b.id)
	`

	stmtGetBanners = `
		SELECT b.id, b.content, b.created_dttm, b.updated_dttm, b.is_active, b.feature_id, array_agg(tag_id) AS tag_ids FROM 
		banners AS b JOIN banners_x_tags AS bt ON b.id=bt.banner_id AND b.updated_dttm=bt.banner_updated_dttm
		WHERE b.updated_dttm=(SELECT MAX(updated_dttm) FROM banners WHERE id=b.id)
			AND (b.feature_id=$1 OR $1=0)
			AND (b.id, b.updated_dttm) IN(SELECT bt2.banner_id, bt2.banner_updated_dttm FROM banners_x_tags as bt2 WHERE bt2.tag_id=$2 OR $2=0)
		GROUP BY b.id,b.content,b.created_dttm,b.updated_dttm
		LIMIT $3
		OFFSET $4
	`
)

func (r *BannerRepository) GetBannerByFeatureAndTagID(ctx context.Context, featureID, tagID uint32) (*domains.Banner, error) {
	fn := `BannerRepository.GetBannersByKeys`

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

func (r *BannerRepository) GetBanners(ctx context.Context, filter *filters.BannerFilter) ([]*domains.BannerWithTagIDs, error) {
	fn := `BannerRepository.GetBanners`

	var banners []*domains.BannerWithTagIDs
	rows, err := r.db.QueryContext(ctx, stmtGetBanners, filter.ByFeatureID, filter.ByTagID, filter.Pagination.Limit, filter.Pagination.Offset)
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

// TODO: DELETE
func (r *BannerRepository) Test() error {
	fn := `BannerRepository.Test`

	stmt := `
	SELECT t.banner_updated_dttm
    FROM banners_x_tags AS t
    WHERE t.banner_id = 1
	`

	_, err := r.db.Query(stmt)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
