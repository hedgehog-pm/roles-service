package role

import (
	"context"

	"github.com/hedgehog-pm/roles-service/internal/domain"
)

type Usecase struct {
	repository repository
	transactor transactor
}

func New(repository repository, transactor transactor) *Usecase {
	return &Usecase{repository: repository, transactor: transactor}
}

func (s *Usecase) Get(ctx context.Context, id int) (*domain.Role, error) {
	if id <= 0 {
		return nil, domain.FieldError(domain.CodeInvalidInput, "id", "ID must be positive")
	}

	return s.repository.Get(ctx, id)
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
