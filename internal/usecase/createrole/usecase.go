package createrole

import (
	"context"

	"github.com/hedgehog-pm/roles-service/internal/domain"
)

type Usecase struct {
	repository repository
	transactor transactor
	catalog    domain.PermissionCatalog
}

func New(repository repository, transactor transactor, catalog domain.PermissionCatalog) *Usecase {
	return &Usecase{repository: repository, transactor: transactor, catalog: catalog}
}

type Params struct {
	Name            string
	Scope           domain.ScopeRef
	PermissionCodes []domain.PermissionCode
}

func (s *Usecase) Do(ctx context.Context, actorID int, input Params) (*domain.Role, error) {
	permissions, err := s.catalog.NewSet(input.Scope.Type(), input.PermissionCodes)
	if err != nil {
		return nil, err
	}

	value, err := domain.NewRole(input.Name, input.Scope, permissions, actorID)
	if err != nil {
		return nil, err
	}

	var result *domain.Role

	err = s.transactor.InTx(ctx, func(ctx context.Context) error {
		var err error

		result, err = s.repository.Insert(ctx, value)

		return err
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
