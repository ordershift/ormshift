package sqlserver

import (
	"database/sql"

	"github.com/ordershift/ormshift/pkg/core"
)

const tableNamesQuery = `
		SELECT
			t.name
		FROM
			sys.tables t
		LEFT JOIN
			sys.extended_properties ep
		ON	ep.major_id = t.[object_id]
		WHERE
			t.is_ms_shipped = 0 AND
			(ep.class_desc IS NULL OR (ep.class_desc <> 'OBJECT_OR_COLUMN' AND
				ep.[name] <> 'microsoft_database_tools_support'))
		ORDER BY
			t.name
	`

func DBSchema(pDB *sql.DB) (*core.DBSchema, error) {
	return core.NewDBSchema(pDB, tableNamesQuery)
}
