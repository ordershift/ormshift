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

type InteroperateSQLCommandWithNamedArgsFunc func(pSQLCommand string, pNamedArgs ...sql.NamedArg) (string, []any)

type genericSQLBuilder struct {
	ColumnDefinitionFunc                    ColumnDefinitionFunc
	InteroperateSQLCommandWithNamedArgsFunc InteroperateSQLCommandWithNamedArgsFunc
	QuoteIdentifierFunc                     QuoteIdentifierFunc
}

func NewGenericSQLBuilder(
	pColumnDefinitionFunc ColumnDefinitionFunc,
	pQuoteIdentifierFunc QuoteIdentifierFunc,
	pInteroperateSQLCommandWithNamedArgsFunc InteroperateSQLCommandWithNamedArgsFunc,
) ormshift.SQLBuilder {
	return &genericSQLBuilder{
		ColumnDefinitionFunc:                    pColumnDefinitionFunc,
		QuoteIdentifierFunc:                     pQuoteIdentifierFunc,
		InteroperateSQLCommandWithNamedArgsFunc: pInteroperateSQLCommandWithNamedArgsFunc,
	}
}

func (sb *genericSQLBuilder) CreateTable(pTable schema.Table) string {
	columns := ""
	pkColumns := ""
	for _, column := range pTable.Columns() {
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
		columns += fmt.Sprintf("PRIMARY KEY (%s)", pkColumns)
	}
	return fmt.Sprintf("CREATE TABLE %s (%s);", sb.QuoteIdentifier(pTable.Name()), columns)
}

func (sb *genericSQLBuilder) DropTable(pTableName string) string {
	return fmt.Sprintf("DROP TABLE %s;", sb.QuoteIdentifier(pTableName))
}

func (sb *genericSQLBuilder) AlterTableAddColumn(pTableName string, pColumn schema.Column) string {
	return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s;", sb.QuoteIdentifier(pTableName), sb.columnDefinition(pColumn))
}

func (sb *genericSQLBuilder) AlterTableDropColumn(pTableName, pColumnName string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;", sb.QuoteIdentifier(pTableName), sb.QuoteIdentifier(pColumnName))
}

func (sb *genericSQLBuilder) ColumnTypeAsString(pColumnType schema.ColumnType) string {
	// Generic implementation, should be overridden by specific SQL builders
	return fmt.Sprintf("<<TYPE_%d>>", pColumnType)
}

func (sb *genericSQLBuilder) columnDefinition(pColumn schema.Column) string {
	if sb.ColumnDefinitionFunc != nil {
		return sb.ColumnDefinitionFunc(pColumn)
	}
	return fmt.Sprintf("%s %s", sb.QuoteIdentifier(pColumn.Name()), sb.ColumnTypeAsString(pColumn.Type()))
}

func (sb *genericSQLBuilder) Insert(pTableName string, pColumns []string) string {
	return fmt.Sprintf("insert into %s (%s) values (%s)", sb.QuoteIdentifier(pTableName), sb.columnsList(pColumns), sb.namesList(pColumns))
}

func (sb *genericSQLBuilder) InsertWithValues(pTableName string, pColumnsValues ormshift.ColumnsValues) (string, []any) {
	insertSQL := sb.Insert(pTableName, pColumnsValues.ToColumns())
	insertArgs := pColumnsValues.ToNamedArgs()
	return sb.InteroperateSQLCommandWithNamedArgs(insertSQL, insertArgs...)
}

func (sb *genericSQLBuilder) Update(pTableName string, pColumns, pColumnsWhere []string) string {
	update := fmt.Sprintf("update %s set %s ", sb.QuoteIdentifier(pTableName), sb.columnEqualNameList(pColumns, ","))
	if len(pColumnsWhere) > 0 {
		update += fmt.Sprintf("where %s", sb.columnEqualNameList(pColumnsWhere, " and ")) // NOSONAR go:S1192 - duplicate tradeoff accepted
	}
	return update
}

func (sb *genericSQLBuilder) UpdateWithValues(pTableName string, pColumns, pColumnsWhere []string, pValues ormshift.ColumnsValues) (string, []any) {
	updateSQL := sb.Update(pTableName, pColumns, pColumnsWhere)
	updateArgs := pValues.ToNamedArgs()
	return sb.InteroperateSQLCommandWithNamedArgs(updateSQL, updateArgs...)
}

