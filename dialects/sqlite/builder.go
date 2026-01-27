package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/internal"
	"github.com/ordershift/ormshift/schema"
)

type sqliteBuilder struct {
	generic ormshift.SQLBuilder
}

func newSQLiteBuilder() ormshift.SQLBuilder {
	sb := sqliteBuilder{}
	sb.generic = internal.NewGenericSQLBuilder(sb.columnDefinition, nil, nil)
	return &sb
}

func (sb *sqliteBuilder) CreateTable(pTable schema.Table) string {
	lColumns := ""
	lPKColumns := ""
	lHasAutoIncrementColumn := false
	for _, lColumn := range pTable.Columns() {
		if lColumns != "" {
			lColumns += ","
		}
		lColumns += sb.columnDefinition(lColumn)

		if lColumn.PrimaryKey() {
			if lPKColumns != "" {
				lPKColumns += ","
			}
			lPKColumns += sb.QuoteIdentifier(lColumn.Name())
		}

		if !lHasAutoIncrementColumn {
			lHasAutoIncrementColumn = lColumn.AutoIncrement()
		}
	}

	if !lHasAutoIncrementColumn && lPKColumns != "" {
		if lColumns != "" {
			lColumns += ","
		}
		lPKConstraintName := sb.QuoteIdentifier("PK_" + pTable.Name())
		lColumns += fmt.Sprintf("CONSTRAINT %s PRIMARY KEY (%s)", lPKConstraintName, lPKColumns)
	}
	return fmt.Sprintf("CREATE TABLE %s (%s);", sb.QuoteIdentifier(pTable.Name()), lColumns)
}

func (sb *sqliteBuilder) DropTable(pTableName string) string {
	return sb.generic.DropTable(pTableName)
}

func (sb *sqliteBuilder) AlterTableAddColumn(pTableName string, pColumn schema.Column) string {
	return sb.generic.AlterTableAddColumn(pTableName, pColumn)
}

func (sb *sqliteBuilder) AlterTableDropColumn(pTableName, pColumnName string) string {
	return sb.generic.AlterTableDropColumn(pTableName, pColumnName)
}

func (sb *sqliteBuilder) ColumnTypeAsString(pColumnType schema.ColumnType) string {
	switch pColumnType {
	case schema.Varchar:
		return "TEXT"
	case schema.Boolean:
		return "INTEGER"
	case schema.Integer:
		return "INTEGER"
	case schema.DateTime:
		return "DATETIME"
	case schema.Monetary:
		return "REAL"
	case schema.Decimal:
		return "REAL"
	case schema.Binary:
		return "BLOB"
	default:
		return "TEXT"
	}
}

func (sb *sqliteBuilder) columnDefinition(pColumn schema.Column) string {
	lColumnDef := fmt.Sprintf("%s %s", sb.QuoteIdentifier(pColumn.Name()), sb.ColumnTypeAsString(pColumn.Type()))
	if pColumn.NotNull() {
		lColumnDef += " NOT NULL"
	}
	if pColumn.AutoIncrement() {
		lColumnDef += " PRIMARY KEY AUTOINCREMENT"
	}
	return lColumnDef
}

func (sb *sqliteBuilder) Insert(pTableName string, pColumns []string) string {
	return sb.generic.Insert(pTableName, pColumns)
}

func (sb *sqliteBuilder) InsertWithValues(pTableName string, pColumnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.InsertWithValues(pTableName, pColumnsValues)
}

func (sb *sqliteBuilder) Update(pTableName string, pColumns, pColumnsWhere []string) string {
	return sb.generic.Update(pTableName, pColumns, pColumnsWhere)
}

func (sb *sqliteBuilder) UpdateWithValues(pTableName string, pColumns, pColumnsWhere []string, pValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.UpdateWithValues(pTableName, pColumns, pColumnsWhere, pValues)
}

func (sb *sqliteBuilder) Delete(pTableName string, pColumnsWhere []string) string {
	return sb.generic.Delete(pTableName, pColumnsWhere)
}

func (sb *sqliteBuilder) DeleteWithValues(pTableName string, pWhereColumnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.DeleteWithValues(pTableName, pWhereColumnsValues)
}

func (sb *sqliteBuilder) Select(pTableName string, pColumns, pColumnsWhere []string) string {
	return sb.generic.Select(pTableName, pColumns, pColumnsWhere)
}

func (sb *sqliteBuilder) SelectWithValues(pTableName string, pColumns []string, pWhereColumnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.SelectWithValues(pTableName, pColumns, pWhereColumnsValues)
}

func (sb *sqliteBuilder) SelectWithPagination(pSQLSelectCommand string, pRowsPerPage, pPageNumber uint) string {
	return sb.generic.SelectWithPagination(pSQLSelectCommand, pRowsPerPage, pPageNumber)
}

func (sb *sqliteBuilder) QuoteIdentifier(pIdentifier string) string {
	return sb.generic.QuoteIdentifier(pIdentifier)
}

func (sb *sqliteBuilder) InteroperateSQLCommandWithNamedArgs(pSQLCommand string, pNamedArgs ...sql.NamedArg) (string, []any) {
	return sb.generic.InteroperateSQLCommandWithNamedArgs(pSQLCommand, pNamedArgs...)
}
