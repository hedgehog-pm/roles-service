package roleassignment

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hedgehog-pm/roles-service/internal/domain"
	rolesdb "github.com/hedgehog-pm/roles-service/internal/generated/bob/roleassignment"
	"github.com/hedgehog-pm/roles-service/internal/storage/postgres"
	"github.com/hedgehog-pm/roles-service/internal/usecase/roleassignment"
)

type Repository struct{ pool *pgxpool.Pool }

func New(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

type assignmentPayload struct {
	RoleID int `json:"role_id"`
	UserID int `json:"user_id"`
}

func (s *Repository) Assign(ctx context.Context, a domain.RoleAssignment) error {
	_, err := s.Ensure(ctx, a)
	return err
}

func (s *Repository) Ensure(ctx context.Context, a domain.RoleAssignment) (roleassignment.Result, error) {
	var result roleassignment.Result

	tx := postgres.Executor(ctx, s.pool)

	lock := rolesdb.LockAssignmentRole(a.RoleID())
	if _, err := postgres.QueryOne(ctx, tx, lock, lock.Scanner); err != nil {
		return roleassignment.Result{}, postgres.MapError(err)
	}

	query := rolesdb.GetAssignment(a.RoleID(), a.UserID())

	row, err := postgres.QueryOne(ctx, tx, query, query.Scanner)
	if err == nil {
		value, restoreErr := domain.RestoreRoleAssignment(row.RoleID, row.UserID, row.AssignedBy, row.AssignedAt)
		result.Assignment = value

		if restoreErr != nil {
			return roleassignment.Result{}, postgres.MapError(restoreErr)
		}

		return result, nil
	}

	if !errors.Is(postgres.MapError(err), domain.ErrNotFound) {
		return roleassignment.Result{}, postgres.MapError(err)
	}

	if _, err := postgres.QueryExec(ctx, tx, rolesdb.InsertAssignment([]rolesdb.InsertAssignment_Rows{{RoleID: a.RoleID(), UserID: a.UserID(), AssignedBy: a.AssignedBy(), AssignedAt: a.AssignedAt()}})); err != nil {
		return roleassignment.Result{}, postgres.MapError(err)
	}

	row, err = postgres.QueryOne(ctx, tx, query, query.Scanner)
	if err != nil {
		return roleassignment.Result{}, postgres.MapError(err)
	}

	result.Assignment, err = domain.RestoreRoleAssignment(row.RoleID, row.UserID, row.AssignedBy, row.AssignedAt)
	if err != nil {
		return roleassignment.Result{}, postgres.MapError(err)
	}

	result.Created = true

	if err := event(ctx, tx, "RoleAssigned", a.RoleID(), assignmentPayload{RoleID: a.RoleID(), UserID: a.UserID()}); err != nil {
		return roleassignment.Result{}, postgres.MapError(err)
	}

	return result, nil
}

func (s *Repository) Revoke(ctx context.Context, roleID, userID int) error {
	tx := postgres.Executor(ctx, s.pool)

	query := rolesdb.LockAssignmentRole(roleID)
	if _, err := postgres.QueryOne(ctx, tx, query, query.Scanner); err != nil {
		return postgres.MapError(err)
	}

	affected, err := postgres.QueryExec(ctx, tx, rolesdb.DeleteAssignment(roleID, userID))
	if err != nil {
		return postgres.MapError(err)
	}

	if affected.RowsAffected() == 0 {
		return nil
	}

	return postgres.MapError(event(ctx, tx, "RoleRevoked", roleID, assignmentPayload{RoleID: roleID, UserID: userID}))
}

func (s *Repository) ListAssignments(ctx context.Context, roleID int) ([]domain.RoleAssignment, error) {
	query := rolesdb.ListAssignments(roleID)

	rows, err := postgres.QueryAll(ctx, postgres.Executor(ctx, s.pool), query, query.Scanner)
	if err != nil {
		return nil, postgres.MapError(err)
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
