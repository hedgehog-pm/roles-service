package listroleassignments

import (
	"context"

	"github.com/hedgehog-pm/roles-service/internal/pagination"
)

type (
	reader interface {
		Page(ctx context.Context, filter Params) (pagination.Page[Assignment], error)
	}
	transactor interface {
		InReadTx(ctx context.Context, fn func(context.Context) error) error
	}
)
