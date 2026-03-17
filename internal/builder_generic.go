package internal

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/schema"
)

type ColumnDefinitionFunc func(schema.Column) string

type QuoteIdentifierFunc func(string) string

type InteroperateSQLCommandWithNamedArgsFunc func(sql string, args ...sql.NamedArg) (string, []any)

type genericSQLBuilder struct {
	ColumnDefinitionFunc                    ColumnDefinitionFunc
	InteroperateSQLCommandWithNamedArgsFunc InteroperateSQLCommandWithNamedArgsFunc
	QuoteIdentifierFunc                     QuoteIdentifierFunc
}

func NewGenericSQLBuilder(
	columnDefinitionFunc ColumnDefinitionFunc,
	quoteIdentifierFunc QuoteIdentifierFunc,
	interoperateSQLCommandWithNamedArgsFunc InteroperateSQLCommandWithNamedArgsFunc,
) ormshift.SQLBuilder {
	return &genericSQLBuilder{
		ColumnDefinitionFunc:                    columnDefinitionFunc,
		QuoteIdentifierFunc:                     quoteIdentifierFunc,
		InteroperateSQLCommandWithNamedArgsFunc: interoperateSQLCommandWithNamedArgsFunc,
	}
}

func (sb *genericSQLBuilder) CreateTable(table schema.Table) string {
	columns := sb.buildCreateTableColumnDefs(table)
	columns = sb.appendPKConstraint(columns, table)
	columns = sb.appendFKConstraints(columns, table)
	columns = sb.appendUCConstraints(columns, table)
	return fmt.Sprintf("CREATE TABLE %s (%s);", sb.QuoteIdentifier(table.Name()), columns)
}

func (sb *genericSQLBuilder) buildCreateTableColumnDefs(table schema.Table) string {
	columns := ""
	for _, column := range table.Columns() {
		if columns != "" {
			columns += ","
		}
		columns += sb.columnDefinition(column)
	}
	return columns
}

func (sb *genericSQLBuilder) appendPKConstraint(columns string, table schema.Table) string {
	pk := table.PrimaryKey()
	if pk == nil {
		return columns
	}
	return columns + fmt.Sprintf(", CONSTRAINT %s PRIMARY KEY (%s)", sb.QuoteIdentifier(pk.Name()), sb.quotedColumnList(pk.Columns()))
}

func (sb *genericSQLBuilder) appendFKConstraints(columns string, table schema.Table) string {
	for _, fk := range table.ForeignKeys() {
		columns += fmt.Sprintf(", CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)",
			sb.QuoteIdentifier(fk.Name()), sb.quotedColumnList(fk.FromColumns()), sb.QuoteIdentifier(fk.ToTable()), sb.quotedColumnList(fk.ToColumns()))
	}
	return columns
}

func (sb *genericSQLBuilder) appendUCConstraints(columns string, table schema.Table) string {
	for _, uc := range table.UniqueConstraints() {
		columns += fmt.Sprintf(", CONSTRAINT %s UNIQUE (%s)", sb.QuoteIdentifier(uc.Name()), sb.quotedColumnList(uc.Columns()))
	}
	return columns
}

func (sb *genericSQLBuilder) quotedColumnList(cols []string) string {
	parts := make([]string, len(cols))
	for i, col := range cols {
		parts[i] = sb.QuoteIdentifier(col)
	}
	return strings.Join(parts, ",")
}

func (sb *genericSQLBuilder) DropTable(table string) string {
	return fmt.Sprintf("DROP TABLE %s;", sb.QuoteIdentifier(table))
}

func (sb *genericSQLBuilder) AlterTableAddColumn(table string, column schema.Column) string {
	// When column has no user default, columnDefinition does not add DEFAULT; for NotNull we add a hardcoded default here.
	// When column has Default() set, columnDefinition (generic or dialect) already includes " DEFAULT ..." in the column def.
	defaultValue := ""
	if column.Default() == "" && column.NotNull() {
		defaultValue = " DEFAULT "
		switch column.Type() {
		case schema.Boolean, schema.Integer:
			defaultValue += "0"
		case schema.DateTime:
			defaultValue += "'1900-01-01 00:00:00.000000'"
		case schema.Monetary, schema.Decimal:
			defaultValue += "0.0"
		case schema.DateTimeOffSet:
			defaultValue += "'1900-01-01 00:00:00.000000 +00:00'"
		default:
			defaultValue += "''"
		}
	}
	return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s%s;", sb.QuoteIdentifier(table), sb.columnDefinition(column), defaultValue)
}

func (sb *genericSQLBuilder) AlterTableDropColumn(table, column string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;", sb.QuoteIdentifier(table), sb.QuoteIdentifier(column))
}

func (sb *genericSQLBuilder) ColumnTypeAsString(columnType schema.ColumnType) string {
	// Generic implementation, should be overridden by specific SQL builders
	return fmt.Sprintf("<<TYPE_%d>>", columnType)
}

func (sb *genericSQLBuilder) columnDefinition(column schema.Column) string {
	if sb.ColumnDefinitionFunc != nil {
		return sb.ColumnDefinitionFunc(column)
	}
	columnDef := fmt.Sprintf("%s %s", sb.QuoteIdentifier(column.Name()), sb.ColumnTypeAsString(column.Type()))
	if column.Default() != "" {
		columnDef += " DEFAULT " + column.Default()
	}
	if column.Check() != "" {
		columnDef += " CHECK (" + column.Check() + ")"
	}
	return columnDef
}

