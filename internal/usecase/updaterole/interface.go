package updaterole

import (
	"context"

	"github.com/hedgehog-pm/roles-service/internal/domain"
)

type repository interface {
	Get(ctx context.Context, id int) (*domain.Role, error)
	Update(ctx context.Context, value *domain.Role, expected int64) error
}
type transactor interface {
	InTx(ctx context.Context, fn func(context.Context) error) error
}
