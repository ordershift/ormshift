package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/ordershift/ormshift/pkg/core"
)

type sqliteSQLBuilder struct {
	generic *core.GenericSQLBuilder
}

func (sb sqliteSQLBuilder) CreateTable(pTable core.Table) string {
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

func (sb sqliteSQLBuilder) DropTable(pTableName core.TableName) string {
	return sb.withGeneric().DropTable(pTableName)
}

func (sb sqliteSQLBuilder) AlterTableAddColumn(pTableName core.TableName, pColumn core.Column) string {
	return sb.withGeneric().AlterTableAddColumn(pTableName, pColumn)
}

func (sb sqliteSQLBuilder) AlterTableDropColumn(pTableName core.TableName, pColumnName core.ColumnName) string {
	return sb.withGeneric().AlterTableDropColumn(pTableName, pColumnName)
}

func (sb sqliteSQLBuilder) ColumnTypeAsString(pColumnType core.ColumnType) string {
	switch pColumnType {
	case core.Varchar:
		return "TEXT"
	case core.Boolean:
		return "INTEGER"
	case core.Integer:
		return "INTEGER"
	case core.DateTime:
		return "DATETIME"
	case core.Monetary:
		return "REAL"
	case core.Decimal:
		return "REAL"
	case core.Binary:
		return "BLOB"
	default:
		return "TEXT"
	}
}

func (sb sqliteSQLBuilder) columnDefinition(pColumn core.Column) string {
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

func (sb sqliteSQLBuilder) InsertWithValues(pTableName string, pColumnsValues core.ColumnsValues) (string, []any) {
	return sb.withGeneric().InsertWithValues(pTableName, pColumnsValues)
}

func (sb sqliteSQLBuilder) Update(pTableName string, pColumns, pColumnsWhere []string) string {
	return sb.withGeneric().Update(pTableName, pColumns, pColumnsWhere)
}

func (sb sqliteSQLBuilder) UpdateWithValues(pTableName string, pColumns, pColumnsWhere []string, pValues core.ColumnsValues) (string, []any) {
	return sb.withGeneric().UpdateWithValues(pTableName, pColumns, pColumnsWhere, pValues)
}

func (sb sqliteSQLBuilder) Delete(pTableName string, pColumnsWhere []string) string {
	return sb.withGeneric().Delete(pTableName, pColumnsWhere)
}

func (sb sqliteSQLBuilder) DeleteWithValues(pTableName string, pWhereColumnsValues core.ColumnsValues) (string, []any) {
	return sb.withGeneric().DeleteWithValues(pTableName, pWhereColumnsValues)
}

func (sb sqliteSQLBuilder) Select(pTableName string, pColumns, pColumnsWhere []string) string {
	return sb.withGeneric().Select(pTableName, pColumns, pColumnsWhere)
}

func (sb sqliteSQLBuilder) SelectWithValues(pTableName string, pColumns []string, pWhereColumnsValues core.ColumnsValues) (string, []any) {
	return sb.withGeneric().SelectWithValues(pTableName, pColumns, pWhereColumnsValues)
}

func (sb sqliteSQLBuilder) SelectWithPagination(pSQLSelectCommand string, pRowsPerPage, pPageNumber uint) string {
	return sb.withGeneric().SelectWithPagination(pSQLSelectCommand, pRowsPerPage, pPageNumber)
}

func (sb sqliteSQLBuilder) InteroperateSQLCommandWithNamedArgs(pSQLCommand string, pNamedArgs ...sql.NamedArg) (string, []any) {
	return sb.withGeneric().InteroperateSQLCommandWithNamedArgs(pSQLCommand, pNamedArgs...)
}

func (sb sqliteSQLBuilder) withGeneric() core.GenericSQLBuilder {
	if sb.generic == nil {
		temp := core.NewGenericSQLBuilder(sb.columnDefinition, nil)
		sb.generic = &temp
	}
	return *sb.generic
}
