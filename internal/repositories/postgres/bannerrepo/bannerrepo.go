package bannerrepo

// import (
// 	"banner_service/internal/domains"
// 	"database/sql"
// )

// type BannerRepository struct {
// 	db *sql.DB
// }

// func New(db *sql.DB) *BannerRepository {
// 	return &BannerRepository{
// 		db: db,
// 	}
// }

// func (r *BannerRepository) GetBannerByFeatureAndTagID(featureID, tagID uint32) *domains.Banner

// func (r *BannerRepository) GetBannersByKeys(keys []domains.BannerKey) (*domains.BannerWithKey, error) {
// 	fn := `BannerRepository.GetBannersByKeys`

// 	stmt := `
// 		SELECT
// 	`
// }
