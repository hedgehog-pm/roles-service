package roletemplate

import (
	"fmt"

	"github.com/hedgehog-pm/roles-service/internal/domain"
	rolesdb "github.com/hedgehog-pm/roles-service/internal/generated/bob/roletemplate"
)

func (s *Repository) restore(row rolesdb.GetTemplateRow) (*domain.RoleTemplate, error) {
	scopeType := domain.ScopeType(row.TargetScopeType)

	permissions, err := s.catalog.NewSet(scopeType, permissionCodes(row.Codes))
	if err != nil {
		return nil, fmt.Errorf("restore templates permissions: %w", err)
	}

	return &domain.RoleTemplate{
		ID: row.ID, Name: row.Name, SystemID: row.SystemID, TargetScopeType: scopeType,
		Permissions: permissions, Version: row.Version, CreatedBy: row.CreatedBy,
		CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}, nil
}
