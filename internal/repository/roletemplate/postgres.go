package roletemplate

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hedgehog-pm/roles-service/internal/domain"
	rolesdb "github.com/hedgehog-pm/roles-service/internal/generated/bob/roletemplate"
	"github.com/hedgehog-pm/roles-service/internal/pagination"
	"github.com/hedgehog-pm/roles-service/internal/storage/postgres"
	"github.com/hedgehog-pm/roles-service/internal/usecase/roletemplate"
)

type Repository struct {
	pool    *pgxpool.Pool
	catalog domain.PermissionCatalog
}

func New(pool *pgxpool.Pool, catalog domain.PermissionCatalog) *Repository {
	return &Repository{pool: pool, catalog: catalog}
}

func (s *Repository) Get(ctx context.Context, id int) (*domain.RoleTemplate, error) {
	query := rolesdb.GetTemplate(id)

	row, err := postgres.QueryOne(ctx, postgres.Executor(ctx, s.pool), query, query.Scanner)
	if err != nil {
		return nil, postgres.MapError(err)
	}

	return s.restore(row)
}

func (s *Repository) Insert(ctx context.Context, value *domain.RoleTemplate) (*domain.RoleTemplate, error) {
	tx := postgres.Executor(ctx, s.pool)
	query := rolesdb.InsertTemplate([]rolesdb.InsertTemplate_Rows{{
		Name: value.Name, TargetScopeType: string(value.TargetScopeType), SystemID: value.SystemID,
		Version: value.Version, CreatedBy: value.CreatedBy, CreatedAt: value.CreatedAt, UpdatedAt: value.UpdatedAt,
	}})

	row, err := postgres.QueryOne(ctx, tx, query, query.Scanner)
	if err != nil {
		return nil, postgres.MapError(err)
	}

	saved := &domain.RoleTemplate{
		ID: row.ID, Name: row.Name, TargetScopeType: domain.ScopeType(row.TargetScopeType), SystemID: row.SystemID,
		Permissions: value.Permissions, Version: row.Version, CreatedBy: row.CreatedBy, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}
	if err := saveTemplatePermissions(ctx, tx, saved); err != nil {
		return nil, postgres.MapError(err)
	}

	return saved, nil
}

func (s *Repository) Update(ctx context.Context, value *domain.RoleTemplate, expected int64) error {
	tx := postgres.Executor(ctx, s.pool)

	affected, err := postgres.QueryExec(ctx, tx, rolesdb.UpdateTemplate(value.Name, value.Version, value.UpdatedAt, value.ID, expected))
	if err != nil {
		return postgres.MapError(err)
	}

	if affected.RowsAffected() != 1 {
		query := rolesdb.GetTemplateForUpdate(value.ID)
		if _, err := postgres.QueryOne(ctx, tx, query, query.Scanner); err != nil {
			return postgres.MapError(err)
		}

		return domain.ErrConflict
	}

	if value.Version == expected {
		return nil
	}

	if _, err := postgres.QueryExec(ctx, tx, rolesdb.ClearTemplatePermissions(value.ID)); err != nil {
		return postgres.MapError(err)
	}

	return postgres.MapError(saveTemplatePermissions(ctx, tx, value))
}

func saveTemplatePermissions(ctx context.Context, tx executor, value *domain.RoleTemplate) error {
	for _, code := range value.Permissions.Codes() {
		_, err := postgres.QueryExec(ctx, tx, rolesdb.InsertTemplatePermission([]rolesdb.InsertTemplatePermission_Rows{{
			TemplateID: value.ID, PermissionCode: string(code),
		}}))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Repository) Delete(ctx context.Context, id int, expected int64) error {
	tx := postgres.Executor(ctx, s.pool)
	query := rolesdb.GetTemplateForUpdate(id)

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

	_, err = postgres.QueryExec(ctx, tx, rolesdb.DeleteTemplate(id))

	return postgres.MapError(err)
}

func (s *Repository) Page(ctx context.Context, filter roletemplate.Filter, options pagination.Options) (pagination.Page[*domain.RoleTemplate], error) {
	result := pagination.Page[*domain.RoleTemplate]{Items: make([]*domain.RoleTemplate, 0), Limit: options.Limit, Offset: options.Offset}

	tx := postgres.Executor(ctx, s.pool)
	count := rolesdb.CountTemplates(string(filter.TargetScopeType), filter.SystemID, options.Search)

	total, err := postgres.QueryOne(ctx, tx, count, count.Scanner)
	if err != nil {
		return result, postgres.MapError(err)
	}

	result.Total = total
	query := rolesdb.PageTemplates(string(filter.TargetScopeType), filter.SystemID, options.Search, options.Sort, options.Direction, options.Limit, options.Offset)

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
