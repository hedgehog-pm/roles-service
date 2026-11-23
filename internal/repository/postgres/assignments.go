package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hedgehog-pm/roles-service/internal/domain"
	rolesdb "github.com/hedgehog-pm/roles-service/internal/generated/bob"
)

type Assignments struct{ pool *pgxpool.Pool }

func NewAssignments(pool *pgxpool.Pool) *Assignments { return &Assignments{pool: pool} }

type assignmentPayload struct {
	RoleID int `json:"role_id"`
	UserID int `json:"user_id"`
}

func (s *Assignments) Assign(ctx context.Context, a domain.RoleAssignment) error {
	_, err := s.Ensure(ctx, a)
	return err
}

func (s *Assignments) Ensure(ctx context.Context, a domain.RoleAssignment) (domain.AssignmentResult, error) {
	result := domain.AssignmentResult{}

	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		lock := rolesdb.LockAssignmentRole(a.RoleID())
		if _, err := queryOne(ctx, tx, lock, lock.Scanner); err != nil {
			return err
		}

		query := rolesdb.GetAssignment(a.RoleID(), a.UserID())

		row, err := queryOne(ctx, tx, query, query.Scanner)
		if err == nil {
			value, restoreErr := domain.RestoreRoleAssignment(row.RoleID, row.UserID, row.AssignedBy, row.AssignedAt)
			result.Assignment = value

			return restoreErr
		}

		if !errors.Is(storageError(err), domain.ErrNotFound) {
			return err
		}

		if _, err := queryExec(ctx, tx, rolesdb.InsertAssignment([]rolesdb.InsertAssignment_Rows{{RoleID: a.RoleID(), UserID: a.UserID(), AssignedBy: a.AssignedBy(), AssignedAt: a.AssignedAt()}})); err != nil {
			return err
		}

		row, err = queryOne(ctx, tx, query, query.Scanner)
		if err != nil {
			return err
		}

		result.Assignment, err = domain.RestoreRoleAssignment(row.RoleID, row.UserID, row.AssignedBy, row.AssignedAt)
		if err != nil {
			return err
		}

		result.Created = true

		return event(ctx, tx, "RoleAssigned", a.RoleID(), assignmentPayload{RoleID: a.RoleID(), UserID: a.UserID()})
	})
	if err != nil {
		return domain.AssignmentResult{}, storageError(err)
	}

	return result, nil
}

func (s *Assignments) Revoke(ctx context.Context, roleID, userID int) error {
	return storageError(pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		query := rolesdb.LockAssignmentRole(roleID)
		if _, err := queryOne(ctx, tx, query, query.Scanner); err != nil {
			return err
		}

		affected, err := queryExec(ctx, tx, rolesdb.DeleteAssignment(roleID, userID))
		if err != nil {
			return err
		}

		if affected.RowsAffected() == 0 {
			return nil
		}

		return event(ctx, tx, "RoleRevoked", roleID, assignmentPayload{RoleID: roleID, UserID: userID})
	}))
}

func (s *Assignments) ListAssignments(ctx context.Context, roleID int) ([]domain.RoleAssignment, error) {
	query := rolesdb.ListAssignments(roleID)

	rows, err := queryAll(ctx, s.pool, query, query.Scanner)
	if err != nil {
		return nil, storageError(err)
	}

	var result []domain.RoleAssignment

	for _, row := range rows {
		a, err := domain.RestoreRoleAssignment(roleID, row.UserID, row.AssignedBy, row.AssignedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, a)
	}

	return result, nil
}
