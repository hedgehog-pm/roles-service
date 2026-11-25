package domain

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
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
	if name == "" || utf8.RuneCountInString(name) > 200 {
		return nil, FieldError(CodeInvalidInput, "name", "name must contain 1 to 200 Unicode characters")
	}

	if len(permissions.codes) == 0 {
		return nil, FieldError(CodeEmptyPermissions, "permission_codes", "at least one permission is required")
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
	if name == "" || utf8.RuneCountInString(name) > 200 {
		return FieldError(CodeInvalidInput, "name", "name must contain 1 to 200 Unicode characters")
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
		return FieldError(CodeEmptyPermissions, "permission_codes", "at least one permission is required")
	}

	if len(permissions.codes) == 0 {
		return FieldError(CodeEmptyPermissions, "permission_codes", "at least one permission is required")
	}

	if r.Permissions.Equal(permissions) {
		return nil
	}

	r.Permissions = permissions
	r.Version++
	r.UpdatedAt = time.Now().UTC()

	return nil
}
