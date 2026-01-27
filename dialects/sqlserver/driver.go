package sqlserver

import (
	"database/sql"
	"fmt"

	// Blank import to register the SQL Server driver
	_ "github.com/microsoft/go-mssqldb"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/schema"
)

type sqlserverDriver struct{}

func Driver() ormshift.DatabaseDriver {
	return sqlserverDriver{}
}

func (d sqlserverDriver) Name() string {
	return "sqlserver"
}

func (d sqlserverDriver) ConnectionString(pParams ormshift.ConnectionParams) string {
	lHostInstanceAndPort := pParams.Host
	if pParams.Instance != "" {
		lHostInstanceAndPort = fmt.Sprintf("%s\\%s", pParams.Host, pParams.Instance)
	}
	if pParams.Port > 0 {
		lHostInstanceAndPort += fmt.Sprintf(";port=%d", pParams.Port)
	}
	return fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", lHostInstanceAndPort, pParams.User, pParams.Password, pParams.Database)
}

func (d sqlserverDriver) SQLBuilder() ormshift.SQLBuilder {
	return newSQLServerBuilder()
}

func (d sqlserverDriver) DBSchema(pDB *sql.DB) (*schema.DBSchema, error) {
	return schema.NewDBSchema(pDB, tableNamesQuery, columnTypesQueryFunc)
}
