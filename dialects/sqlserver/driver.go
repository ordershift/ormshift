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

func (d *sqlserverDriver) ConnectionString(pParams ormshift.ConnectionParams) string {
	hostInstanceAndPort := pParams.Host
	if pParams.Instance != "" {
		hostInstanceAndPort = fmt.Sprintf("%s\\%s", pParams.Host, pParams.Instance)
	}
	if pParams.Port > 0 {
		hostInstanceAndPort += fmt.Sprintf(";port=%d", pParams.Port)
	}
	return fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", hostInstanceAndPort, pParams.User, pParams.Password, pParams.Database)
}

func (d *sqlserverDriver) SQLBuilder() ormshift.SQLBuilder {
	return d.sqlBuilder
}

func (d *sqlserverDriver) DBSchema(pDB *sql.DB) (*schema.DBSchema, error) {
	return schema.NewDBSchema(pDB, tableNamesQuery, columnTypesQueryFunc(d.sqlBuilder))
}
