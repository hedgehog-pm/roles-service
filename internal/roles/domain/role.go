package domain

import (
	"fmt"
	"strings"
	"time"
)

type Role struct {
	ID          int
	Name        string
	Permissions PermissionSet
	Scope       ScopeRef
	Version     int64
	CreatedBy   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewRole(
	name string,
	scope ScopeRef,
	permissions PermissionSet,
	createdBy int,
) (*Role, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, NewError(CodeInvalidInput, "role name is empty")
	}

	if len(permissions.codes) == 0 {
		return nil, NewError(CodeInvalidInput, "role must have at least one permission")
	}

	if _, err := NewScopeRef(scope.Type(), scope.ID()); err != nil {
		return nil, fmt.Errorf("invalid role scope: %w", err)
	}

	if permissions.ScopeType() != scope.Type() {
		return nil, NewError(CodeInvalidInput, "permissions scope does not match role scope")
	}

	if createdBy <= 0 {
		return nil, NewError(CodeInvalidInput, "creator ID must be positive")
	}

	now := time.Now().UTC()

	return &Role{
		Name:        name,
		Permissions: permissions,
		Scope:       scope,
		Version:     1,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (r *Role) Rename(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return NewError(CodeInvalidInput, "role name is empty")
	}

	if name == r.Name {
		return nil
	}

	r.Name = name
	r.Version++
	r.UpdatedAt = time.Now().UTC()

	return nil
}

func (r *Role) SetPermissions(permissions PermissionSet) error {
	if permissions.ScopeType() != r.Scope.Type() {
		return NewError(CodeInvalidInput, "permissions scope does not match role scope")
	}

	if permissions.codes == nil {
		return NewError(CodeInvalidInput, "permission set is not initialized")
	}

	if len(permissions.codes) == 0 {
		return NewError(CodeInvalidInput, "role must have at least one permission")
	}

	if r.Permissions.Equal(permissions) {
		return nil
	}

	r.Permissions = permissions
	r.Version++
	r.UpdatedAt = time.Now().UTC()

	return nil
}
