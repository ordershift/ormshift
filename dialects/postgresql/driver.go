package postgresql

import (
	"database/sql"
	"fmt"

	// Blank import to register the PostgreSQL driver
	_ "github.com/lib/pq"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/schema"
)

type postgresqlDriver struct{}

func Driver() ormshift.DatabaseDriver {
	return postgresqlDriver{}
}

func (d postgresqlDriver) Name() string {
	return "postgres"
}

func (d postgresqlDriver) ConnectionString(pParams ormshift.ConnectionParams) string {
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

func (d postgresqlDriver) SQLBuilder() ormshift.SQLBuilder {
	return newPostgreSQLBuilder()
}

func (d postgresqlDriver) DBSchema(pDB *sql.DB) (*schema.DBSchema, error) {
	return schema.NewDBSchema(pDB, tableNamesQuery, columnTypesQueryFunc)
}
