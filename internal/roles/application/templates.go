package application

import (
	"context"

	"github.com/hedgehog-pm/roles-service/internal/roles/domain"
)

type (
	TemplateFilter struct {
		SystemID        int
		TargetScopeType domain.ScopeType
	}
	TemplateRepository interface {
		Page(context.Context, TemplateFilter, ListOptions) (Page[*domain.RoleTemplate], error)
		Get(context.Context, int) (*domain.RoleTemplate, error)

		Insert(context.Context, *domain.RoleTemplate) (*domain.RoleTemplate, error)
		Update(context.Context, *domain.RoleTemplate, int64) error
		Delete(context.Context, int, int64) error
	}
)

type Templates struct {
	repository TemplateRepository
	catalog    domain.PermissionCatalog
}

func NewTemplates(repository TemplateRepository, catalog domain.PermissionCatalog) *Templates {
	return &Templates{repository: repository, catalog: catalog}
}

type CreateTemplateInput struct {
	Name            string
	SystemID        int
	TargetScopeType domain.ScopeType
	PermissionCodes []domain.PermissionCode
}

func (s *Templates) Create(ctx context.Context, actorID int, input CreateTemplateInput) (*domain.RoleTemplate, error) {
	permissions, err := s.catalog.NewSet(input.TargetScopeType, input.PermissionCodes)
	if err != nil {
		return nil, err
	}

	value, err := domain.NewRoleTemplate(input.SystemID, input.Name, input.TargetScopeType, permissions, actorID)
	if err != nil {
		return nil, err
	}

	return s.repository.Insert(ctx, value)
}

type UpdateTemplateInput struct {
	ID              int
	Version         int64
	Name            string
	PermissionCodes []domain.PermissionCode
}

func (s *Templates) Update(ctx context.Context, actorID int, input UpdateTemplateInput) (*domain.RoleTemplate, error) {
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
		return nil, ErrConflict
	}

	permissions, err := s.catalog.NewSet(value.TargetScopeType, input.PermissionCodes)
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

func (s *Templates) Get(ctx context.Context, id int) (*domain.RoleTemplate, error) {
	if id <= 0 {
		return nil, domain.FieldError(domain.CodeInvalidInput, "id", "ID must be positive")
	}

	return s.repository.Get(ctx, id)
}

func (s *Templates) Delete(ctx context.Context, actorID, id int, version int64) error {
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
