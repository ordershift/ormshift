package internal

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/schema"
)

type ColumnDefinitionFunc func(schema.Column) string

type InteroperateSQLCommandWithNamedArgsFunc func(pSQLCommand string, pNamedArgs ...sql.NamedArg) (string, []any)

type genericSQLBuilder struct {
	ColumnDefinitionFunc                    ColumnDefinitionFunc
	InteroperateSQLCommandWithNamedArgsFunc InteroperateSQLCommandWithNamedArgsFunc
}

func NewGenericSQLBuilder(pColumnDefinitionFunc ColumnDefinitionFunc, pInteroperateSQLCommandWithNamedArgsFunc InteroperateSQLCommandWithNamedArgsFunc) ormshift.SQLBuilder {
	return genericSQLBuilder{
		ColumnDefinitionFunc:                    pColumnDefinitionFunc,
		InteroperateSQLCommandWithNamedArgsFunc: pInteroperateSQLCommandWithNamedArgsFunc,
	}
}

func (sb genericSQLBuilder) CreateTable(pTable schema.Table) string {
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
		lColumns += fmt.Sprintf("PRIMARY KEY (%s)", lPKColumns)
	}
	return fmt.Sprintf("CREATE TABLE %s (%s);", pTable.Name().String(), lColumns)
}

func (sb genericSQLBuilder) DropTable(pTableName schema.TableName) string {
	return fmt.Sprintf("DROP TABLE %s;", pTableName.String())
}

func (sb genericSQLBuilder) AlterTableAddColumn(pTableName schema.TableName, pColumn schema.Column) string {
	return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s;", pTableName.String(), sb.columnDefinition(pColumn))
}

func (sb genericSQLBuilder) AlterTableDropColumn(pTableName schema.TableName, pColumnName schema.ColumnName) string {
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;", pTableName.String(), pColumnName.String())
}

func (sb genericSQLBuilder) ColumnTypeAsString(pColumnType schema.ColumnType) string {
	// Generic implementation, should be overridden by specific SQL builders
	return fmt.Sprintf("<<TYPE_%d>>", pColumnType)
}

func (sb genericSQLBuilder) columnDefinition(pColumn schema.Column) string {
	if sb.ColumnDefinitionFunc != nil {
		return sb.ColumnDefinitionFunc(pColumn)
	}
	return fmt.Sprintf("%s %s", pColumn.Name().String(), sb.ColumnTypeAsString(pColumn.Type()))
}

func (sb genericSQLBuilder) Insert(pTableName string, pColumns []string) string {
	return fmt.Sprintf("insert into %s (%s) values (%s)", pTableName, sb.columnsList(pColumns), sb.namesList(pColumns))
}

func (sb genericSQLBuilder) InsertWithValues(pTableName string, pColumnsValues ormshift.ColumnsValues) (string, []any) {
	lInsertSQL := sb.Insert(pTableName, pColumnsValues.ToColumns())
	lInsertArgs := pColumnsValues.ToNamedArgs()
	return sb.InteroperateSQLCommandWithNamedArgs(lInsertSQL, lInsertArgs...)
}

func (sb genericSQLBuilder) Update(pTableName string, pColumns, pColumnsWhere []string) string {
	lUpdate := fmt.Sprintf("update %s set %s ", pTableName, sb.columnEqualNameList(pColumns, ","))
	if len(pColumnsWhere) > 0 {
		lUpdate += fmt.Sprintf("where %s", sb.columnEqualNameList(pColumnsWhere, " and ")) // NOSONAR go:S1192 - duplicate tradeoff accepted
	}
	return lUpdate
}

func (sb genericSQLBuilder) UpdateWithValues(pTableName string, pColumns, pColumnsWhere []string, pValues ormshift.ColumnsValues) (string, []any) {
	lUpdateSQL := sb.Update(pTableName, pColumns, pColumnsWhere)
	lUpdateArgs := pValues.ToNamedArgs()
	return sb.InteroperateSQLCommandWithNamedArgs(lUpdateSQL, lUpdateArgs...)
}

