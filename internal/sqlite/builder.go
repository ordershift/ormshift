package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/internal"
)

type sqliteSQLBuilder struct {
	generic *internal.GenericSQLBuilder
}

func (sb sqliteSQLBuilder) CreateTable(pTable ormshift.Table) string {
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

func (sb sqliteSQLBuilder) DropTable(pTableName ormshift.TableName) string {
	return sb.withGeneric().DropTable(pTableName)
}

func (sb sqliteSQLBuilder) AlterTableAddColumn(pTableName ormshift.TableName, pColumn ormshift.Column) string {
	return sb.withGeneric().AlterTableAddColumn(pTableName, pColumn)
}

func (sb sqliteSQLBuilder) AlterTableDropColumn(pTableName ormshift.TableName, pColumnName ormshift.ColumnName) string {
	return sb.withGeneric().AlterTableDropColumn(pTableName, pColumnName)
}

func (sb sqliteSQLBuilder) ColumnTypeAsString(pColumnType ormshift.ColumnType) string {
	switch pColumnType {
	case ormshift.Varchar:
		return "TEXT"
	case ormshift.Boolean:
		return "INTEGER"
	case ormshift.Integer:
		return "INTEGER"
	case ormshift.DateTime:
		return "DATETIME"
	case ormshift.Monetary:
		return "REAL"
	case ormshift.Decimal:
		return "REAL"
	case ormshift.Binary:
		return "BLOB"
	default:
		return "TEXT"
	}
}

func (sb sqliteSQLBuilder) columnDefinition(pColumn ormshift.Column) string {
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
	return sb.withGeneric().Insert(pTableName, pColumns)
}

func (sb sqliteSQLBuilder) InsertWithValues(pTableName string, pColumnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.withGeneric().InsertWithValues(pTableName, pColumnsValues)
}

func (sb sqliteSQLBuilder) Update(pTableName string, pColumns, pColumnsWhere []string) string {
	return sb.withGeneric().Update(pTableName, pColumns, pColumnsWhere)
}

func (sb sqliteSQLBuilder) UpdateWithValues(pTableName string, pColumns, pColumnsWhere []string, pValues ormshift.ColumnsValues) (string, []any) {
	return sb.withGeneric().UpdateWithValues(pTableName, pColumns, pColumnsWhere, pValues)
}

func (sb sqliteSQLBuilder) Delete(pTableName string, pColumnsWhere []string) string {
	return sb.withGeneric().Delete(pTableName, pColumnsWhere)
}

func (sb sqliteSQLBuilder) DeleteWithValues(pTableName string, pWhereColumnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.withGeneric().DeleteWithValues(pTableName, pWhereColumnsValues)
}

func (sb sqliteSQLBuilder) Select(pTableName string, pColumns, pColumnsWhere []string) string {
	return sb.withGeneric().Select(pTableName, pColumns, pColumnsWhere)
}

func (sb sqliteSQLBuilder) SelectWithValues(pTableName string, pColumns []string, pWhereColumnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.withGeneric().SelectWithValues(pTableName, pColumns, pWhereColumnsValues)
}

func (sb sqliteSQLBuilder) SelectWithPagination(pSQLSelectCommand string, pRowsPerPage, pPageNumber uint) string {
	return sb.withGeneric().SelectWithPagination(pSQLSelectCommand, pRowsPerPage, pPageNumber)
}

func (sb sqliteSQLBuilder) InteroperateSQLCommandWithNamedArgs(pSQLCommand string, pNamedArgs ...sql.NamedArg) (string, []any) {
	return sb.withGeneric().InteroperateSQLCommandWithNamedArgs(pSQLCommand, pNamedArgs...)
}

func (sb sqliteSQLBuilder) withGeneric() internal.GenericSQLBuilder {
	if sb.generic == nil {
		temp := internal.NewGenericSQLBuilder(sb.columnDefinition, nil)
		sb.generic = &temp
	}
	return *sb.generic
}
