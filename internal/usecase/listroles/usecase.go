package listroles

import (
	"context"

	"github.com/hedgehog-pm/roles-service/internal/domain"
	"github.com/hedgehog-pm/roles-service/internal/pagination"
)

type (
	Params struct {
		pagination.Options
		PermissionCode domain.PermissionCode
	}
	Usecase struct {
		reader     reader
		transactor transactor
		catalog    domain.PermissionCatalog
	}
)

func New(reader reader, transactor transactor, catalog domain.PermissionCatalog) *Usecase {
	return &Usecase{reader: reader, transactor: transactor, catalog: catalog}
}

func (s *Usecase) Do(ctx context.Context, scope domain.ScopeRef, options Params) (pagination.Page[*domain.Role], error) {
	if _, err := domain.NewScopeRef(scope.Type(), scope.ID()); err != nil {
		return pagination.Page[*domain.Role]{}, err
	}

	if err := options.Validate(); err != nil {
		return pagination.Page[*domain.Role]{}, err
	}

	if options.PermissionCode != "" {
		if _, err := s.catalog.NewSet(scope.Type(), []domain.PermissionCode{options.PermissionCode}); err != nil {
			return pagination.Page[*domain.Role]{}, err
		}
	}

	var result pagination.Page[*domain.Role]

	err := s.transactor.InReadTx(ctx, func(ctx context.Context) error {
		var err error

		result, err = s.reader.Page(ctx, scope, options)

		return err
	})

	return result, err
}
