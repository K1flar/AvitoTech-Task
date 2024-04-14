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
	_, err := a.db.Exec(`DELETE FROM features`)
	panicOnErr(err)
	_, err = a.db.Exec(`DELETE FROM tags`)
	panicOnErr(err)
	_, err = a.db.Exec(`DELETE FROM banners`)
	panicOnErr(err)
	_, err = a.db.Exec(`DELETE FROM banner_x_tags`)
	panicOnErr(err)
	_, err = a.db.Exec("ALTER SEQUENCE banners_id_seq RESTART WITH 1")
	panicOnErr(err)
}

func (a *testApp) mustAddBanner(banner *domains.Banner, tagIDs []int) {
	var pqTagIDs pq.Int64Array
	for _, tagID := range tagIDs {
		pqTagIDs = append(pqTagIDs, int64(tagID))
	}
	_, err := a.db.Exec(`
		INSERT INTO features(id)
		VALUES ($1)
		ON CONFLICT (id) DO NOTHING
	`, banner.FeatureID)
	panicOnErr(err)
	_, err = a.db.Exec(`
		INSERT INTO tags(id)
		SELECT unnest($1::INTEGER[])
		ON CONFLICT (id) DO NOTHING
	`, pqTagIDs)
	panicOnErr(err)
	_, err = a.db.Exec(`
		INSERT INTO banners(id, content, is_active, feature_id)
		VALUES ($1, $2::JSONB, $3, $4)
		RETURNING id
	`, banner.ID, banner.Content, banner.IsActive, banner.FeatureID)
	panicOnErr(err)
	_, err = a.db.Exec(`
		INSERT INTO banner_x_tag(banner_id, tag_id, feature_id)
		SELECT $1 AS banner_id, unnest($2::INTEGER[]), $3
	`, banner.ID, pqTagIDs, banner.FeatureID)
	panicOnErr(err)
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
