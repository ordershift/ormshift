package sqlite

import (
	"database/sql"

	"github.com/ordershift/ormshift/pkg/core"
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

func DBSchema(pDB *sql.DB) (*core.DBSchema, error) {
	return core.NewDBSchema(pDB, tableNamesQuery)
}
