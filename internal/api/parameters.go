package api

import (
	"github.com/samber/lo"

	"github.com/hedgehog-pm/roles-service/internal/pagination"
)

func pageParams(limit, offset *int) pagination.Pagination {
	size := 20
	if limit != nil {
		size = *limit
	}

	return pagination.Pagination{Limit: size, Offset: lo.FromPtr(offset)}
}