func (sb *genericSQLBuilder) Insert(table string, columns []string) string {
	return fmt.Sprintf("insert into %s (%s) values (%s)", sb.QuoteIdentifier(table), sb.columnsList(columns), sb.namesList(columns))
}

func (sb *genericSQLBuilder) InsertWithValues(table string, values ormshift.ColumnsValues) (string, []any) {
	insertSQL := sb.Insert(table, values.ToColumns())
	insertArgs := values.ToNamedArgs()
	return sb.InteroperateSQLCommandWithNamedArgs(insertSQL, insertArgs...)
}

func (sb *genericSQLBuilder) Update(table string, columns, where []string) string {
	update := fmt.Sprintf("update %s set %s ", sb.QuoteIdentifier(table), sb.columnEqualNameList(columns, ","))
	if len(where) > 0 {
		update += fmt.Sprintf("where %s", sb.columnEqualNameList(where, " and ")) // NOSONAR go:S1192 - duplicate tradeoff accepted
	}
	return update
}

func (sb *genericSQLBuilder) UpdateWithValues(table string, columns, where []string, values ormshift.ColumnsValues) (string, []any) {
	updateSQL := sb.Update(table, columns, where)
	updateArgs := values.ToNamedArgs()
	return sb.InteroperateSQLCommandWithNamedArgs(updateSQL, updateArgs...)
}

func (sb *genericSQLBuilder) Delete(table string, where []string) string {
	delete := fmt.Sprintf("delete from %s ", sb.QuoteIdentifier(table))
	if len(where) > 0 {
		delete += fmt.Sprintf("where %s", sb.columnEqualNameList(where, " and ")) // NOSONAR go:S1192 - duplicate tradeoff accepted
	}
	return delete
}

func (sb *genericSQLBuilder) DeleteWithValues(table string, where ormshift.ColumnsValues) (string, []any) {
	deleteSQL := sb.Delete(table, where.ToColumns())
	deleteArgs := where.ToNamedArgs()
	return sb.InteroperateSQLCommandWithNamedArgs(deleteSQL, deleteArgs...)
}

func (sb *genericSQLBuilder) Select(table string, columns, where []string) string {
	update := fmt.Sprintf("select %s from %s ", sb.columnsList(columns), sb.QuoteIdentifier(table))
	if len(where) > 0 {
		update += fmt.Sprintf("where %s", sb.columnEqualNameList(where, " and ")) // NOSONAR go:S1192 - duplicate tradeoff accepted
	}
	return update
}

func (sb *genericSQLBuilder) SelectWithValues(table string, columns []string, where ormshift.ColumnsValues) (string, []any) {
	selectSQL := sb.Select(table, columns, where.ToColumns())
	selectArgs := where.ToNamedArgs()
	return sb.InteroperateSQLCommandWithNamedArgs(selectSQL, selectArgs...)
}

func (sb *genericSQLBuilder) SelectWithPagination(sql string, size, number uint) string {
	selectWithPagination := sql
	if size > 0 {
		selectWithPagination += fmt.Sprintf(" LIMIT %d", size)
		if number > 1 {
			selectWithPagination += fmt.Sprintf(" OFFSET %d", size*(number-1))
		}
	}
	return selectWithPagination
}

func (sb *genericSQLBuilder) columnsList(columns []string) string {
	quotedColumns := []string{}
	for _, col := range columns {
		quotedColumns = append(quotedColumns, sb.QuoteIdentifier(col))
	}
	return strings.Join(quotedColumns, ",")
}

func (sb *genericSQLBuilder) namesList(columns []string) string {
	names := []string{}
	for _, column := range columns {
		names = append(names, "@"+column)
	}
	return strings.Join(names, ",")
}

func (sb *genericSQLBuilder) columnEqualNameList(columns []string, separator string) string {
	columnEqualNameList := ""
	for _, column := range columns {
		if columnEqualNameList != "" {
			columnEqualNameList += separator
		}
		columnEqualNameList += fmt.Sprintf("%s = @%s", sb.QuoteIdentifier(column), column)
	}
	return columnEqualNameList
}

func (sb *genericSQLBuilder) QuoteIdentifier(identifier string) string {
	if sb.QuoteIdentifierFunc != nil {
		return sb.QuoteIdentifierFunc(identifier)
	}

	// Most databases uses double quotes: "identifier" (PostgreSQL, SQLite, etc.)
	// Escape rule: double quote becomes two double quotes
	// Example: users -> "users", table"name -> "table""name"
	identifier = strings.ReplaceAll(identifier, `"`, `""`)
	return fmt.Sprintf(`"%s"`, identifier)
}

func (sb *genericSQLBuilder) InteroperateSQLCommandWithNamedArgs(sql string, args ...sql.NamedArg) (string, []any) {
	if sb.InteroperateSQLCommandWithNamedArgsFunc != nil {
		return sb.InteroperateSQLCommandWithNamedArgsFunc(sql, args...)
	}

	a := []any{}
	for _, param := range args {
		a = append(a, param)
	}
	return sql, a
}
