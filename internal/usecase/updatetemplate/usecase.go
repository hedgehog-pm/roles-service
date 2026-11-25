package updatetemplate

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
	ID              int
	Version         int64
	Name            string
	PermissionCodes []domain.PermissionCode
}

func (s *Usecase) Do(ctx context.Context, actorID int, input Params) (*domain.RoleTemplate, error) {
	if actorID <= 0 {
		return nil, domain.NewError(domain.CodeInvalidInput, "actor ID must be positive")
	}

	if input.ID <= 0 || input.Version <= 0 {
		return nil, domain.FieldError(domain.CodeInvalidInput, "expected_version", "ID and expected_version must be positive")
	}

	var result *domain.RoleTemplate

	err := s.transactor.InTx(ctx, func(ctx context.Context) error {
		value, err := s.repository.Get(ctx, input.ID)
		if err != nil {
			return err
		}

		if value.Version != input.Version {
			return domain.ErrConflict
		}

		permissions, err := s.catalog.NewSet(value.TargetScopeType, input.PermissionCodes)
		if err != nil {
			return err
		}

		if err = value.Rename(input.Name); err != nil {
			return err
		}

		if err = value.SetPermissions(permissions); err != nil {
			return err
		}

		if value.Version != input.Version {
			value.Version = input.Version + 1
		}

		if err = s.repository.Update(ctx, value, input.Version); err != nil {
			return err
		}

		result = value

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
