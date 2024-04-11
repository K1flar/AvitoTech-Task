package domains

import "time"

type Banner struct {
	ID        uint32
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsActive  bool
	FeatureID uint32
}

type BannerKey struct {
	FeatureID uint32
	TagID     uint32
}

type BannerWithTagIDs struct {
	Banner
	TagIDs []uint32
}
