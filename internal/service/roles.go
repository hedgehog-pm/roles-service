package application

import (
	"context"

	"github.com/hedgehog-pm/roles-service/internal/domain"
)

type RoleRepository interface {
	Page(context.Context, domain.ScopeRef, domain.ListOptions) (domain.Page[*domain.Role], error)
	Get(context.Context, int) (*domain.Role, error)

	Insert(context.Context, *domain.Role) (*domain.Role, error)
	Update(context.Context, *domain.Role, int64) error
	Delete(context.Context, int, int64) error
}
type Roles struct {
	repository RoleRepository
	catalog    domain.PermissionCatalog
}

func NewRoles(repository RoleRepository, catalog domain.PermissionCatalog) *Roles {
	return &Roles{repository: repository, catalog: catalog}
}

type CreateRoleInput struct {
	Name            string
	Scope           domain.ScopeRef
	PermissionCodes []domain.PermissionCode
}

func (s *Roles) Create(ctx context.Context, actorID int, input CreateRoleInput) (*domain.Role, error) {
	permissions, err := s.catalog.NewSet(input.Scope.Type(), input.PermissionCodes)
	if err != nil {
		return nil, err
	}

	value, err := domain.NewRole(input.Name, input.Scope, permissions, actorID)
	if err != nil {
		return nil, err
	}

	return s.repository.Insert(ctx, value)
}

type UpdateRoleInput struct {
	ID              int
	Version         int64
	Name            string
	PermissionCodes []domain.PermissionCode
}

func (s *Roles) Update(ctx context.Context, actorID int, input UpdateRoleInput) (*domain.Role, error) {
	if actorID <= 0 {
		return nil, domain.NewError(domain.CodeInvalidInput, "actor ID must be positive")
	}

	if input.ID <= 0 || input.Version <= 0 {
		return nil, domain.FieldError(domain.CodeInvalidInput, "expected_version", "ID and expected_version must be positive")
	}

	value, err := s.repository.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if value.Version != input.Version {
		return nil, domain.ErrConflict
	}

	permissions, err := s.catalog.NewSet(value.Scope.Type(), input.PermissionCodes)
	if err != nil {
		return nil, err
	}

	if err = value.Rename(input.Name); err != nil {
		return nil, err
	}

	if err = value.SetPermissions(permissions); err != nil {
		return nil, err
	}

	if value.Version != input.Version {
		value.Version = input.Version + 1
	}

	if err = s.repository.Update(ctx, value, input.Version); err != nil {
		return nil, err
	}

	return value, nil
}

func (s *Roles) Get(ctx context.Context, id int) (*domain.Role, error) {
	if id <= 0 {
		return nil, domain.FieldError(domain.CodeInvalidInput, "id", "ID must be positive")
	}

	return s.repository.Get(ctx, id)
}

func (s *Roles) Delete(ctx context.Context, actorID, id int, version int64) error {
	if actorID <= 0 {
		return domain.NewError(domain.CodeInvalidInput, "actor ID must be positive")
	}

	if id <= 0 {
		return domain.FieldError(domain.CodeInvalidInput, "id", "ID must be positive")
	}

	if version < 1 {
		return domain.NewError(domain.CodeInvalidInput, "version must be positive")
	}

	return s.repository.Delete(ctx, id, version)
}
