package application

import (
	"context"
	"time"

	"github.com/hedgehog-pm/roles-service/internal/roles/domain"
)

type Pagination struct {
	Limit  int
	Offset int
}

func (p *Pagination) Validate() error {
	if p.Limit == 0 {
		p.Limit = 20
	}

	if p.Limit < 1 || p.Limit > 100 || p.Offset < 0 {
		return domain.FieldError(domain.CodeInvalidInput, "pagination", "limit must be 1 to 100 and offset nonnegative")
	}

	return nil
}

type ListOptions struct {
	Pagination
	Search         string
	Sort           string
	Direction      string
	PermissionCode domain.PermissionCode
}

func (o *ListOptions) Validate() error {
	if err := o.Pagination.Validate(); err != nil {
		return err
	}

	if o.Sort == "" {
		o.Sort = "created_at"
	}

	if o.Direction == "" {
		o.Direction = "desc"
	}

	switch o.Sort {
	case "name", "created_at", "updated_at":
	default:
		return domain.FieldError(domain.CodeInvalidInput, "sort", "unsupported sort field")
	}

	if o.Direction != "asc" && o.Direction != "desc" {
		return domain.FieldError(domain.CodeInvalidInput, "direction", "direction must be asc or desc")
	}

	return nil
}

type (
	Page[T any] struct {
		Items  []T
		Total  int64
		Limit  int
		Offset int
	}
	AssignmentView struct {
		RoleID     int
		UserID     int
		RoleName   string
		Scope      domain.ScopeRef
		AssignedBy int
		AssignedAt time.Time
	}
	AssignmentFilter struct {
		RoleID    int
		UserID    int
		ScopeType domain.ScopeType
		ScopeID   int
		Pagination
	}
)

func (s *Roles) Page(ctx context.Context, scope domain.ScopeRef, options ListOptions) (Page[*domain.Role], error) {
	if _, err := domain.NewScopeRef(scope.Type(), scope.ID()); err != nil {
		return Page[*domain.Role]{}, err
	}

	if err := options.Validate(); err != nil {
		return Page[*domain.Role]{}, err
	}

	if options.PermissionCode != "" {
		if _, err := s.catalog.NewSet(scope.Type(), []domain.PermissionCode{options.PermissionCode}); err != nil {
			return Page[*domain.Role]{}, err
		}
	}

	return s.repository.Page(ctx, scope, options)
}

func (s *Templates) Page(ctx context.Context, filter TemplateFilter, options ListOptions) (Page[*domain.RoleTemplate], error) {
	if _, err := domain.NewScopeRef(filter.TargetScopeType, filter.SystemID); err != nil {
		return Page[*domain.RoleTemplate]{}, err
	}

	if err := options.Validate(); err != nil {
		return Page[*domain.RoleTemplate]{}, err
	}

	return s.repository.Page(ctx, filter, options)
}

func (s *Assignments) Page(ctx context.Context, filter AssignmentFilter) (Page[AssignmentView], error) {
	if filter.RoleID < 0 || filter.UserID < 0 || (filter.RoleID == 0 && (filter.ScopeID <= 0 || filter.ScopeType == "")) {
		return Page[AssignmentView]{}, domain.NewError(domain.CodeInvalidInput, "role or user ID is required")
	}

	if filter.ScopeID < 0 || (filter.ScopeID > 0 && filter.ScopeType == "") {
		return Page[AssignmentView]{}, domain.FieldError(domain.CodeInvalidInput, "scope_id", "scope_id requires scope_type")
	}

	if filter.ScopeType != "" {
		if _, err := domain.NewScopeRef(filter.ScopeType, 1); err != nil {
			return Page[AssignmentView]{}, err
		}
	}

	if err := filter.Validate(); err != nil {
		return Page[AssignmentView]{}, err
	}

	return s.repository.Page(ctx, filter)
}

type AssignmentResult struct {
	Assignment domain.RoleAssignment
	Created    bool
}

func (s *Assignments) Ensure(ctx context.Context, actorID, roleID, userID int) (AssignmentResult, error) {
	assignment, err := domain.NewRoleAssignment(roleID, userID, actorID)
	if err != nil {
		return AssignmentResult{}, err
	}

	return s.repository.Ensure(ctx, assignment)
}
