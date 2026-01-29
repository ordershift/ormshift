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

func (sb *sqliteBuilder) DropTable(tableName string) string {
	return sb.generic.DropTable(tableName)
}

func (sb *sqliteBuilder) AlterTableAddColumn(tableName string, column schema.Column) string {
	return sb.generic.AlterTableAddColumn(tableName, column)
}

func (sb *sqliteBuilder) AlterTableDropColumn(tableName, columnName string) string {
	return sb.generic.AlterTableDropColumn(tableName, columnName)
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

func (sb *sqliteBuilder) Insert(tableName string, columns []string) string {
	return sb.generic.Insert(tableName, columns)
}

func (sb *sqliteBuilder) InsertWithValues(tableName string, columnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.InsertWithValues(tableName, columnsValues)
}

func (sb *sqliteBuilder) Update(tableName string, columns, columnsWhere []string) string {
	return sb.generic.Update(tableName, columns, columnsWhere)
}

func (sb *sqliteBuilder) UpdateWithValues(tableName string, columns, columnsWhere []string, values ormshift.ColumnsValues) (string, []any) {
	return sb.generic.UpdateWithValues(tableName, columns, columnsWhere, values)
}

func (sb *sqliteBuilder) Delete(tableName string, columnsWhere []string) string {
	return sb.generic.Delete(tableName, columnsWhere)
}

func (sb *sqliteBuilder) DeleteWithValues(tableName string, whereColumnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.DeleteWithValues(tableName, whereColumnsValues)
}

func (sb *sqliteBuilder) Select(tableName string, columns, columnsWhere []string) string {
	return sb.generic.Select(tableName, columns, columnsWhere)
}

func (sb *sqliteBuilder) SelectWithValues(tableName string, columns []string, whereColumnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.SelectWithValues(tableName, columns, whereColumnsValues)
}

func (sb *sqliteBuilder) SelectWithPagination(sqlSelectCommand string, rowsPerPage, pageNumber uint) string {
	return sb.generic.SelectWithPagination(sqlSelectCommand, rowsPerPage, pageNumber)
}

func (sb *sqliteBuilder) QuoteIdentifier(identifier string) string {
	return sb.generic.QuoteIdentifier(identifier)
}

func (sb *sqliteBuilder) InteroperateSQLCommandWithNamedArgs(sqlCommand string, namedArgs ...sql.NamedArg) (string, []any) {
	return sb.generic.InteroperateSQLCommandWithNamedArgs(sqlCommand, namedArgs...)
}
