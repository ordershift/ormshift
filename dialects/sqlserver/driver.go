package sqlserver

import (
	"database/sql"
	"fmt"

	_ "github.com/microsoft/go-mssqldb"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/schema"
)

type SQLServerDriver struct{}

func (d SQLServerDriver) Name() string {
	return "sqlserver"
}

func (d SQLServerDriver) ConnectionString(pParams ormshift.ConnectionParams) string {
	lHostInstanceAndPort := pParams.Host
	if pParams.Instance != "" {
		lHostInstanceAndPort = fmt.Sprintf("%s\\%s", pParams.Host, pParams.Instance)
	}
	if pParams.Port > 0 {
		lHostInstanceAndPort += fmt.Sprintf(";port=%d", pParams.Port)
	}
	return fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", lHostInstanceAndPort, pParams.User, pParams.Password, pParams.Database)
}

func (d SQLServerDriver) SQLBuilder() ormshift.SQLBuilder {
	return sqlserverSQLBuilder{}
}

func (d SQLServerDriver) DBSchema(pDB *sql.DB) (*schema.DBSchema, error) {
	return schema.NewDBSchema(pDB, tableNamesQuery)
}
