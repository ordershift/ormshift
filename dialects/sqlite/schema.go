package sqlite

import "fmt"

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

func columnTypesQueryFunc(pTableName string) string {
	return fmt.Sprintf("SELECT * FROM %s WHERE 1=0", QuoteIdentifier(pTableName))
}
