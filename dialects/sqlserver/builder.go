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
	sb.generic = internal.NewGenericSQLBuilder(sb.columnDefinition, sb.QuoteIdentifier, nil)
	return &sb
}

func (sb *sqlserverBuilder) CreateTable(table schema.Table) string {
	columns := ""
	pkColumns := ""
	for _, column := range table.Columns() {
		if columns != "" {
			columns += ","
		}
		columns += sb.columnDefinition(column)

		if column.PrimaryKey() {
			if pkColumns != "" {
				pkColumns += ","
			}
			pkColumns += sb.QuoteIdentifier(column.Name())
		}
	}

	if pkColumns != "" {
		if columns != "" {
			columns += ","
		}
		pkConstraintName := sb.QuoteIdentifier("PK_" + table.Name())
		columns += fmt.Sprintf("CONSTRAINT %s PRIMARY KEY (%s)", pkConstraintName, pkColumns)
	}
	return fmt.Sprintf("CREATE TABLE %s (%s);", sb.QuoteIdentifier(table.Name()), columns)
}

func (sb *sqlserverBuilder) DropTable(tableName string) string {
	return sb.generic.DropTable(tableName)
}

func (sb *sqlserverBuilder) AlterTableAddColumn(tableName string, column schema.Column) string {
	return sb.generic.AlterTableAddColumn(tableName, column)
}

func (sb *sqlserverBuilder) AlterTableDropColumn(tableName, columnName string) string {
	return sb.generic.AlterTableDropColumn(tableName, columnName)
}

func (sb *sqlserverBuilder) ColumnTypeAsString(columnType schema.ColumnType) string {
	switch columnType {
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

func (sb *sqlserverBuilder) columnDefinition(column schema.Column) string {
	columnDef := sb.QuoteIdentifier(column.Name())
	if column.Type() == schema.Varchar {
		columnDef += fmt.Sprintf(" %s(%d)", sb.ColumnTypeAsString(column.Type()), column.Size())
	} else {
		columnDef += fmt.Sprintf(" %s", sb.ColumnTypeAsString(column.Type()))
	}
	if column.NotNull() {
		columnDef += " NOT NULL"
	}
	if column.AutoIncrement() {
		columnDef += " IDENTITY (1, 1)"
	}
	return columnDef
}

func (sb *sqlserverBuilder) Insert(tableName string, columns []string) string {
	return sb.generic.Insert(tableName, columns)
}

func (sb *sqlserverBuilder) InsertWithValues(tableName string, columnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.InsertWithValues(tableName, columnsValues)
}

func (sb *sqlserverBuilder) Update(tableName string, columns, columnsWhere []string) string {
	return sb.generic.Update(tableName, columns, columnsWhere)
}

func (sb *sqlserverBuilder) UpdateWithValues(tableName string, columns, columnsWhere []string, values ormshift.ColumnsValues) (string, []any) {
	return sb.generic.UpdateWithValues(tableName, columns, columnsWhere, values)
}

func (sb *sqlserverBuilder) Delete(tableName string, columnsWhere []string) string {
	return sb.generic.Delete(tableName, columnsWhere)
}

func (sb *sqlserverBuilder) DeleteWithValues(tableName string, whereColumnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.DeleteWithValues(tableName, whereColumnsValues)
}

func (sb *sqlserverBuilder) Select(tableName string, columns, columnsWhere []string) string {
	return sb.generic.Select(tableName, columns, columnsWhere)
}

func (sb *sqlserverBuilder) SelectWithValues(tableName string, columns []string, whereColumnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.SelectWithValues(tableName, columns, whereColumnsValues)
}

func (sb *sqlserverBuilder) SelectWithPagination(sqlSelectCommand string, rowsPerPage, pageNumber uint) string {
	selectWithPagination := sqlSelectCommand
	if rowsPerPage > 0 {
		offSet := uint(0)
		if pageNumber > 1 {
			offSet = rowsPerPage * (pageNumber - 1)
		}
		selectWithPagination += fmt.Sprintf(" OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", offSet, rowsPerPage)
	}
	return selectWithPagination
}

func (sb *sqlserverBuilder) QuoteIdentifier(identifier string) string {
	// SQL Server uses square brackets: [identifier]
	// Escape rule: ] becomes ]]
	// Example: users -> [users], table]name -> [table]]name]
	identifier = strings.ReplaceAll(identifier, "]", "]]")
	return fmt.Sprintf("[%s]", identifier)
}

func (sb *sqlserverBuilder) InteroperateSQLCommandWithNamedArgs(sqlCommand string, namedArgs ...sql.NamedArg) (string, []any) {
	return sb.generic.InteroperateSQLCommandWithNamedArgs(sqlCommand, namedArgs...)
}
