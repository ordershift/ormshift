package sqlserver

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/internal"
	"github.com/ordershift/ormshift/schema"
)

type sqlserverBuilder struct {
	generic ormshift.SQLBuilder
}

func newSQLServerBuilder() ormshift.SQLBuilder {
	sb := sqlserverBuilder{}
	sb.generic = internal.NewGenericSQLBuilder(sb.columnDefinition, QuoteIdentifier, nil)
	return sb
}

func (sb sqlserverBuilder) CreateTable(pTable schema.Table) string {
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
			lPKColumns += QuoteIdentifier(lColumn.Name())
		}
	}

	if lPKColumns != "" {
		if lColumns != "" {
			lColumns += ","
		}
		lPKConstraintName := QuoteIdentifier("PK_" + pTable.Name())
		lColumns += fmt.Sprintf("CONSTRAINT %s PRIMARY KEY (%s)", lPKConstraintName, lPKColumns)
	}
	return fmt.Sprintf("CREATE TABLE %s (%s);", QuoteIdentifier(pTable.Name()), lColumns)
}

func (sb sqlserverBuilder) DropTable(pTableName string) string {
	return sb.generic.DropTable(pTableName)
}

func (sb sqlserverBuilder) AlterTableAddColumn(pTableName string, pColumn schema.Column) string {
	return sb.generic.AlterTableAddColumn(pTableName, pColumn)
}

func (sb sqlserverBuilder) AlterTableDropColumn(pTableName string, pColumnName string) string {
	return sb.generic.AlterTableDropColumn(pTableName, pColumnName)
}

func (sb sqlserverBuilder) ColumnTypeAsString(pColumnType schema.ColumnType) string {
	switch pColumnType {
	case schema.Varchar:
		return "VARCHAR"
	case schema.Boolean:
		return "BIT"
	case schema.Integer:
		return "BIGINT"
	case schema.DateTime:
		return "DATETIME2(6)"
	case schema.Monetary:
		return "MONEY"
	case schema.Decimal:
		return "FLOAT"
	case schema.Binary:
		return "VARBINARY(MAX)"
	default:
		return "VARCHAR"
	}
}

func (sb sqlserverBuilder) columnDefinition(pColumn schema.Column) string {
	lColumnDef := QuoteIdentifier(pColumn.Name())
	if pColumn.Type() == schema.Varchar {
		lColumnDef += fmt.Sprintf(" %s(%d)", sb.ColumnTypeAsString(pColumn.Type()), pColumn.Size())
	} else {
		lColumnDef += fmt.Sprintf(" %s", sb.ColumnTypeAsString(pColumn.Type()))
	}
	if pColumn.NotNull() {
		lColumnDef += " NOT NULL"
	}
	if pColumn.AutoIncrement() {
		lColumnDef += " IDENTITY (1, 1)"
	}
	return lColumnDef
}

func (sb sqlserverBuilder) Insert(pTableName string, pColumns []string) string {
	return sb.generic.Insert(pTableName, pColumns)
}

func (sb sqlserverBuilder) InsertWithValues(pTableName string, pColumnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.InsertWithValues(pTableName, pColumnsValues)
}

func (sb sqlserverBuilder) Update(pTableName string, pColumns, pColumnsWhere []string) string {
	return sb.generic.Update(pTableName, pColumns, pColumnsWhere)
}

func (sb sqlserverBuilder) UpdateWithValues(pTableName string, pColumns, pColumnsWhere []string, pValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.UpdateWithValues(pTableName, pColumns, pColumnsWhere, pValues)
}

func (sb sqlserverBuilder) Delete(pTableName string, pColumnsWhere []string) string {
	return sb.generic.Delete(pTableName, pColumnsWhere)
}

func (sb sqlserverBuilder) DeleteWithValues(pTableName string, pWhereColumnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.DeleteWithValues(pTableName, pWhereColumnsValues)
}

func (sb sqlserverBuilder) Select(pTableName string, pColumns, pColumnsWhere []string) string {
	return sb.generic.Select(pTableName, pColumns, pColumnsWhere)
}

func (sb sqlserverBuilder) SelectWithValues(pTableName string, pColumns []string, pWhereColumnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.SelectWithValues(pTableName, pColumns, pWhereColumnsValues)
}

func (sb sqlserverBuilder) SelectWithPagination(pSQLSelectCommand string, pRowsPerPage, pPageNumber uint) string {
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

func (sb sqlserverBuilder) InteroperateSQLCommandWithNamedArgs(pSQLCommand string, pNamedArgs ...sql.NamedArg) (string, []any) {
	return sb.generic.InteroperateSQLCommandWithNamedArgs(pSQLCommand, pNamedArgs...)
}

func QuoteIdentifier(pIdentifier string) string {
	// SQL Server uses square brackets: [identifier]
	// Escape rule: ] becomes ]]
	// Example: users -> [users], table]name -> [table]]name]
	pIdentifier = strings.ReplaceAll(pIdentifier, "]", "]]")
	return fmt.Sprintf("[%s]", pIdentifier)
}
