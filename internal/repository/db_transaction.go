package repository

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func RunInTransaction(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
	sqlTx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(sqlTx); err != nil {
		_ = sqlTx.Rollback()
		return err
	}
	if err := sqlTx.Commit(); err != nil {
		_ = sqlTx.Rollback()
		return err
	}
	return nil
}
