package application

import (
	"context"

	"github.com/hedgehog-pm/roles-service/internal/domain"
)

func (s *Roles) Page(ctx context.Context, scope domain.ScopeRef, options domain.ListOptions) (domain.Page[*domain.Role], error) {
	if _, err := domain.NewScopeRef(scope.Type(), scope.ID()); err != nil {
		return domain.Page[*domain.Role]{}, err
	}

	if err := options.Validate(); err != nil {
		return domain.Page[*domain.Role]{}, err
	}

	if options.PermissionCode != "" {
		if _, err := s.catalog.NewSet(scope.Type(), []domain.PermissionCode{options.PermissionCode}); err != nil {
			return domain.Page[*domain.Role]{}, err
		}
	}

	return s.repository.Page(ctx, scope, options)
}

func (s *Templates) Page(ctx context.Context, filter domain.TemplateFilter, options domain.ListOptions) (domain.Page[*domain.RoleTemplate], error) {
	if _, err := domain.NewScopeRef(filter.TargetScopeType, filter.SystemID); err != nil {
		return domain.Page[*domain.RoleTemplate]{}, err
	}

	if err := options.Validate(); err != nil {
		return domain.Page[*domain.RoleTemplate]{}, err
	}

	return s.repository.Page(ctx, filter, options)
}

func (s *Assignments) Page(ctx context.Context, filter domain.AssignmentFilter) (domain.Page[domain.AssignmentView], error) {
	if filter.RoleID < 0 || filter.UserID < 0 || (filter.RoleID == 0 && (filter.ScopeID <= 0 || filter.ScopeType == "")) {
		return domain.Page[domain.AssignmentView]{}, domain.NewError(domain.CodeInvalidInput, "role or user ID is required")
	}

	if filter.ScopeID < 0 || (filter.ScopeID > 0 && filter.ScopeType == "") {
		return domain.Page[domain.AssignmentView]{}, domain.FieldError(domain.CodeInvalidInput, "scope_id", "scope_id requires scope_type")
	}

	if filter.ScopeType != "" {
		if _, err := domain.NewScopeRef(filter.ScopeType, 1); err != nil {
			return domain.Page[domain.AssignmentView]{}, err
		}
	}

	if err := filter.Validate(); err != nil {
		return domain.Page[domain.AssignmentView]{}, err
	}

	return s.repository.Page(ctx, filter)
}

func (s *Assignments) Ensure(ctx context.Context, actorID, roleID, userID int) (domain.AssignmentResult, error) {
	assignment, err := domain.NewRoleAssignment(roleID, userID, actorID)
	if err != nil {
		return domain.AssignmentResult{}, err
	}

	return s.repository.Ensure(ctx, assignment)
}
