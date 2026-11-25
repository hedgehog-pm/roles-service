package permission

import "github.com/hedgehog-pm/roles-service/internal/domain"

type Usecase struct{ catalog domain.PermissionCatalog }

func New(catalog domain.PermissionCatalog) *Usecase {
	return &Usecase{catalog: catalog}
}

func (s *Usecase) List(scope domain.ScopeType) ([]domain.PermissionDescriptor, error) {
	if _, err := s.catalog.NewSet(scope, nil); err != nil {
		return nil, err
	}

	return s.catalog.List(scope), nil
}

func (s *Usecase) Describe(set domain.PermissionSet) []domain.PermissionDescriptor {
	result := make([]domain.PermissionDescriptor, 0, len(set.Codes()))
	for _, code := range set.Codes() {
		if value, ok := s.catalog.LookUp(code); ok {
			result = append(result, value)
		}
	}

	return result
}
