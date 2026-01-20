package sqlite

import (
	"database/sql"

	"github.com/ordershift/ormshift"
)

const tableNamesQuery = `
		SELECT
			name 
		FROM
			sqlite_master
		WHERE
			type = 'table'
		ORDER BY
			name
	`

func DBSchema(pDB *sql.DB) (*ormshift.DBSchema, error) {
	return ormshift.NewDBSchema(pDB, tableNamesQuery)
}
