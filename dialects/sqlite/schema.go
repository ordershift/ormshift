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

func columnTypesQueryFunc(sqlBuilder ormshift.SQLBuilder) func(string) string {
	return func(table string) string {
		return fmt.Sprintf("SELECT * FROM %s WHERE 1=0", sqlBuilder.QuoteIdentifier(table))
	}
}
