package core

import (
	"context"
	"database/sql"
)

// SQLExecutor defines the interface for executing SQL commands, relying on database/sql package.
type SQLExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}
