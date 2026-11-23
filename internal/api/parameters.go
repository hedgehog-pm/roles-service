package api

import (
	"github.com/hedgehog-pm/roles-service/internal/domain"
)

func optional[T any](v *T) T {
	if v != nil {
		return *v
	}

	var zero T

	return zero
}

func pagination(limit, offset *int) domain.Pagination {
	size := 20
	if limit != nil {
		size = *limit
	}

	return domain.Pagination{Limit: size, Offset: optional(offset)}
}