func (sb *genericSQLBuilder) Delete(pTableName string, pColumnsWhere []string) string {
	delete := fmt.Sprintf("delete from %s ", sb.QuoteIdentifier(pTableName))
	if len(pColumnsWhere) > 0 {
		delete += fmt.Sprintf("where %s", sb.columnEqualNameList(pColumnsWhere, " and ")) // NOSONAR go:S1192 - duplicate tradeoff accepted
	}
	return delete
}

func (sb *genericSQLBuilder) DeleteWithValues(pTableName string, pWhereColumnsValues ormshift.ColumnsValues) (string, []any) {
	deleteSQL := sb.Delete(pTableName, pWhereColumnsValues.ToColumns())
	deleteArgs := pWhereColumnsValues.ToNamedArgs()
	return sb.InteroperateSQLCommandWithNamedArgs(deleteSQL, deleteArgs...)
}

func (sb *genericSQLBuilder) Select(pTableName string, pColumns, pColumnsWhere []string) string {
	update := fmt.Sprintf("select %s from %s ", sb.columnsList(pColumns), sb.QuoteIdentifier(pTableName))
	if len(pColumnsWhere) > 0 {
		update += fmt.Sprintf("where %s", sb.columnEqualNameList(pColumnsWhere, " and ")) // NOSONAR go:S1192 - duplicate tradeoff accepted
	}
	return update
}

func (sb *genericSQLBuilder) SelectWithValues(pTableName string, pColumns []string, pWhereColumnsValues ormshift.ColumnsValues) (string, []any) {
	selectSQL := sb.Select(pTableName, pColumns, pWhereColumnsValues.ToColumns())
	selectArgs := pWhereColumnsValues.ToNamedArgs()
	return sb.InteroperateSQLCommandWithNamedArgs(selectSQL, selectArgs...)
}

func (sb *genericSQLBuilder) SelectWithPagination(pSQLSelectCommand string, pRowsPerPage, pPageNumber uint) string {
	selectWithPagination := pSQLSelectCommand
	if pRowsPerPage > 0 {
		selectWithPagination += fmt.Sprintf(" LIMIT %d", pRowsPerPage)
		if pPageNumber > 1 {
			selectWithPagination += fmt.Sprintf(" OFFSET %d", pRowsPerPage*(pPageNumber-1))
		}
	}
	return selectWithPagination
}

func (sb *genericSQLBuilder) columnsList(pColumns []string) string {
	quotedColumns := []string{}
	for _, col := range pColumns {
		quotedColumns = append(quotedColumns, sb.QuoteIdentifier(col))
	}
	return strings.Join(quotedColumns, ",")
}

func (sb *genericSQLBuilder) namesList(pColumns []string) string {
	names := []string{}
	for _, column := range pColumns {
		names = append(names, "@"+column)
	}
	return strings.Join(names, ",")
}

func (sb *genericSQLBuilder) columnEqualNameList(pColumns []string, pSeparator string) string {
	columnEqualNameList := ""
	for _, column := range pColumns {
		if columnEqualNameList != "" {
			columnEqualNameList += pSeparator
		}
		columnEqualNameList += fmt.Sprintf("%s = @%s", sb.QuoteIdentifier(column), column)
	}
	return columnEqualNameList
}

func (sb *genericSQLBuilder) QuoteIdentifier(pIdentifier string) string {
	if sb.QuoteIdentifierFunc != nil {
		return sb.QuoteIdentifierFunc(pIdentifier)
	}

	// Most databases uses double quotes: "identifier" (PostgreSQL, SQLite, etc.)
	// Escape rule: double quote becomes two double quotes
	// Example: users -> "users", table"name -> "table""name"
	pIdentifier = strings.ReplaceAll(pIdentifier, `"`, `""`)
	return fmt.Sprintf(`"%s"`, pIdentifier)
}

func (sb *genericSQLBuilder) InteroperateSQLCommandWithNamedArgs(pSQLCommand string, pNamedArgs ...sql.NamedArg) (string, []any) {
	if sb.InteroperateSQLCommandWithNamedArgsFunc != nil {
		return sb.InteroperateSQLCommandWithNamedArgsFunc(pSQLCommand, pNamedArgs...)
	}

	sqlCommand := pSQLCommand
	args := []any{}
	for _, param := range pNamedArgs {
		args = append(args, param)
	}
	return sqlCommand, args
}
