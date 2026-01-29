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

func (sb *sqlserverBuilder) DropTable(table string) string {
	return sb.generic.DropTable(table)
}

func (sb *sqlserverBuilder) AlterTableAddColumn(table string, column schema.Column) string {
	return sb.generic.AlterTableAddColumn(table, column)
}

func (sb *sqlserverBuilder) AlterTableDropColumn(table, column string) string {
	return sb.generic.AlterTableDropColumn(table, column)
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

func (sb *sqlserverBuilder) Insert(table string, columns []string) string {
	return sb.generic.Insert(table, columns)
}

func (sb *sqlserverBuilder) InsertWithValues(table string, values ormshift.ColumnsValues) (string, []any) {
	return sb.generic.InsertWithValues(table, values)
}

func (sb *sqlserverBuilder) Update(table string, columns, where []string) string {
	return sb.generic.Update(table, columns, where)
}

func (sb *sqlserverBuilder) UpdateWithValues(table string, columns, where []string, values ormshift.ColumnsValues) (string, []any) {
	return sb.generic.UpdateWithValues(table, columns, where, values)
}

func (sb *sqlserverBuilder) Delete(table string, where []string) string {
	return sb.generic.Delete(table, where)
}

func (sb *sqlserverBuilder) DeleteWithValues(table string, where ormshift.ColumnsValues) (string, []any) {
	return sb.generic.DeleteWithValues(table, where)
}

func (sb *sqlserverBuilder) Select(table string, columns, where []string) string {
	return sb.generic.Select(table, columns, where)
}

func (sb *sqlserverBuilder) SelectWithValues(table string, columns []string, where ormshift.ColumnsValues) (string, []any) {
	return sb.generic.SelectWithValues(table, columns, where)
}

func (sb *sqlserverBuilder) SelectWithPagination(sql string, size, number uint) string {
	selectWithPagination := sql
	if size > 0 {
		offset := uint(0)
		if number > 1 {
			offset = size * (number - 1)
		}
		selectWithPagination += fmt.Sprintf(" OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", offset, size)
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

func (sb *sqlserverBuilder) InteroperateSQLCommandWithNamedArgs(sql string, args ...sql.NamedArg) (string, []any) {
	return sb.generic.InteroperateSQLCommandWithNamedArgs(sql, args...)
}
