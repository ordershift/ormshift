package sqlserver

import (
	"database/sql"
	"fmt"

	// Blank import to register the SQL Server driver
	_ "github.com/microsoft/go-mssqldb"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/schema"
)

type sqlserverDriver struct {
	sqlBuilder ormshift.SQLBuilder
}

func Driver() ormshift.DatabaseDriver {
	return &sqlserverDriver{
		sqlBuilder: newSQLServerBuilder(),
	}
}

func (d *sqlserverDriver) Name() string {
	return "sqlserver"
}

func (d *sqlserverDriver) ConnectionString(params ormshift.ConnectionParams) string {
	hostInstanceAndPort := params.Host
	if params.Instance != "" {
		hostInstanceAndPort = fmt.Sprintf("%s\\%s", params.Host, params.Instance)
	}
	if params.Port > 0 {
		hostInstanceAndPort += fmt.Sprintf(";port=%d", params.Port)
	}
	return fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", hostInstanceAndPort, params.User, params.Password, params.Database)
}

func (d *sqlserverDriver) SQLBuilder() ormshift.SQLBuilder {
	return d.sqlBuilder
}

func (d *sqlserverDriver) DBSchema(db *sql.DB) (*schema.DBSchema, error) {
	return schema.NewDBSchema(db, tableNamesQuery, columnTypesQueryFunc(d.sqlBuilder))
}
