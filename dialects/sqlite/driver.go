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

func (d *sqliteDriver) ConnectionString(pParams ormshift.ConnectionParams) string {
	if pParams.InMemory {
		return ":memory:"
	}
	connectionWithAuth := ""
	if pParams.User != "" {
		connectionWithAuth += fmt.Sprintf("_auth&_auth_user=%s&", pParams.User)
		if pParams.Password != "" {
			connectionWithAuth += fmt.Sprintf("_auth_pass=%s&", pParams.Password)
		}
	}
	return fmt.Sprintf("file:%s.db?%s_locking=NORMAL", pParams.Database, connectionWithAuth)
}

func (d *sqliteDriver) SQLBuilder() ormshift.SQLBuilder {
	return d.sqlBuilder
}

func (d *sqliteDriver) DBSchema(pDB *sql.DB) (*schema.DBSchema, error) {
	return schema.NewDBSchema(pDB, tableNamesQuery, columnTypesQueryFunc(d.sqlBuilder))
}
