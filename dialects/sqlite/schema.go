package sqlite

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
