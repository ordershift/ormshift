package sqlite

import (
	"database/sql"
	"fmt"

	// Blank import to register the SQLite driver
	_ "modernc.org/sqlite"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/schema"
)

type sqliteDriver struct {
	sqlBuilder ormshift.SQLBuilder
}

func Driver() ormshift.DatabaseDriver {
	return &sqliteDriver{
		sqlBuilder: newSQLiteBuilder(),
	}
}

func (d *sqliteDriver) Name() string {
	return "sqlite"
}

func (d *sqliteDriver) ConnectionString(params ormshift.ConnectionParams) string {
	if params.InMemory {
		return ":memory:"
	}
	connectionWithAuth := ""
	if params.User != "" {
		connectionWithAuth += fmt.Sprintf("_auth&_auth_user=%s&", params.User)
		if params.Password != "" {
			connectionWithAuth += fmt.Sprintf("_auth_pass=%s&", params.Password)
		}
	}
	return fmt.Sprintf("file:%s.db?%s_locking=NORMAL", params.Database, connectionWithAuth)
}

func (d *sqliteDriver) SQLBuilder() ormshift.SQLBuilder {
	return d.sqlBuilder
}

func (d *sqliteDriver) DBSchema(db *sql.DB) (*schema.DBSchema, error) {
	return schema.NewDBSchema(db, tableNamesQuery, columnTypesQueryFunc(d.sqlBuilder))
}
