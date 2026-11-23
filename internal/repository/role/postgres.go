package role

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hedgehog-pm/roles-service/internal/domain"
	rolesdb "github.com/hedgehog-pm/roles-service/internal/generated/bob/role"
	"github.com/hedgehog-pm/roles-service/internal/pagination"
	"github.com/hedgehog-pm/roles-service/internal/storage/postgres"
	"github.com/hedgehog-pm/roles-service/internal/usecase/listroles"
)

type Repository struct {
	pool    *pgxpool.Pool
	catalog domain.PermissionCatalog
}

func New(pool *pgxpool.Pool, catalog domain.PermissionCatalog) *Repository {
	return &Repository{pool: pool, catalog: catalog}
}

func (s *Repository) Get(ctx context.Context, id int) (*domain.Role, error) {
	query := rolesdb.GetRole(id)

	row, err := postgres.QueryOne(ctx, postgres.Executor(ctx, s.pool), query, query.Scanner)
	if err != nil {
		return nil, postgres.MapError(err)
	}

	return s.restore(row)
}

func (s *Repository) Insert(ctx context.Context, value *domain.Role) (*domain.Role, error) {
	tx := postgres.Executor(ctx, s.pool)
	query := rolesdb.InsertRole([]rolesdb.InsertRole_Rows{{
		Name: value.Name, ScopeType: string(value.Scope.Type()), ScopeID: value.Scope.ID(),
		Version: value.Version, CreatedBy: value.CreatedBy, CreatedAt: value.CreatedAt, UpdatedAt: value.UpdatedAt,
	}})

	row, err := postgres.QueryOne(ctx, tx, query, query.Scanner)
	if err != nil {
		return nil, postgres.MapError(err)
	}

	saved := &domain.Role{
		ID: row.ID, Name: row.Name, Scope: value.Scope,
		Permissions: value.Permissions, Version: row.Version, CreatedBy: row.CreatedBy, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}
	if err := saveRolePermissions(ctx, tx, saved); err != nil {
		return nil, postgres.MapError(err)
	}

	if err := roleEvent(ctx, tx, "RoleCreated", saved); err != nil {
		return nil, postgres.MapError(err)
	}

	return saved, nil
}

func (s *Repository) Update(ctx context.Context, value *domain.Role, expected int64) error {
	tx := postgres.Executor(ctx, s.pool)

	affected, err := postgres.QueryExec(ctx, tx, rolesdb.UpdateRole(value.Name, value.Version, value.UpdatedAt, value.ID, expected))
	if err != nil {
		return postgres.MapError(err)
	}

	if affected.RowsAffected() != 1 {
		query := rolesdb.GetRoleForUpdate(value.ID)
		if _, err := postgres.QueryOne(ctx, tx, query, query.Scanner); err != nil {
			return postgres.MapError(err)
		}

		return domain.ErrConflict
	}

	if value.Version == expected {
		return nil
	}

	if _, err := postgres.QueryExec(ctx, tx, rolesdb.ClearRolePermissions(value.ID)); err != nil {
		return postgres.MapError(err)
	}

	if err := saveRolePermissions(ctx, tx, value); err != nil {
		return postgres.MapError(err)
	}

	return postgres.MapError(roleEvent(ctx, tx, "RoleUpdated", value))
}

func saveRolePermissions(ctx context.Context, tx executor, value *domain.Role) error {
	for _, code := range value.Permissions.Codes() {
		_, err := postgres.QueryExec(ctx, tx, rolesdb.InsertRolePermission([]rolesdb.InsertRolePermission_Rows{{
			RoleID: value.ID, PermissionCode: string(code),
		}}))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Repository) Delete(ctx context.Context, id int, expected int64) error {
	tx := postgres.Executor(ctx, s.pool)
	query := rolesdb.GetRoleForUpdate(id)

	row, err := postgres.QueryOne(ctx, tx, query, query.Scanner)
	if err != nil {
		return postgres.MapError(err)
	}

	value, err := s.restore(row)
	if err != nil {
		return postgres.MapError(err)
	}

	if value.Version != expected {
		return domain.ErrConflict
	}

	if _, err := postgres.QueryExec(ctx, tx, rolesdb.DeleteRole(id)); err != nil {
		return postgres.MapError(err)
	}

	return postgres.MapError(roleEvent(ctx, tx, "RoleDeleted", value))
}

func (s *Repository) Page(ctx context.Context, scope domain.ScopeRef, options listroles.Params) (pagination.Page[*domain.Role], error) {
	result := pagination.Page[*domain.Role]{Items: make([]*domain.Role, 0), Limit: options.Limit, Offset: options.Offset}

	tx := postgres.Executor(ctx, s.pool)
	count := rolesdb.CountRoles(string(scope.Type()), scope.ID(), options.Search, string(options.PermissionCode))

	total, err := postgres.QueryOne(ctx, tx, count, count.Scanner)
	if err != nil {
		return result, postgres.MapError(err)
	}

	result.Total = total
	query := rolesdb.PageRoles(string(scope.Type()), scope.ID(), options.Search, string(options.PermissionCode), options.Sort, options.Direction, options.Limit, options.Offset)

	rows, err := postgres.QueryAll(ctx, tx, query, query.Scanner)
	if err != nil {
		return result, postgres.MapError(err)
	}

	for _, row := range rows {
		value, err := s.restore(row)
		if err != nil {
			return result, postgres.MapError(err)
		}

		result.Items = append(result.Items, value)
	}

	return result, nil
}
