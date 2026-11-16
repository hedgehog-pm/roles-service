package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Role struct {
	id               uuid.UUID
	name             string
	permissions      PermissionSet
	scope            ScopeRef
	sourceTemplateID *uuid.UUID
	version          int64
	createdBy        uuid.UUID
	createdAt        time.Time
	updatedAt        time.Time
}

func NewRole(
	name string,
	scope ScopeRef,
	permissions PermissionSet,
	createdBy uuid.UUID,
) (*Role, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("role name is empty")
	}

	if len(permissions.codes) == 0 {
		return nil, errors.New("role must have at least one permission")
	}

	if _, err := NewScopeRef(scope.Type(), scope.ID()); err != nil {
		return nil, fmt.Errorf("invalid role scope: %w", err)
	}

	if permissions.ScopeType() != scope.Type() {
		return nil, errors.New("permissions scope does not match role scope")
	}

	if createdBy == uuid.Nil {
		return nil, errors.New("creator ID is empty")
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("generate role ID: %w", err)
	}

	now := time.Now().UTC()

	return &Role{
		id:          id,
		name:        name,
		permissions: permissions,
		scope:       scope,
		version:     1,
		createdBy:   createdBy,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

func (r *Role) Rename(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("role name is empty")
	}

	if name == r.name {
		return nil
	}

	r.name = name
	r.version++
	r.updatedAt = time.Now().UTC()

	return nil
}

func (r *Role) SetPermissions(permissions PermissionSet) error {
	if permissions.ScopeType() != r.scope.Type() {
		return errors.New("permissions scope does not match role scope")
	}

	if permissions.codes == nil {
		return errors.New("permission set is not initialized")
	}

	if r.permissions.Equal(permissions) {
		return nil
	}

	if len(permissions.codes) == 0 {
		return errors.New("role must have at least one permission")
	}

	r.permissions = permissions
	r.version++
	r.updatedAt = time.Now().UTC()

	return nil
}

func (r *Role) ID() uuid.UUID {
	return r.id
}

func (r *Role) Name() string {
	return r.name
}

func (r *Role) Scope() ScopeRef {
	return r.scope
}

func (r *Role) Permissions() PermissionSet {
	return r.permissions
}

func (r *Role) SourceTemplateID() *uuid.UUID {
	if r.sourceTemplateID == nil {
		return nil
	}

	id := *r.sourceTemplateID

	return &id
}

func (r *Role) Version() int64 {
	return r.version
}

func (r *Role) CreatedBy() uuid.UUID {
	return r.createdBy
}

func (r *Role) CreatedAt() time.Time {
	return r.createdAt
}

func (r *Role) UpdatedAt() time.Time {
	return r.updatedAt
}
