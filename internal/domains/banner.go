package domains

import "time"

type Banner struct {
	ID        uint32
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsActive  bool
}

type BannerKey struct {
	FeatureID uint32
	TagID     uint32
}

type BannerWithKey struct {
	BannerKey
	Banner
}
