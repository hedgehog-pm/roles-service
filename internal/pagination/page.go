package pagination

import "github.com/hedgehog-pm/roles-service/internal/domain"

type Pagination struct {
	Limit  int
	Offset int
}

func (p *Pagination) Validate() error {
	if p.Limit == 0 {
		p.Limit = 20
	}

	if p.Limit < 1 || p.Limit > 100 || p.Offset < 0 {
		return domain.FieldError(domain.CodeInvalidInput, "pagination", "limit must be 1 to 100 and offset nonnegative")
	}

	return nil
}

type Options struct {
	Pagination
	Search    string
	Sort      string
	Direction string
}

func (o *Options) Validate() error {
	if err := o.Pagination.Validate(); err != nil {
		return err
	}

	if o.Sort == "" {
		o.Sort = "created_at"
	}

	if o.Direction == "" {
		o.Direction = "desc"
	}

	switch o.Sort {
	case "name", "created_at", "updated_at":
	default:
		return domain.FieldError(domain.CodeInvalidInput, "sort", "unsupported sort field")
	}

	if o.Direction != "asc" && o.Direction != "desc" {
		return domain.FieldError(domain.CodeInvalidInput, "direction", "direction must be asc or desc")
	}

	return nil
}

type Page[T any] struct {
	Items  []T
	Total  int64
	Limit  int
	Offset int
}
