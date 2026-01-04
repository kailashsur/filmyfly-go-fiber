package utils

import "math"

type Pagination struct {
	Page        int   `json:"page"`
	Limit       int   `json:"limit"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"totalPages"`
	HasNextPage bool  `json:"hasNextPage"`
	HasPrevPage bool  `json:"hasPrevPage"`
}

// NewPagination creates a pagination object
func NewPagination(page, limit int, total int64) *Pagination {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &Pagination{
		Page:        page,
		Limit:       limit,
		Total:       total,
		TotalPages:  totalPages,
		HasNextPage: page < totalPages,
		HasPrevPage: page > 1,
	}
}

// GetOffset calculates the database offset for pagination
func GetOffset(page, limit int) int {
	return (page - 1) * limit
}
