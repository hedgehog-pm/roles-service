package roletemplate

import (
	"context"

	"github.com/hedgehog-pm/roles-service/internal/domain"
	"github.com/hedgehog-pm/roles-service/internal/pagination"
)

type Usecase struct {
	repository repository
	transactor transactor
}

func New(repository repository, transactor transactor) *Usecase {
	return &Usecase{repository: repository, transactor: transactor}
}

func (s *Usecase) Get(ctx context.Context, id int) (*domain.RoleTemplate, error) {
	if id <= 0 {
		return nil, domain.FieldError(domain.CodeInvalidInput, "id", "ID must be positive")
	}

	return s.repository.Get(ctx, id)
}

func (s *Usecase) Page(ctx context.Context, filter Filter, options pagination.Options) (pagination.Page[*domain.RoleTemplate], error) {
	if _, err := domain.NewScopeRef(filter.TargetScopeType, filter.SystemID); err != nil {
		return pagination.Page[*domain.RoleTemplate]{}, err
	}

	if err := options.Validate(); err != nil {
		return pagination.Page[*domain.RoleTemplate]{}, err
	}

	var result pagination.Page[*domain.RoleTemplate]

	err := s.transactor.InReadTx(ctx, func(ctx context.Context) error {
		var err error

		result, err = s.repository.Page(ctx, filter, options)

		return err
	})

	return result, err
}

func (s *Usecase) Delete(ctx context.Context, actorID, id int, version int64) error {
	if actorID <= 0 {
		return domain.NewError(domain.CodeInvalidInput, "actor ID must be positive")
	}

	if id <= 0 {
		return domain.FieldError(domain.CodeInvalidInput, "id", "ID must be positive")
	}

	if version < 1 {
		return domain.NewError(domain.CodeInvalidInput, "version must be positive")
	}

	return s.transactor.InTx(ctx, func(ctx context.Context) error { return s.repository.Delete(ctx, id, version) })
}

func (s *Usecase) CheckSystem(ctx context.Context, id, systemID int) error {
	value, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	if value.SystemID != systemID {
		return domain.ErrNotFound
	}

	return nil
}
