package application

import (
	"context"

	"github.com/hedgehog-pm/roles-service/internal/domain"
)

type AssignmentRepository interface {
	Ensure(context.Context, domain.RoleAssignment) (domain.AssignmentResult, error)
	Page(context.Context, domain.AssignmentFilter) (domain.Page[domain.AssignmentView], error)
	Assign(context.Context, domain.RoleAssignment) error
	Revoke(context.Context, int, int) error
	ListAssignments(context.Context, int) ([]domain.RoleAssignment, error)
}
type Assignments struct{ repository AssignmentRepository }

func NewAssignments(repository AssignmentRepository) *Assignments {
	return &Assignments{repository: repository}
}

func (s *Assignments) Assign(ctx context.Context, actorID, roleID, userID int) error {
	assignment, err := domain.NewRoleAssignment(roleID, userID, actorID)
	if err != nil {
		return err
	}

	return s.repository.Assign(ctx, assignment)
}

func (s *Assignments) Revoke(ctx context.Context, actorID, roleID, userID int) error {
	if actorID <= 0 || roleID <= 0 || userID <= 0 {
		return domain.NewError(domain.CodeInvalidInput, "assignment IDs must be positive")
	}

	return s.repository.Revoke(ctx, roleID, userID)
}

func (s *Assignments) List(ctx context.Context, roleID int) ([]domain.RoleAssignment, error) {
	if roleID <= 0 {
		return nil, domain.NewError(domain.CodeInvalidInput, "role ID must be positive")
	}

	return s.repository.ListAssignments(ctx, roleID)
}
