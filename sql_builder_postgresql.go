package ormshift

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
)

type postgresqlSQLBuilder struct{}

func (sb postgresqlSQLBuilder) CreateTable(pTable Table) string {
	return sb.generic().CreateTable(pTable)
}

func (sb postgresqlSQLBuilder) DropTable(pTableName TableName) string {
	return sb.generic().DropTable(pTableName)
}

func (sb postgresqlSQLBuilder) AlterTableAddColumn(pTableName TableName, pColumn Column) string {
	return sb.generic().AlterTableAddColumn(pTableName, pColumn)
}

func (sb postgresqlSQLBuilder) AlterTableDropColumn(pTableName TableName, pColumnName ColumnName) string {
	return sb.generic().AlterTableDropColumn(pTableName, pColumnName)
}

func (sb postgresqlSQLBuilder) ColumnTypeAsString(pColumnType ColumnType) string {
	switch pColumnType {
	case Varchar:
		return "VARCHAR"
	case Boolean:
		return "SMALLINT"
	case Integer:
		return "BIGINT"
	case DateTime:
		return "TIMESTAMP(6)"
	case Monetary:
		return "NUMERIC(17,2)"
	case Decimal:
		return "DOUBLE PRECISION"
	case Binary:
		return "BYTEA"
	default:
		return "VARCHAR"
	}
}

func (sb postgresqlSQLBuilder) columnDefinition(pColumn Column) string {
	lColumnDef := pColumn.Name().String()
	if pColumn.Autoincrement() {
		lColumnDef += " BIGSERIAL"
	} else {
		if pColumn.Type() == Varchar {
			lColumnDef += fmt.Sprintf(" %s(%d)", sb.ColumnTypeAsString(pColumn.Type()), pColumn.Size())
		} else {
			lColumnDef += fmt.Sprintf(" %s", sb.ColumnTypeAsString(pColumn.Type()))
		}
	}
	if pColumn.NotNull() {
		lColumnDef += " NOT NULL"
	}
	return lColumnDef
}

func (sb postgresqlSQLBuilder) Insert(pTableName string, pColumns []string) string {
	return sb.generic().Insert(pTableName, pColumns)
}

func (sb postgresqlSQLBuilder) InsertWithValues(pTableName string, pColumnsValues ColumnsValues) (string, []any) {
	return sb.generic().InsertWithValues(pTableName, pColumnsValues)
}

func (sb postgresqlSQLBuilder) Update(pTableName string, pColumns, pColumnsWhere []string) string {
	return sb.generic().Update(pTableName, pColumns, pColumnsWhere)
}

func (sb postgresqlSQLBuilder) UpdateWithValues(pTableName string, pColumns, pColumnsWhere []string, pValues ColumnsValues) (string, []any) {
	return sb.generic().UpdateWithValues(pTableName, pColumns, pColumnsWhere, pValues)
}

func (sb postgresqlSQLBuilder) Delete(pTableName string, pColumnsWhere []string) string {
	return sb.generic().Delete(pTableName, pColumnsWhere)
}

func (sb postgresqlSQLBuilder) DeleteWithValues(pTableName string, pWhereColumnsValues ColumnsValues) (string, []any) {
	return sb.generic().DeleteWithValues(pTableName, pWhereColumnsValues)
}

func (sb postgresqlSQLBuilder) Select(pTableName string, pColumns, pColumnsWhere []string) string {
	return sb.generic().Select(pTableName, pColumns, pColumnsWhere)
}

func (sb postgresqlSQLBuilder) SelectWithValues(pTableName string, pColumns []string, pWhereColumnsValues ColumnsValues) (string, []any) {
	return sb.generic().SelectWithValues(pTableName, pColumns, pWhereColumnsValues)
}

func (sb postgresqlSQLBuilder) SelectWithPagination(pSQLSelectCommand string, pRowsPerPage, pPageNumber uint) string {
	return sb.generic().SelectWithPagination(pSQLSelectCommand, pRowsPerPage, pPageNumber)
}

func (sb postgresqlSQLBuilder) InteroperateSQLCommandWithNamedArgs(pSQLCommand string, pNamedArgs ...sql.NamedArg) (string, []any) {
	lSQLCommand := pSQLCommand
	lArgs := []any{}
	lIndices := map[string]int{}
	for i, lParametro := range pNamedArgs {
		lIndices[strings.ToLower(lParametro.Name)] = i + 1
		lValorBooleano, ok := lParametro.Value.(bool)
		if ok {
			if lValorBooleano {
				lArgs = append(lArgs, int(1))
			} else {
				lArgs = append(lArgs, int(0))
			}
		} else {
			lArgs = append(lArgs, lParametro.Value)
		}
	}
	lRegex := regexp.MustCompile(`@([a-zA-Z_][a-zA-Z0-9_]*)`)
	lSQLCommand = lRegex.ReplaceAllStringFunc(lSQLCommand, func(m string) string {
		lNome := m[1:]
		if idx, ok := lIndices[strings.ToLower(lNome)]; ok {
			return fmt.Sprintf("$%d", idx)
		}
		return m
	})
	return lSQLCommand, lArgs
}

func (sb postgresqlSQLBuilder) generic() genericSQLBuilder {
	return newGenericSQLBuilder(sb.columnDefinition, sb.InteroperateSQLCommandWithNamedArgs)
}
