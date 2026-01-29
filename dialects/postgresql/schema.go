package postgresql

import (
	"fmt"

	"github.com/ordershift/ormshift"
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

func columnTypesQueryFunc(sqlBuilder ormshift.SQLBuilder) func(string) string {
	return func(tableName string) string {
		return fmt.Sprintf("SELECT * FROM %s WHERE 1=0", sqlBuilder.QuoteIdentifier(tableName))
	}
}
