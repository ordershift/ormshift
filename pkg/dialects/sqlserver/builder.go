package sqlserver

import (
	"database/sql"
	"fmt"

	"github.com/ordershift/ormshift/pkg/core"
)

type sqlserverSQLBuilder struct {
	generic *core.GenericSQLBuilder
}

func (sb sqlserverSQLBuilder) CreateTable(pTable core.Table) string {
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

func (sb sqlserverSQLBuilder) DropTable(pTableName core.TableName) string {
	return sb.withGeneric().DropTable(pTableName)
}

func (sb sqlserverSQLBuilder) AlterTableAddColumn(pTableName core.TableName, pColumn core.Column) string {
	return sb.withGeneric().AlterTableAddColumn(pTableName, pColumn)
}

func (sb sqlserverSQLBuilder) AlterTableDropColumn(pTableName core.TableName, pColumnName core.ColumnName) string {
	return sb.withGeneric().AlterTableDropColumn(pTableName, pColumnName)
}

func (sb sqlserverSQLBuilder) ColumnTypeAsString(pColumnType core.ColumnType) string {
	switch pColumnType {
	case core.Varchar:
		return "VARCHAR"
	case core.Boolean:
		return "BIT"
	case core.Integer:
		return "BIGINT"
	case core.DateTime:
		return "DATETIME2(6)"
	case core.Monetary:
		return "MONEY"
	case core.Decimal:
		return "FLOAT"
	case core.Binary:
		return "VARBINARY(MAX)"
	default:
		return "VARCHAR"
	}
}

func (sb sqlserverSQLBuilder) columnDefinition(pColumn core.Column) string {
	lColumnDef := pColumn.Name().String()
	if pColumn.Type() == core.Varchar {
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
	return sb.withGeneric().Insert(pTableName, pColumns)
}

func (sb sqlserverSQLBuilder) InsertWithValues(pTableName string, pColumnsValues core.ColumnsValues) (string, []any) {
	return sb.withGeneric().InsertWithValues(pTableName, pColumnsValues)
}

func (sb sqlserverSQLBuilder) Update(pTableName string, pColumns, pColumnsWhere []string) string {
	return sb.withGeneric().Update(pTableName, pColumns, pColumnsWhere)
}

func (sb sqlserverSQLBuilder) UpdateWithValues(pTableName string, pColumns, pColumnsWhere []string, pValues core.ColumnsValues) (string, []any) {
	return sb.withGeneric().UpdateWithValues(pTableName, pColumns, pColumnsWhere, pValues)
}

func (sb sqlserverSQLBuilder) Delete(pTableName string, pColumnsWhere []string) string {
	return sb.withGeneric().Delete(pTableName, pColumnsWhere)
}

func (sb sqlserverSQLBuilder) DeleteWithValues(pTableName string, pWhereColumnsValues core.ColumnsValues) (string, []any) {
	return sb.withGeneric().DeleteWithValues(pTableName, pWhereColumnsValues)
}

func (sb sqlserverSQLBuilder) Select(pTableName string, pColumns, pColumnsWhere []string) string {
	return sb.withGeneric().Select(pTableName, pColumns, pColumnsWhere)
}

func (sb sqlserverSQLBuilder) SelectWithValues(pTableName string, pColumns []string, pWhereColumnsValues core.ColumnsValues) (string, []any) {
	return sb.withGeneric().SelectWithValues(pTableName, pColumns, pWhereColumnsValues)
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
	return sb.withGeneric().InteroperateSQLCommandWithNamedArgs(pSQLCommand, pNamedArgs...)
}

func (sb sqlserverSQLBuilder) withGeneric() core.GenericSQLBuilder {
	if sb.generic == nil {
		temp := core.NewGenericSQLBuilder(sb.columnDefinition, sb.InteroperateSQLCommandWithNamedArgs)
		sb.generic = &temp
	}
	return *sb.generic
}
