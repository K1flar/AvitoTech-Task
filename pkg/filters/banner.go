package filters

import (
	"net/http"
	"strconv"
)

const (
	QueryFeatureIDName = "feature_id"
	QueryTagIDName     = "tag_id"
)

type BannerFilter struct {
	Pagination            *Pagination
	ByFeatureID           int
	MustContainsFeatureID bool
	ByTagID               int
	MustContainsTagID     bool
}

func NewBannerFilterFromRequest(r *http.Request) *BannerFilter {
	q := r.URL.Query()
	filter := &BannerFilter{
		Pagination: NewPaginationFromRequest(r),
	}

	if featureID, err := strconv.Atoi(q.Get(QueryFeatureIDName)); err == nil {
		filter.ByFeatureID = featureID
		filter.MustContainsFeatureID = true
	}
	if tagID, err := strconv.Atoi(q.Get(QueryTagIDName)); err == nil {
		filter.ByFeatureID = tagID
		filter.MustContainsTagID = true
	}

	return filter
}
