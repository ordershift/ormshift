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

func (sb *sqliteBuilder) CreateTable(table schema.Table) string {
	columns := ""
	pkColumns := ""
	hasAutoIncrementColumn := false
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

		if !hasAutoIncrementColumn {
			hasAutoIncrementColumn = column.AutoIncrement()
		}
	}

	if !hasAutoIncrementColumn && pkColumns != "" {
		if columns != "" {
			columns += ","
		}
		pkConstraintName := sb.QuoteIdentifier("PK_" + table.Name())
		columns += fmt.Sprintf("CONSTRAINT %s PRIMARY KEY (%s)", pkConstraintName, pkColumns)
	}
	return fmt.Sprintf("CREATE TABLE %s (%s);", sb.QuoteIdentifier(table.Name()), columns)
}

func (sb *sqliteBuilder) DropTable(table string) string {
	return sb.generic.DropTable(table)
}

func (sb *sqliteBuilder) AlterTableAddColumn(table string, column schema.Column) string {
	return sb.generic.AlterTableAddColumn(table, column)
}

func (sb *sqliteBuilder) AlterTableDropColumn(table, column string) string {
	return sb.generic.AlterTableDropColumn(table, column)
}

func (sb *sqliteBuilder) ColumnTypeAsString(columnType schema.ColumnType) string {
	switch columnType {
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

func (sb *sqliteBuilder) columnDefinition(column schema.Column) string {
	columnDef := fmt.Sprintf("%s %s", sb.QuoteIdentifier(column.Name()), sb.ColumnTypeAsString(column.Type()))
	if column.NotNull() {
		columnDef += " NOT NULL"
	}
	if column.AutoIncrement() {
		columnDef += " PRIMARY KEY AUTOINCREMENT"
	}
	return columnDef
}

func (sb *sqliteBuilder) Insert(table string, columns []string) string {
	return sb.generic.Insert(table, columns)
}

func (sb *sqliteBuilder) InsertWithValues(table string, values ormshift.ColumnsValues) (string, []any) {
	return sb.generic.InsertWithValues(table, values)
}

func (sb *sqliteBuilder) Update(table string, columns, where []string) string {
	return sb.generic.Update(table, columns, where)
}

func (sb *sqliteBuilder) UpdateWithValues(table string, columns, where []string, values ormshift.ColumnsValues) (string, []any) {
	return sb.generic.UpdateWithValues(table, columns, where, values)
}

func (sb *sqliteBuilder) Delete(table string, where []string) string {
	return sb.generic.Delete(table, where)
}

func (sb *sqliteBuilder) DeleteWithValues(table string, where ormshift.ColumnsValues) (string, []any) {
	return sb.generic.DeleteWithValues(table, where)
}

func (sb *sqliteBuilder) Select(table string, columns, where []string) string {
	return sb.generic.Select(table, columns, where)
}

func (sb *sqliteBuilder) SelectWithValues(table string, columns []string, where ormshift.ColumnsValues) (string, []any) {
	return sb.generic.SelectWithValues(table, columns, where)
}

func (sb *sqliteBuilder) SelectWithPagination(sqlSelectCommand string, rowsPerPage, pageNumber uint) string {
	return sb.generic.SelectWithPagination(sqlSelectCommand, rowsPerPage, pageNumber)
}

func (sb *sqliteBuilder) QuoteIdentifier(identifier string) string {
	return sb.generic.QuoteIdentifier(identifier)
}

func (sb *sqliteBuilder) InteroperateSQLCommandWithNamedArgs(sql string, args ...sql.NamedArg) (string, []any) {
	return sb.generic.InteroperateSQLCommandWithNamedArgs(sql, args...)
}
