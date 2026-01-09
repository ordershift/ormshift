package ormshift

import (
	"database/sql"
	"fmt"
)

type sqlserverSQLBuilder struct{}

func (sb sqlserverSQLBuilder) CreateTable(pTable Table) string {
	lColumns := ""
	lPKColumns := ""
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
	}

	if lPKColumns != "" {
		if lColumns != "" {
			lColumns += ","
		}
		lColumns += fmt.Sprintf("CONSTRAINT PK_%s PRIMARY KEY (%s)", pTable.Name().String(), lPKColumns)
	}
	return fmt.Sprintf("CREATE TABLE %s (%s);", pTable.Name().String(), lColumns)
}

func (sb sqlserverSQLBuilder) DropTable(pTableName TableName) string {
	return sb.generic().DropTable(pTableName)
}

func (sb sqlserverSQLBuilder) AlterTableAddColumn(pTableName TableName, pColumn Column) string {
	return sb.generic().AlterTableAddColumn(pTableName, pColumn)
}

func (sb sqlserverSQLBuilder) AlterTableDropColumn(pTableName TableName, pColumnName ColumnName) string {
	return sb.generic().AlterTableDropColumn(pTableName, pColumnName)
}

func (sb sqlserverSQLBuilder) ColumnTypeAsString(pColumnType ColumnType) string {
	switch pColumnType {
	case Varchar:
		return "VARCHAR"
	case Boolean:
		return "BIT"
	case Integer:
		return "BIGINT"
	case DateTime:
		return "DATETIME2(6)"
	case Monetary:
		return "MONEY"
	case Decimal:
		return "FLOAT"
	case Binary:
		return "VARBINARY(MAX)"
	default:
		return "VARCHAR"
	}
}

func (sb sqlserverSQLBuilder) columnDefinition(pColumn Column) string {
	lColumnDef := pColumn.Name().String()
	if pColumn.Type() == Varchar {
		lColumnDef += fmt.Sprintf(" %s(%d)", sb.ColumnTypeAsString(pColumn.Type()), pColumn.Size())
	} else {
		lColumnDef += fmt.Sprintf(" %s", sb.ColumnTypeAsString(pColumn.Type()))
	}
	if pColumn.NotNull() {
		lColumnDef += " NOT NULL"
	}
	if pColumn.Autoincrement() {
		lColumnDef += " IDENTITY (1, 1)"
	}
	return lColumnDef
}

func (sb sqlserverSQLBuilder) Insert(pTableName string, pColumns []string) string {
	return sb.generic().Insert(pTableName, pColumns)
}

func (sb sqlserverSQLBuilder) InsertWithValues(pTableName string, pColumnsValues ColumnsValues) (string, []any) {
	return sb.generic().InsertWithValues(pTableName, pColumnsValues)
}

func (sb sqlserverSQLBuilder) Update(pTableName string, pColumns, pColumnsWhere []string) string {
	return sb.generic().Update(pTableName, pColumns, pColumnsWhere)
}

func (sb sqlserverSQLBuilder) UpdateWithValues(pTableName string, pColumns, pColumnsWhere []string, pValues ColumnsValues) (string, []any) {
	return sb.generic().UpdateWithValues(pTableName, pColumns, pColumnsWhere, pValues)
}

func (sb sqlserverSQLBuilder) Delete(pTableName string, pColumnsWhere []string) string {
	return sb.generic().Delete(pTableName, pColumnsWhere)
}

func (sb sqlserverSQLBuilder) DeleteWithValues(pTableName string, pWhereColumnsValues ColumnsValues) (string, []any) {
	return sb.generic().DeleteWithValues(pTableName, pWhereColumnsValues)
}

func (sb sqlserverSQLBuilder) Select(pTableName string, pColumns, pColumnsWhere []string) string {
	return sb.generic().Select(pTableName, pColumns, pColumnsWhere)
}

func (sb sqlserverSQLBuilder) SelectWithValues(pTableName string, pColumns []string, pWhereColumnsValues ColumnsValues) (string, []any) {
	return sb.generic().SelectWithValues(pTableName, pColumns, pWhereColumnsValues)
}

func (sb sqlserverSQLBuilder) SelectWithPagination(pSQLSelectCommand string, pRowsPerPage, pPageNumber uint) string {
	lSelectWithPagination := pSQLSelectCommand
	if pRowsPerPage > 0 {
		lOffSet := uint(0)
		if pPageNumber > 1 {
			lOffSet = pRowsPerPage * (pPageNumber - 1)
		}
		lSelectWithPagination += fmt.Sprintf(" OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", lOffSet, pRowsPerPage)
	}
	return lSelectWithPagination
}

func (sb sqlserverSQLBuilder) InteroperateSQLCommandWithNamedArgs(pSQLCommand string, pNamedArgs ...sql.NamedArg) (string, []any) {
	return sb.generic().InteroperateSQLCommandWithNamedArgs(pSQLCommand, pNamedArgs...)
}

func (sb sqlserverSQLBuilder) generic() genericSQLBuilder {
	return newGenericSQLBuilder(sb.columnDefinition, nil)
}
