package domains

import "time"

type Content string

func (c *Content) MarshalJSON() ([]byte, error) {
	return []byte(*c), nil
}

func (c *Content) UnmarshalJSON(data []byte) error {
	*c = Content(string(data))
	return nil
}

func (c *Content) String() string {
	return string(*c)
}

type Banner struct {
	ID        uint32    `json:"banner_id"`
	Content   Content   `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsActive  bool      `json:"is_active"`
	FeatureID uint32    `json:"feature_id"`
}

type BannerKey struct {
	FeatureID uint32
	TagID     uint32
}

type BannerWithTagIDs struct {
	Banner
	TagIDs []uint32 `json:"tag_ids"`
}
