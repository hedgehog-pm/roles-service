package roleassignment

import (
	"context"

	"github.com/hedgehog-pm/roles-service/internal/domain"
)

type repository interface {
	Ensure(ctx context.Context, value domain.RoleAssignment) (Result, error)
	Revoke(ctx context.Context, roleID, userID int) error
	ListAssignments(ctx context.Context, roleID int) ([]domain.RoleAssignment, error)
}
type transactor interface {
	InTx(ctx context.Context, fn func(context.Context) error) error
}
