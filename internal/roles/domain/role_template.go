package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RoleTemplate struct {
	id              uuid.UUID
	systemID        uuid.UUID
	name            string
	targetScopeType ScopeType
	permissions     PermissionSet
	version         int64
	createdBy       uuid.UUID
	createdAt       time.Time
	updatedAt       time.Time
}

func NewRoleTemplate(
	systemID uuid.UUID,
	name string,
	targetScopeType ScopeType,
	permissions PermissionSet,
	createdBy uuid.UUID,
) (*RoleTemplate, error) {
	if systemID == uuid.Nil {
		return nil, errors.New("system ID is empty")
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("role template name is empty")
	}

	switch targetScopeType {
	case ScopeSystem, ScopeWorkspace, ScopeProject:
	default:
		return nil, errors.New("unknown target scope type")
	}

	if permissions.ScopeType() != targetScopeType {
		return nil, errors.New(
			"permissions scope does not match template target scope",
		)
	}

	if permissions.codes == nil {
		return nil, errors.New("permission set is not initialized")
	}

	if len(permissions.codes) == 0 {
		return nil, errors.New("role template permissions are empty")
	}

	if createdBy == uuid.Nil {
		return nil, errors.New("creator ID is empty")
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("generate role template ID: %w", err)
	}

	now := time.Now().UTC()

	return &RoleTemplate{
		id:              id,
		systemID:        systemID,
		name:            name,
		targetScopeType: targetScopeType,
		permissions:     permissions,
		version:         1,
		createdBy:       createdBy,
		createdAt:       now,
		updatedAt:       now,
	}, nil
}

func (t *RoleTemplate) Rename(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("role template name is empty")
	}

	if name == t.name {
		return nil
	}

	t.name = name
	t.version++
	t.updatedAt = time.Now().UTC()

	return nil
}

func (t *RoleTemplate) SetPermissions(permissions PermissionSet) error {
	if permissions.ScopeType() != t.targetScopeType {
		return errors.New(
			"permissions scope does not match template target scope",
		)
	}

	if permissions.codes == nil {
		return errors.New("permission set is not initialized")
	}

	if len(permissions.codes) == 0 {
		return errors.New("role template must have at least one permission")
	}

	if t.permissions.Equal(permissions) {
		return nil
	}

	t.permissions = permissions
	t.version++
	t.updatedAt = time.Now().UTC()

	return nil
}

func (t *RoleTemplate) ID() uuid.UUID {
	return t.id
}

func (t *RoleTemplate) SystemID() uuid.UUID {
	return t.systemID
}

func (t *RoleTemplate) Name() string {
	return t.name
}

func (t *RoleTemplate) TargetScopeType() ScopeType {
	return t.targetScopeType
}

func (t *RoleTemplate) Permissions() PermissionSet {
	return t.permissions
}

func (t *RoleTemplate) Version() int64 {
	return t.version
}

func (t *RoleTemplate) CreatedBy() uuid.UUID {
	return t.createdBy
}

func (t *RoleTemplate) CreatedAt() time.Time {
	return t.createdAt
}

func (t *RoleTemplate) UpdatedAt() time.Time {
	return t.updatedAt
}
