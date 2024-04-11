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

	pagination := NewPaginationFromRequest(r)

	featureID, _ := strconv.Atoi(q.Get(QueryFeatureIDName))
	tagID, _ := strconv.Atoi(q.Get(QueryTagIDName))

	return NewBannerFilter(pagination, featureID, tagID)
}

func NewBannerFilter(pagination *Pagination, featureID, tagID int) *BannerFilter {
	filter := &BannerFilter{
		Pagination: pagination,
	}

	if featureID > 0 {
		filter.ByFeatureID = featureID
		filter.MustContainsFeatureID = true
	}

	if tagID > 0 {
		filter.ByTagID = tagID
		filter.MustContainsTagID = true
	}

	return filter
}

func (f *BannerFilter) Validate() {
	f.Pagination.Validate()

	if f.ByFeatureID <= 0 {
		f.MustContainsFeatureID = false
	}

	if f.ByTagID <= 0 {
		f.MustContainsTagID = false
	}
}
