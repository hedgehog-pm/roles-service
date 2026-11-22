package application

import "github.com/hedgehog-pm/roles-service/internal/roles/domain"

type Permissions struct{ catalog domain.PermissionCatalog }

func NewPermissions(catalog domain.PermissionCatalog) *Permissions {
	return &Permissions{catalog: catalog}
}

func (s *Permissions) List(scope domain.ScopeType) ([]domain.PermissionDescriptor, error) {
	if _, err := s.catalog.NewSet(scope, nil); err != nil {
		return nil, err
	}

	return s.catalog.List(scope), nil
}

func (s *Permissions) Describe(set domain.PermissionSet) []domain.PermissionDescriptor {
	result := make([]domain.PermissionDescriptor, 0, len(set.Codes()))
	for _, code := range set.Codes() {
		if value, ok := s.catalog.LookUp(code); ok {
			result = append(result, value)
		}
	}

	return result
}
