package roletemplate

import (
	"context"

	"github.com/hedgehog-pm/roles-service/internal/domain"
	"github.com/hedgehog-pm/roles-service/internal/pagination"
)

type repository interface {
	Delete(ctx context.Context, id int, version int64) error
	Get(ctx context.Context, id int) (*domain.RoleTemplate, error)
	Page(ctx context.Context, filter Filter, options pagination.Options) (pagination.Page[*domain.RoleTemplate], error)
}
type transactor interface {
	InTx(ctx context.Context, fn func(context.Context) error) error
	InReadTx(ctx context.Context, fn func(context.Context) error) error
}
