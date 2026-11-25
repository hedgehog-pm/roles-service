package role

import (
	"fmt"

	"github.com/hedgehog-pm/roles-service/internal/domain"
	rolesdb "github.com/hedgehog-pm/roles-service/internal/generated/bob/role"
)

func (s *Repository) restore(row rolesdb.GetRoleRow) (*domain.Role, error) {
	scopeType := domain.ScopeType(row.ScopeType)

	permissions, err := s.catalog.NewSet(scopeType, permissionCodes(row.Codes))
	if err != nil {
		return nil, fmt.Errorf("restore roles permissions: %w", err)
	}

	scope, err := domain.NewScopeRef(scopeType, row.ScopeID)
	if err != nil {
		return nil, err
	}

	return &domain.Role{
		ID: row.ID, Name: row.Name, Scope: scope,
		Permissions: permissions, Version: row.Version, CreatedBy: row.CreatedBy,
		CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}, nil
}
