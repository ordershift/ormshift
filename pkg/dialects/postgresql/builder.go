package postgresql

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/ordershift/ormshift/pkg/core"
)

type postgresqlSQLBuilder struct {
	generic *core.GenericSQLBuilder
}

func (sb postgresqlSQLBuilder) CreateTable(pTable core.Table) string {
	return sb.withGeneric().CreateTable(pTable)
}

func (sb postgresqlSQLBuilder) DropTable(pTableName core.TableName) string {
	return sb.withGeneric().DropTable(pTableName)
}

func (sb postgresqlSQLBuilder) AlterTableAddColumn(pTableName core.TableName, pColumn core.Column) string {
	return sb.withGeneric().AlterTableAddColumn(pTableName, pColumn)
}

func (sb postgresqlSQLBuilder) AlterTableDropColumn(pTableName core.TableName, pColumnName core.ColumnName) string {
	return sb.withGeneric().AlterTableDropColumn(pTableName, pColumnName)
}

func (sb postgresqlSQLBuilder) ColumnTypeAsString(pColumnType core.ColumnType) string {
	switch pColumnType {
	case core.Varchar:
		return "VARCHAR"
	case core.Boolean:
		return "SMALLINT"
	case core.Integer:
		return "BIGINT"
	case core.DateTime:
		return "TIMESTAMP(6)"
	case core.Monetary:
		return "NUMERIC(17,2)"
	case core.Decimal:
		return "DOUBLE PRECISION"
	case core.Binary:
		return "BYTEA"
	default:
		return "VARCHAR"
	}
}

func (sb postgresqlSQLBuilder) columnDefinition(pColumn core.Column) string {
	lColumnDef := pColumn.Name().String()
	if pColumn.Autoincrement() {
		lColumnDef += " BIGSERIAL"
	} else {
		if pColumn.Type() == core.Varchar {
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
	return sb.withGeneric().Insert(pTableName, pColumns)
}

func (sb postgresqlSQLBuilder) InsertWithValues(pTableName string, pColumnsValues core.ColumnsValues) (string, []any) {
	return sb.withGeneric().InsertWithValues(pTableName, pColumnsValues)
}

func (sb postgresqlSQLBuilder) Update(pTableName string, pColumns, pColumnsWhere []string) string {
	return sb.withGeneric().Update(pTableName, pColumns, pColumnsWhere)
}

func (sb postgresqlSQLBuilder) UpdateWithValues(pTableName string, pColumns, pColumnsWhere []string, pValues core.ColumnsValues) (string, []any) {
	return sb.withGeneric().UpdateWithValues(pTableName, pColumns, pColumnsWhere, pValues)
}

func (sb postgresqlSQLBuilder) Delete(pTableName string, pColumnsWhere []string) string {
	return sb.withGeneric().Delete(pTableName, pColumnsWhere)
}

func (sb postgresqlSQLBuilder) DeleteWithValues(pTableName string, pWhereColumnsValues core.ColumnsValues) (string, []any) {
	return sb.withGeneric().DeleteWithValues(pTableName, pWhereColumnsValues)
}

func (sb postgresqlSQLBuilder) Select(pTableName string, pColumns, pColumnsWhere []string) string {
	return sb.withGeneric().Select(pTableName, pColumns, pColumnsWhere)
}

func (sb postgresqlSQLBuilder) SelectWithValues(pTableName string, pColumns []string, pWhereColumnsValues core.ColumnsValues) (string, []any) {
	return sb.withGeneric().SelectWithValues(pTableName, pColumns, pWhereColumnsValues)
}

func (sb postgresqlSQLBuilder) SelectWithPagination(pSQLSelectCommand string, pRowsPerPage, pPageNumber uint) string {
	return sb.withGeneric().SelectWithPagination(pSQLSelectCommand, pRowsPerPage, pPageNumber)
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

func (sb postgresqlSQLBuilder) withGeneric() core.GenericSQLBuilder {
	if sb.generic == nil {
		temp := core.NewGenericSQLBuilder(sb.columnDefinition, sb.InteroperateSQLCommandWithNamedArgs)
		sb.generic = &temp
	}
	return *sb.generic
}
