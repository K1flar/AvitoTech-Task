package filters

import "net/http"

const (
	QueryLimitName  = "limit"
	QueryOffsetName = "offset"

	DefaultLimit  = 10
	DefaultOffset = 0
)

type Pagination struct {
	Limit  int
	Offset int
}

func NewPagination(limit, offset int) *Pagination {
	return &Pagination{
		Limit:  limit,
		Offset: offset,
	}
}

func NewPaginationFromRequest(r *http.Request) *Pagination {
	q := r.URL.Query()
	limit := parseIntWithDefaultValue(q.Get(QueryLimitName), DefaultLimit)
	offset := parseIntWithDefaultValue(q.Get(QueryOffsetName), DefaultOffset)
	return NewPagination(limit, offset)
}

func (p *Pagination) Validate() {
	if p.Limit < 0 {
		p.Limit = DefaultLimit
	}

	if p.Offset < 0 {
		p.Offset = DefaultOffset
	}
}
