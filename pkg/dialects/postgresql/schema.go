package postgresql

import (
	"database/sql"

	"github.com/ordershift/ormshift/pkg/core"
)

const tableNamesQuery = `
		SELECT
			table_name
		FROM
			information_schema.tables
		WHERE
			table_type = 'BASE TABLE' AND
			table_schema NOT IN ('pg_catalog', 'information_schema')
		ORDER BY
			table_name
	`

func DBSchema(pDB *sql.DB) (*core.DBSchema, error) {
	return core.NewDBSchema(pDB, tableNamesQuery)
}
