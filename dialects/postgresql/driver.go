package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/schema"
)

type PostgreSQLDriver struct{}

func (d PostgreSQLDriver) Name() string {
	return "postgres"
}

func (d PostgreSQLDriver) ConnectionString(pParams ormshift.ConnectionParams) string {
	lHost := pParams.Host
	if lHost == "" {
		lHost = "localhost"
	}
	lPort := pParams.Port
	if lPort == 0 {
		lPort = 5432
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", lHost, lPort, pParams.User, pParams.Password, pParams.Database)
}

func (d PostgreSQLDriver) SQLBuilder() ormshift.SQLBuilder {
	return postgresqlSQLBuilder{}
}

func (d PostgreSQLDriver) DBSchema(pDB *sql.DB) (*schema.DBSchema, error) {
	return schema.NewDBSchema(pDB, tableNamesQuery)
}