func (sb genericSQLBuilder) Delete(pTableName string, pColumnsWhere []string) string {
	lDelete := fmt.Sprintf("delete from %s ", pTableName)
	if len(pColumnsWhere) > 0 {
		lDelete += fmt.Sprintf("where %s", sb.columnEqualNameList(pColumnsWhere, " and ")) // NOSONAR go:S1192 - duplicate tradeoff accepted
	}
	return lDelete
}

func (sb genericSQLBuilder) DeleteWithValues(pTableName string, pWhereColumnsValues ormshift.ColumnsValues) (string, []any) {
	lDeleteSQL := sb.Delete(pTableName, pWhereColumnsValues.ToColumns())
	lDeleteArgs := pWhereColumnsValues.ToNamedArgs()
	return sb.InteroperateSQLCommandWithNamedArgs(lDeleteSQL, lDeleteArgs...)
}

func (sb genericSQLBuilder) Select(pTableName string, pColumns, pColumnsWhere []string) string {
	lUpdate := fmt.Sprintf("select %s from %s ", sb.columnsList(pColumns), pTableName)
	if len(pColumnsWhere) > 0 {
		lUpdate += fmt.Sprintf("where %s", sb.columnEqualNameList(pColumnsWhere, " and ")) // NOSONAR go:S1192 - duplicate tradeoff accepted
	}
	return lUpdate
}

func (sb genericSQLBuilder) SelectWithValues(pTableName string, pColumns []string, pWhereColumnsValues ormshift.ColumnsValues) (string, []any) {
	lSelectSQL := sb.Select(pTableName, pColumns, pWhereColumnsValues.ToColumns())
	lSelectArgs := pWhereColumnsValues.ToNamedArgs()
	return sb.InteroperateSQLCommandWithNamedArgs(lSelectSQL, lSelectArgs...)
}

func (sb genericSQLBuilder) SelectWithPagination(pSQLSelectCommand string, pRowsPerPage, pPageNumber uint) string {
	lSelectWithPagination := pSQLSelectCommand
	if pRowsPerPage > 0 {
		lSelectWithPagination += fmt.Sprintf(" LIMIT %d", pRowsPerPage)
		if pPageNumber > 1 {
			lSelectWithPagination += fmt.Sprintf(" OFFSET %d", pRowsPerPage*(pPageNumber-1))
		}
	}
	return lSelectWithPagination
}

func (sb genericSQLBuilder) columnsList(pColumns []string) string {
	return strings.Join(pColumns, ",")
}

func (sb genericSQLBuilder) namesList(pColumns []string) string {
	lNames := []string{}
	for _, lColumn := range pColumns {
		lNames = append(lNames, "@"+lColumn)
	}
	return strings.Join(lNames, ",")
}

func (sb genericSQLBuilder) columnEqualNameList(pColumns []string, pSeparator string) string {
	lColumnEqualNameList := ""
	for _, lColumn := range pColumns {
		if lColumnEqualNameList != "" {
			lColumnEqualNameList += pSeparator
		}
		lColumnEqualNameList += fmt.Sprintf("%s = @%s", lColumn, lColumn)
	}
	return lColumnEqualNameList
}

func (sb genericSQLBuilder) InteroperateSQLCommandWithNamedArgs(pSQLCommand string, pNamedArgs ...sql.NamedArg) (string, []any) {
	if sb.InteroperateSQLCommandWithNamedArgsFunc != nil {
		return sb.InteroperateSQLCommandWithNamedArgsFunc(pSQLCommand, pNamedArgs...)
	}
	lSQLCommand := pSQLCommand
	lArgs := []any{}
	for _, lParam := range pNamedArgs {
		lArgs = append(lArgs, lParam)
	}
	return lSQLCommand, lArgs
}
