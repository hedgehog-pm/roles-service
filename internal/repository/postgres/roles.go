package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hedgehog-pm/roles-service/internal/domain"
	rolesdb "github.com/hedgehog-pm/roles-service/internal/generated/bob"
)

type Roles struct {
	pool    *pgxpool.Pool
	catalog domain.PermissionCatalog
}

func NewRoles(pool *pgxpool.Pool, catalog domain.PermissionCatalog) *Roles {
	return &Roles{pool: pool, catalog: catalog}
}

func (s *Roles) restore(row rolesdb.GetRoleRow) (*domain.Role, error) {
	scopeType := domain.ScopeType(row.ScopeType)

	permissions, err := s.catalog.NewSet(scopeType, permissionCodes(row.Codes))
	if err != nil {
		return nil, fmt.Errorf("restore roles permissions: %w", err)
	}

	scope, err := domain.NewScopeRef(scopeType, row.ScopeID)
	if err != nil {
		return nil, err
	}

	return &domain.Role{
		ID: row.ID, Name: row.Name, Scope: scope,
		Permissions: permissions, Version: row.Version, CreatedBy: row.CreatedBy,
		CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}, nil
}

func (s *Roles) Get(ctx context.Context, id int) (*domain.Role, error) {
	query := rolesdb.GetRole(id)

	row, err := queryOne(ctx, s.pool, query, query.Scanner)
	if err != nil {
		return nil, storageError(err)
	}

	return s.restore(row)
}

func (s *Roles) Insert(ctx context.Context, value *domain.Role) (*domain.Role, error) {
	var saved *domain.Role

	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		query := rolesdb.InsertRole([]rolesdb.InsertRole_Rows{{
			Name: value.Name, ScopeType: string(value.Scope.Type()), ScopeID: value.Scope.ID(),
			Version: value.Version, CreatedBy: value.CreatedBy, CreatedAt: value.CreatedAt, UpdatedAt: value.UpdatedAt,
		}})

		row, err := queryOne(ctx, tx, query, query.Scanner)
		if err != nil {
			return err
		}

		saved = &domain.Role{
			ID: row.ID, Name: row.Name, Scope: value.Scope,
			Permissions: value.Permissions, Version: row.Version, CreatedBy: row.CreatedBy, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
		}
		if err := saveRolePermissions(ctx, tx, saved); err != nil {
			return err
		}

		return roleEvent(ctx, tx, "RoleCreated", saved)
	})
	if err != nil {
		return nil, storageError(err)
	}

	return saved, nil
}

func (s *Roles) Update(ctx context.Context, value *domain.Role, expected int64) error {
	return storageError(pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		affected, err := queryExec(ctx, tx, rolesdb.UpdateRole(value.Name, value.Version, value.UpdatedAt, value.ID, expected))
		if err != nil {
			return err
		}

		if affected.RowsAffected() != 1 {
			query := rolesdb.GetRoleForUpdate(value.ID)
			if _, err := queryOne(ctx, tx, query, query.Scanner); err != nil {
				return err
			}

			return domain.ErrConflict
		}

		if value.Version == expected {
			return nil
		}

		if _, err := queryExec(ctx, tx, rolesdb.ClearRolePermissions(value.ID)); err != nil {
			return err
		}

		if err := saveRolePermissions(ctx, tx, value); err != nil {
			return err
		}

		return roleEvent(ctx, tx, "RoleUpdated", value)
	}))
}

func saveRolePermissions(ctx context.Context, tx pgx.Tx, value *domain.Role) error {
	for _, code := range value.Permissions.Codes() {
		_, err := queryExec(ctx, tx, rolesdb.InsertRolePermission([]rolesdb.InsertRolePermission_Rows{{
			RoleID: value.ID, PermissionCode: string(code),
		}}))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Roles) Delete(ctx context.Context, id int, expected int64) error {
	return storageError(pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		query := rolesdb.GetRoleForUpdate(id)

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

		if _, err := queryExec(ctx, tx, rolesdb.DeleteRole(id)); err != nil {
			return err
		}

		return roleEvent(ctx, tx, "RoleDeleted", value)
	}))
}
