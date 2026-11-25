package createtemplate

import (
	"context"

	"github.com/hedgehog-pm/roles-service/internal/domain"
)

type repository interface {
	Insert(ctx context.Context, value *domain.RoleTemplate) (*domain.RoleTemplate, error)
}
type transactor interface {
	InTx(ctx context.Context, fn func(context.Context) error) error
}
