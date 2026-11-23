package listroles

import (
	"context"

	"github.com/hedgehog-pm/roles-service/internal/domain"
	"github.com/hedgehog-pm/roles-service/internal/pagination"
)

type (
	reader interface {
		Page(ctx context.Context, scope domain.ScopeRef, options Params) (pagination.Page[*domain.Role], error)
	}
	transactor interface {
		InReadTx(ctx context.Context, fn func(context.Context) error) error
	}
)
