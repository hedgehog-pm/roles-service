package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hedgehog-pm/roles-service/internal/domain"
	rolesdb "github.com/hedgehog-pm/roles-service/internal/generated/bob"
)

type Templates struct {
	pool    *pgxpool.Pool
	catalog domain.PermissionCatalog
}

func NewTemplates(pool *pgxpool.Pool, catalog domain.PermissionCatalog) *Templates {
	return &Templates{pool: pool, catalog: catalog}
}

func (s *Templates) restore(row rolesdb.GetTemplateRow) (*domain.RoleTemplate, error) {
	scopeType := domain.ScopeType(row.TargetScopeType)

	permissions, err := s.catalog.NewSet(scopeType, permissionCodes(row.Codes))
	if err != nil {
		return nil, fmt.Errorf("restore templates permissions: %w", err)
	}

	return &domain.RoleTemplate{
		ID: row.ID, Name: row.Name, SystemID: row.SystemID, TargetScopeType: scopeType,
		Permissions: permissions, Version: row.Version, CreatedBy: row.CreatedBy,
		CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}, nil
}

func (s *Templates) Get(ctx context.Context, id int) (*domain.RoleTemplate, error) {
	query := rolesdb.GetTemplate(id)

	row, err := queryOne(ctx, s.pool, query, query.Scanner)
	if err != nil {
		return nil, storageError(err)
	}

	return s.restore(row)
}

func (s *Templates) Insert(ctx context.Context, value *domain.RoleTemplate) (*domain.RoleTemplate, error) {
	var saved *domain.RoleTemplate

	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		query := rolesdb.InsertTemplate([]rolesdb.InsertTemplate_Rows{{
			Name: value.Name, TargetScopeType: string(value.TargetScopeType), SystemID: value.SystemID,
			Version: value.Version, CreatedBy: value.CreatedBy, CreatedAt: value.CreatedAt, UpdatedAt: value.UpdatedAt,
		}})

		row, err := queryOne(ctx, tx, query, query.Scanner)
		if err != nil {
			return err
		}

		saved = &domain.RoleTemplate{
			ID: row.ID, Name: row.Name, TargetScopeType: domain.ScopeType(row.TargetScopeType), SystemID: row.SystemID,
			Permissions: value.Permissions, Version: row.Version, CreatedBy: row.CreatedBy, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
		}
		if err := saveTemplatePermissions(ctx, tx, saved); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, storageError(err)
	}

	return saved, nil
}

func (s *Templates) Update(ctx context.Context, value *domain.RoleTemplate, expected int64) error {
	return storageError(pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		affected, err := queryExec(ctx, tx, rolesdb.UpdateTemplate(value.Name, value.Version, value.UpdatedAt, value.ID, expected))
		if err != nil {
			return err
		}

		if affected.RowsAffected() != 1 {
			query := rolesdb.GetTemplateForUpdate(value.ID)
			if _, err := queryOne(ctx, tx, query, query.Scanner); err != nil {
				return err
			}

			return domain.ErrConflict
		}

		if value.Version == expected {
			return nil
		}

		if _, err := queryExec(ctx, tx, rolesdb.ClearTemplatePermissions(value.ID)); err != nil {
			return err
		}

		return saveTemplatePermissions(ctx, tx, value)
	}))
}

func saveTemplatePermissions(ctx context.Context, tx pgx.Tx, value *domain.RoleTemplate) error {
	for _, code := range value.Permissions.Codes() {
		_, err := queryExec(ctx, tx, rolesdb.InsertTemplatePermission([]rolesdb.InsertTemplatePermission_Rows{{
			TemplateID: value.ID, PermissionCode: string(code),
		}}))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Templates) Delete(ctx context.Context, id int, expected int64) error {
	return storageError(pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		query := rolesdb.GetTemplateForUpdate(id)

		row, err := queryOne(ctx, tx, query, query.Scanner)
		if err != nil {
			return err
		}

		value, err := s.restore(row)
		if err != nil {
			return err
		}

		if value.Version != expected {
			return domain.ErrConflict
		}

		_, err = queryExec(ctx, tx, rolesdb.DeleteTemplate(id))

		return err
	}))
}
