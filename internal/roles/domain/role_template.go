package domain

import (
	"strings"
	"time"
)

type RoleTemplate struct {
	ID              int
	SystemID        int
	Name            string
	TargetScopeType ScopeType
	Permissions     PermissionSet
	Version         int64
	CreatedBy       int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewRoleTemplate(
	systemID int,
	name string,
	targetScopeType ScopeType,
	permissions PermissionSet,
	createdBy int,
) (*RoleTemplate, error) {
	if systemID <= 0 {
		return nil, NewError(CodeInvalidInput, "system ID must be positive")
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, NewError(CodeInvalidInput, "role template name is empty")
	}

	switch targetScopeType {
	case ScopeSystem, ScopeWorkspace, ScopeProject:
	default:
		return nil, NewError(CodeInvalidInput, "unknown target scope type")
	}

	if permissions.ScopeType() != targetScopeType {
		return nil, NewError(CodeInvalidInput,
			"permissions scope does not match template target scope",
		)
	}

	if permissions.codes == nil {
		return nil, NewError(CodeInvalidInput, "permission set is not initialized")
	}

	if len(permissions.codes) == 0 {
		return nil, NewError(CodeInvalidInput, "role template permissions are empty")
	}

	if createdBy <= 0 {
		return nil, NewError(CodeInvalidInput, "creator ID must be positive")
	}

	now := time.Now().UTC()

	return &RoleTemplate{
		SystemID:        systemID,
		Name:            name,
		TargetScopeType: targetScopeType,
		Permissions:     permissions,
		Version:         1,
		CreatedBy:       createdBy,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

func (t *RoleTemplate) Rename(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return NewError(CodeInvalidInput, "role template name is empty")
	}

	if name == t.Name {
		return nil
	}

	t.Name = name
	t.Version++
	t.UpdatedAt = time.Now().UTC()

	return nil
}

func (t *RoleTemplate) SetPermissions(permissions PermissionSet) error {
	if permissions.ScopeType() != t.TargetScopeType {
		return NewError(CodeInvalidInput,
			"permissions scope does not match template target scope",
		)
	}

	if permissions.codes == nil {
		return NewError(CodeInvalidInput, "permission set is not initialized")
	}

	if len(permissions.codes) == 0 {
		return NewError(CodeInvalidInput, "role template must have at least one permission")
	}

	if t.Permissions.Equal(permissions) {
		return nil
	}

	t.Permissions = permissions
	t.Version++
	t.UpdatedAt = time.Now().UTC()

	return nil
}
