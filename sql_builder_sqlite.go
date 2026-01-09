package ormshift

import (
	"database/sql"
	"fmt"
)

type sqliteSQLBuilder struct{}

func (sb sqliteSQLBuilder) CreateTable(pTable Table) string {
	lColumns := ""
	lPKColumns := ""
	lTemColunaComAutoIncremento := false
	for _, lColumn := range pTable.Columns() {
		if lColumns != "" {
			lColumns += ","
		}
		lColumns += sb.columnDefinition(lColumn)

		if lColumn.PrimaryKey() {
			if lPKColumns != "" {
				lPKColumns += ","
			}
			lPKColumns += lColumn.Name().String()
		}

		if !lTemColunaComAutoIncremento {
			lTemColunaComAutoIncremento = lColumn.Autoincrement()
		}
	}

	if !lTemColunaComAutoIncremento && lPKColumns != "" {
		if lColumns != "" {
			lColumns += ","
		}
		lColumns += fmt.Sprintf("CONSTRAINT PK_%s PRIMARY KEY (%s)", pTable.Name().String(), lPKColumns)
	}
	return fmt.Sprintf("CREATE TABLE %s (%s);", pTable.Name().String(), lColumns)
}

func (sb sqliteSQLBuilder) DropTable(pTableName TableName) string {
	return sb.generic().DropTable(pTableName)
}

func (sb sqliteSQLBuilder) AlterTableAddColumn(pTableName TableName, pColumn Column) string {
	return sb.generic().AlterTableAddColumn(pTableName, pColumn)
}

func (sb sqliteSQLBuilder) AlterTableDropColumn(pTableName TableName, pColumnName ColumnName) string {
	return sb.generic().AlterTableDropColumn(pTableName, pColumnName)
}

func (sb sqliteSQLBuilder) ColumnTypeAsString(pColumnType ColumnType) string {
	switch pColumnType {
	case Varchar:
		return "TEXT"
	case Boolean:
		return "INTEGER"
	case Integer:
		return "INTEGER"
	case DateTime:
		return "DATETIME"
	case Monetary:
		return "REAL"
	case Decimal:
		return "REAL"
	case Binary:
		return "BLOB"
	default:
		return "TEXT"
	}
}

func (sb sqliteSQLBuilder) columnDefinition(pColumn Column) string {
	lColumnDef := fmt.Sprintf("%s %s", pColumn.Name().String(), sb.ColumnTypeAsString(pColumn.Type()))
	if pColumn.NotNull() {
		lColumnDef += " NOT NULL"
	}
	if pColumn.Autoincrement() {
		lColumnDef += " PRIMARY KEY AUTOINCREMENT"
	}
	return lColumnDef
}

func (sb sqliteSQLBuilder) Insert(pTableName string, pColumns []string) string {
	return sb.generic().Insert(pTableName, pColumns)
}

func (sb sqliteSQLBuilder) InsertWithValues(pTableName string, pColumnsValues ColumnsValues) (string, []any) {
	return sb.generic().InsertWithValues(pTableName, pColumnsValues)
}

func (sb sqliteSQLBuilder) Update(pTableName string, pColumns, pColumnsWhere []string) string {
	return sb.generic().Update(pTableName, pColumns, pColumnsWhere)
}

func (sb sqliteSQLBuilder) UpdateWithValues(pTableName string, pColumns, pColumnsWhere []string, pValues ColumnsValues) (string, []any) {
	return sb.generic().UpdateWithValues(pTableName, pColumns, pColumnsWhere, pValues)
}

func (sb sqliteSQLBuilder) Delete(pTableName string, pColumnsWhere []string) string {
	return sb.generic().Delete(pTableName, pColumnsWhere)
}

func (sb sqliteSQLBuilder) DeleteWithValues(pTableName string, pWhereColumnsValues ColumnsValues) (string, []any) {
	return sb.generic().DeleteWithValues(pTableName, pWhereColumnsValues)
}

func (sb sqliteSQLBuilder) Select(pTableName string, pColumns, pColumnsWhere []string) string {
	return sb.generic().Select(pTableName, pColumns, pColumnsWhere)
}

func (sb sqliteSQLBuilder) SelectWithValues(pTableName string, pColumns []string, pWhereColumnsValues ColumnsValues) (string, []any) {
	return sb.generic().SelectWithValues(pTableName, pColumns, pWhereColumnsValues)
}

func (sb sqliteSQLBuilder) SelectWithPagination(pSQLSelectCommand string, pRowsPerPage, pPageNumber uint) string {
	return sb.generic().SelectWithPagination(pSQLSelectCommand, pRowsPerPage, pPageNumber)
}

func (sb sqliteSQLBuilder) InteroperateSQLCommandWithNamedArgs(pSQLCommand string, pNamedArgs ...sql.NamedArg) (string, []any) {
	return sb.generic().InteroperateSQLCommandWithNamedArgs(pSQLCommand, pNamedArgs...)
}

func (sb sqliteSQLBuilder) generic() genericSQLBuilder {
	return newGenericSQLBuilder(sb.columnDefinition, nil)
}
