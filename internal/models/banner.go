package models

import "banner_service/internal/domains"

type Banner struct {
	Content   domains.Content `json:"content"`
	IsActive  bool            `json:"is_active"`
	FeatureID uint32          `json:"feature_id"`
	TagIDs    []uint32        `json:"tag_ids"`
}
