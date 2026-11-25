package roleassignment

import (
	"context"

	"github.com/hedgehog-pm/roles-service/internal/domain"
)

type Usecase struct {
	repository repository
	transactor transactor
}

func New(repository repository, transactor transactor) *Usecase {
	return &Usecase{repository: repository, transactor: transactor}
}

func (s *Usecase) List(ctx context.Context, roleID int) ([]domain.RoleAssignment, error) {
	if roleID <= 0 {
		return nil, domain.NewError(domain.CodeInvalidInput, "role ID must be positive")
	}

	return s.repository.ListAssignments(ctx, roleID)
}

func (s *Usecase) Assign(ctx context.Context, actorID, roleID, userID int) (Result, error) {
	assignment, err := domain.NewRoleAssignment(roleID, userID, actorID)
	if err != nil {
		return Result{}, err
	}

	var result Result

	err = s.transactor.InTx(ctx, func(ctx context.Context) error {
		var err error

		result, err = s.repository.Ensure(ctx, assignment)

		return err
	})
	if err != nil {
		return Result{}, err
	}

	return result, nil
}

func (s *Usecase) Revoke(ctx context.Context, actorID, roleID, userID int) error {
	if actorID <= 0 || roleID <= 0 || userID <= 0 {
		return domain.NewError(domain.CodeInvalidInput, "assignment IDs must be positive")
	}

	return s.transactor.InTx(ctx, func(ctx context.Context) error { return s.repository.Revoke(ctx, roleID, userID) })
}

type Result struct {
	Assignment domain.RoleAssignment
	Created    bool
}
