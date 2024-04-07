package domains

import "time"

type Banner struct {
	ID        uint32
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsActive  bool
}
