package listroleassignments

import (
	"context"

	"github.com/hedgehog-pm/roles-service/internal/domain"
	"github.com/hedgehog-pm/roles-service/internal/pagination"
)

type Usecase struct {
	reader     reader
	transactor transactor
}

func New(reader reader, transactor transactor) *Usecase {
	return &Usecase{reader: reader, transactor: transactor}
}

func (s *Usecase) Do(ctx context.Context, filter Params) (pagination.Page[Assignment], error) {
	if filter.RoleID < 0 || filter.UserID < 0 || (filter.RoleID == 0 && (filter.ScopeID <= 0 || filter.ScopeType == "")) {
		return pagination.Page[Assignment]{}, domain.NewError(domain.CodeInvalidInput, "role or user ID is required")
	}

	if filter.ScopeID < 0 || (filter.ScopeID > 0 && filter.ScopeType == "") {
		return pagination.Page[Assignment]{}, domain.FieldError(domain.CodeInvalidInput, "scope_id", "scope_id requires scope_type")
	}

	if filter.ScopeType != "" {
		if _, err := domain.NewScopeRef(filter.ScopeType, 1); err != nil {
			return pagination.Page[Assignment]{}, err
		}
	}

	if err := filter.Validate(); err != nil {
		return pagination.Page[Assignment]{}, err
	}

	var result pagination.Page[Assignment]

	err := s.transactor.InReadTx(ctx, func(ctx context.Context) error {
		var err error

		result, err = s.reader.Page(ctx, filter)

		return err
	})

	return result, err
}
