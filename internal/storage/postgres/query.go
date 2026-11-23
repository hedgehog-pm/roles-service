package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/scan"
	"github.com/stephenafamo/scan/pgxscan"
)

func QueryOne[T any](ctx context.Context, db pgxscan.Queryer, query bob.Query, mapper scan.Mapper[T]) (T, error) {
	sql, args, err := bob.Build(ctx, query)
	if err != nil {
		var zero T
		return zero, err
	}

	return pgxscan.One(ctx, db, mapper, sql, args...)
}

func QueryAll[T any](ctx context.Context, db pgxscan.Queryer, query bob.Query, mapper scan.Mapper[T]) ([]T, error) {
	sql, args, err := bob.Build(ctx, query)
	if err != nil {
		return nil, err
	}

	return pgxscan.All(ctx, db, mapper, sql, args...)
}

func QueryExec(ctx context.Context, tx executor, query bob.Query) (pgconn.CommandTag, error) {
	sql, args, err := bob.Build(ctx, query)
	if err != nil {
		return pgconn.CommandTag{}, err
	}

	return tx.Exec(ctx, sql, args...)
}
