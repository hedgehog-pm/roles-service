package postgres

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/hedgehog-pm/roles-service/internal/roles/application"
)

func storageError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
		return application.ErrNotFound
	}

	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		switch pgErr.Code {
		case "23505", "40001", "40P01":
			return application.ErrConflict
		case "23503":
			return application.ErrNotFound
		}
	}

	return err
}
