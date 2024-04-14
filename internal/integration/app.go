package integration

import (
	"banner_service/cmd/server"
	"banner_service/internal/domains"
	"database/sql"

	"github.com/lib/pq"
)

type testApp struct {
	*server.App
	db *sql.DB
}

func (a *testApp) mustClearDB() {
	a.db.Exec(`DELETE FROM features`)
	a.db.Exec(`DELETE FROM tags`)
	a.db.Exec(`DELETE FROM banners`)
	a.db.Exec(`DELETE FROM banner_x_tags`)
	a.db.Exec("ALTER SEQUENCE banners_id_seq RESTART WITH 1")
}

func (a *testApp) mustAddBanner(banner *domains.Banner, tagIDs []int) {
	var pqTagIDs pq.Int64Array
	for _, tagID := range tagIDs {
		pqTagIDs = append(pqTagIDs, int64(tagID))
	}
	a.db.Exec(`
		INSERT INTO features(id)
		VALUES ($1)
		ON CONFLICT (id) DO NOTHING
	`, banner.FeatureID)
	a.db.Exec(`
		INSERT INTO tags(id)
		SELECT unnest($1::INTEGER[])
		ON CONFLICT (id) DO NOTHING
	`, pqTagIDs)
	a.db.Exec(`
		INSERT INTO banners(id, content, is_active, feature_id)
		VALUES ($1, $2::JSONB, $3, $4)
		RETURNING id
	`, banner.ID, banner.Content, banner.IsActive, banner.FeatureID)
	a.db.Exec(`
		INSERT INTO banner_x_tag(banner_id, tag_id, feature_id)
		SELECT $1 AS banner_id, unnest($2::INTEGER[]), $3
	`, banner.ID, pqTagIDs, banner.FeatureID)
}
