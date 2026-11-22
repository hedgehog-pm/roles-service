package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/hedgehog-pm/roles-service/internal/generated/bob/rolesdb"
	"github.com/hedgehog-pm/roles-service/internal/roles/application"
	"github.com/hedgehog-pm/roles-service/internal/roles/domain"
)

func (s *Roles) Page(ctx context.Context, scope domain.ScopeRef, options application.ListOptions) (application.Page[*domain.Role], error) {
	result := application.Page[*domain.Role]{Items: make([]*domain.Role, 0), Limit: options.Limit, Offset: options.Offset}
	err := pgx.BeginTxFunc(ctx, s.pool, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly}, func(tx pgx.Tx) error {
		count := rolesdb.CountRoles(string(scope.Type()), scope.ID(), options.Search, string(options.PermissionCode))

		total, err := queryOne(ctx, tx, count, count.Scanner)
		if err != nil {
			return err
		}

		result.Total = total
		query := rolesdb.PageRoles(string(scope.Type()), scope.ID(), options.Search, string(options.PermissionCode), options.Sort, options.Direction, options.Limit, options.Offset)

		rows, err := queryAll(ctx, tx, query, query.Scanner)
		if err != nil {
			return err
		}

		for _, row := range rows {
			value, err := s.restore(row)
			if err != nil {
				return err
			}

			result.Items = append(result.Items, value)
		}

		return nil
	})

	return result, storageError(err)
}

func (s *Templates) Page(ctx context.Context, filter application.TemplateFilter, options application.ListOptions) (application.Page[*domain.RoleTemplate], error) {
	result := application.Page[*domain.RoleTemplate]{Items: make([]*domain.RoleTemplate, 0), Limit: options.Limit, Offset: options.Offset}
	err := pgx.BeginTxFunc(ctx, s.pool, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly}, func(tx pgx.Tx) error {
		count := rolesdb.CountTemplates(string(filter.TargetScopeType), filter.SystemID, options.Search)

		total, err := queryOne(ctx, tx, count, count.Scanner)
		if err != nil {
			return err
		}

		result.Total = total
		query := rolesdb.PageTemplates(string(filter.TargetScopeType), filter.SystemID, options.Search, options.Sort, options.Direction, options.Limit, options.Offset)

		rows, err := queryAll(ctx, tx, query, query.Scanner)
		if err != nil {
			return err
		}

		for _, row := range rows {
			value, err := s.restore(row)
			if err != nil {
				return err
			}

			result.Items = append(result.Items, value)
		}

		return nil
	})

	return result, storageError(err)
}

func (s *Assignments) Page(ctx context.Context, filter application.AssignmentFilter) (application.Page[application.AssignmentView], error) {
	result := application.Page[application.AssignmentView]{Items: make([]application.AssignmentView, 0), Limit: filter.Limit, Offset: filter.Offset}
	err := pgx.BeginTxFunc(ctx, s.pool, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly}, func(tx pgx.Tx) error {
		count := rolesdb.CountAssignments(filter.RoleID, filter.UserID, string(filter.ScopeType), filter.ScopeID)

		total, err := queryOne(ctx, tx, count, count.Scanner)
		if err != nil {
			return err
		}

		result.Total = total
		query := rolesdb.PageAssignments(filter.RoleID, filter.UserID, string(filter.ScopeType), filter.ScopeID, filter.Limit, filter.Offset)

		rows, err := queryAll(ctx, tx, query, query.Scanner)
		if err != nil {
			return err
		}

		for _, row := range rows {
			scope, err := domain.NewScopeRef(domain.ScopeType(row.ScopeType), row.ScopeID)
			if err != nil {
				return err
			}

			result.Items = append(result.Items, application.AssignmentView{RoleID: row.RoleID, UserID: row.UserID, RoleName: row.RoleName, Scope: scope, AssignedBy: row.AssignedBy, AssignedAt: row.AssignedAt})
		}

		return nil
	})

	return result, storageError(err)
}
