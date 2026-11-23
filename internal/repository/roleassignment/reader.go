package roleassignment

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hedgehog-pm/roles-service/internal/domain"
	rolesdb "github.com/hedgehog-pm/roles-service/internal/generated/bob/roleassignment"
	"github.com/hedgehog-pm/roles-service/internal/pagination"
	"github.com/hedgehog-pm/roles-service/internal/storage/postgres"
	"github.com/hedgehog-pm/roles-service/internal/usecase/listroleassignments"
)

type Reader struct{ pool *pgxpool.Pool }

func NewReader(pool *pgxpool.Pool) *Reader { return &Reader{pool: pool} }
func (s *Reader) Page(ctx context.Context, filter listroleassignments.Params) (pagination.Page[listroleassignments.Assignment], error) {
	result := pagination.Page[listroleassignments.Assignment]{Items: make([]listroleassignments.Assignment, 0), Limit: filter.Limit, Offset: filter.Offset}

	tx := postgres.Executor(ctx, s.pool)
	count := rolesdb.CountAssignments(filter.RoleID, filter.UserID, string(filter.ScopeType), filter.ScopeID)

	total, err := postgres.QueryOne(ctx, tx, count, count.Scanner)
	if err != nil {
		return result, postgres.MapError(err)
	}

	result.Total = total
	query := rolesdb.PageAssignments(filter.RoleID, filter.UserID, string(filter.ScopeType), filter.ScopeID, filter.Limit, filter.Offset)

	rows, err := postgres.QueryAll(ctx, tx, query, query.Scanner)
	if err != nil {
		return result, postgres.MapError(err)
	}

	for _, row := range rows {
		scope, err := domain.NewScopeRef(domain.ScopeType(row.ScopeType), row.ScopeID)
		if err != nil {
			return result, postgres.MapError(err)
		}

		result.Items = append(result.Items, listroleassignments.Assignment{RoleID: row.RoleID, UserID: row.UserID, RoleName: row.RoleName, Scope: scope, AssignedBy: row.AssignedBy, AssignedAt: row.AssignedAt})
	}

	return result, nil
}
