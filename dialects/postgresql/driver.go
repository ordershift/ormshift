package postgresql

import (
	"database/sql"
	"fmt"

	// Blank import to register the PostgreSQL driver
	_ "github.com/lib/pq"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/schema"
)

type postgresqlDriver struct {
	sqlBuilder ormshift.SQLBuilder
}

func Driver() ormshift.DatabaseDriver {
	return &postgresqlDriver{
		sqlBuilder: newPostgreSQLBuilder(),
	}
}

func (d *postgresqlDriver) Name() string {
	return "postgres"
}

func (d *postgresqlDriver) ConnectionString(params ormshift.ConnectionParams) string {
	host := params.Host
	if host == "" {
		host = "localhost"
	}
	port := params.Port
	if port == 0 {
		port = 5432
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, params.User, params.Password, params.Database)
}

func (d *postgresqlDriver) SQLBuilder() ormshift.SQLBuilder {
	return d.sqlBuilder
}

func (d *postgresqlDriver) DBSchema(db *sql.DB) (*schema.DBSchema, error) {
	return schema.NewDBSchema(db, tableNamesQuery, columnTypesQueryFunc(d.sqlBuilder))
}
