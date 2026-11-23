package domain

import "time"

type Pagination struct {
	Limit  int
	Offset int
}

func (p *Pagination) Validate() error {
	if p.Limit == 0 {
		p.Limit = 20
	}

	if p.Limit < 1 || p.Limit > 100 || p.Offset < 0 {
		return FieldError(CodeInvalidInput, "pagination", "limit must be 1 to 100 and offset nonnegative")
	}

	return nil
}

type ListOptions struct {
	Pagination
	Search         string
	Sort           string
	Direction      string
	PermissionCode PermissionCode
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
		return FieldError(CodeInvalidInput, "sort", "unsupported sort field")
	}

	if o.Direction != "asc" && o.Direction != "desc" {
		return FieldError(CodeInvalidInput, "direction", "direction must be asc or desc")
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
		Scope      ScopeRef
		AssignedBy int
		AssignedAt time.Time
	}
	AssignmentFilter struct {
		RoleID    int
		UserID    int
		ScopeType ScopeType
		ScopeID   int
		Pagination
	}
)

type AssignmentResult struct {
	Assignment RoleAssignment
	Created    bool
}

type TemplateFilter struct {
	SystemID        int
	TargetScopeType ScopeType
}
