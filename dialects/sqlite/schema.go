package sqlite

import (
	"fmt"

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

func columnTypesQueryFunc(pSQLBuilder ormshift.SQLBuilder) func(string) string {
	return func(pTableName string) string {
		return fmt.Sprintf("SELECT * FROM %s WHERE 1=0", pSQLBuilder.QuoteIdentifier(pTableName))
	}
}
