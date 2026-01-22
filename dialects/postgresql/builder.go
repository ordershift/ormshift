package postgresql

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/internal"
	"github.com/ordershift/ormshift/schema"
)

type postgresqlSQLBuilder struct {
	generic *internal.GenericSQLBuilder
}

func (sb postgresqlSQLBuilder) CreateTable(pTable schema.Table) string {
	return sb.withGeneric().CreateTable(pTable)
}

func (sb postgresqlSQLBuilder) DropTable(pTableName schema.TableName) string {
	return sb.withGeneric().DropTable(pTableName)
}

func (sb postgresqlSQLBuilder) AlterTableAddColumn(pTableName schema.TableName, pColumn schema.Column) string {
	return sb.withGeneric().AlterTableAddColumn(pTableName, pColumn)
}

func (sb postgresqlSQLBuilder) AlterTableDropColumn(pTableName schema.TableName, pColumnName schema.ColumnName) string {
	return sb.withGeneric().AlterTableDropColumn(pTableName, pColumnName)
}

func (sb postgresqlSQLBuilder) ColumnTypeAsString(pColumnType schema.ColumnType) string {
	switch pColumnType {
	case schema.Varchar:
		return "VARCHAR"
	case schema.Boolean:
		return "SMALLINT"
	case schema.Integer:
		return "BIGINT"
	case schema.DateTime:
		return "TIMESTAMP(6)"
	case schema.Monetary:
		return "NUMERIC(17,2)"
	case schema.Decimal:
		return "DOUBLE PRECISION"
	case schema.Binary:
		return "BYTEA"
	default:
		return "VARCHAR"
	}
}

func (sb postgresqlSQLBuilder) columnDefinition(pColumn schema.Column) string {
	lColumnDef := pColumn.Name().String()
	if pColumn.Autoincrement() {
		lColumnDef += " BIGSERIAL"
	} else {
		if pColumn.Type() == schema.Varchar {
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

func (sb postgresqlSQLBuilder) InsertWithValues(pTableName string, pColumnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.withGeneric().InsertWithValues(pTableName, pColumnsValues)
}

func (sb postgresqlSQLBuilder) Update(pTableName string, pColumns, pColumnsWhere []string) string {
	return sb.withGeneric().Update(pTableName, pColumns, pColumnsWhere)
}

func (sb postgresqlSQLBuilder) UpdateWithValues(pTableName string, pColumns, pColumnsWhere []string, pValues ormshift.ColumnsValues) (string, []any) {
	return sb.withGeneric().UpdateWithValues(pTableName, pColumns, pColumnsWhere, pValues)
}

func (sb postgresqlSQLBuilder) Delete(pTableName string, pColumnsWhere []string) string {
	return sb.withGeneric().Delete(pTableName, pColumnsWhere)
}

func (sb postgresqlSQLBuilder) DeleteWithValues(pTableName string, pWhereColumnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.withGeneric().DeleteWithValues(pTableName, pWhereColumnsValues)
}

func (sb postgresqlSQLBuilder) Select(pTableName string, pColumns, pColumnsWhere []string) string {
	return sb.withGeneric().Select(pTableName, pColumns, pColumnsWhere)
}

func (sb postgresqlSQLBuilder) SelectWithValues(pTableName string, pColumns []string, pWhereColumnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.withGeneric().SelectWithValues(pTableName, pColumns, pWhereColumnsValues)
}

func (sb postgresqlSQLBuilder) SelectWithPagination(pSQLSelectCommand string, pRowsPerPage, pPageNumber uint) string {
	return sb.withGeneric().SelectWithPagination(pSQLSelectCommand, pRowsPerPage, pPageNumber)
}

func (sb postgresqlSQLBuilder) InteroperateSQLCommandWithNamedArgs(pSQLCommand string, pNamedArgs ...sql.NamedArg) (string, []any) {
	lSQLCommand := pSQLCommand
	lArgs := []any{}
	lIndexes := map[string]int{}
	for i, lParam := range pNamedArgs {
		lIndexes[strings.ToLower(lParam.Name)] = i + 1
		lBooleanValue, lIsBoolean := lParam.Value.(bool)
		if lIsBoolean {
			if lBooleanValue {
				lArgs = append(lArgs, int(1))
			} else {
				lArgs = append(lArgs, int(0))
			}
		} else {
			lArgs = append(lArgs, lParam.Value)
		}
	}
	lRegex := regexp.MustCompile(`@([a-zA-Z_][a-zA-Z0-9_]*)`)
	lSQLCommand = lRegex.ReplaceAllStringFunc(lSQLCommand, func(m string) string {
		lName := m[1:]
		if idx, ok := lIndexes[strings.ToLower(lName)]; ok {
			return fmt.Sprintf("$%d", idx)
		}
		return m
	})
	return lSQLCommand, lArgs
}

func (sb postgresqlSQLBuilder) withGeneric() internal.GenericSQLBuilder {
	if sb.generic == nil {
		temp := internal.NewGenericSQLBuilder(sb.columnDefinition, sb.InteroperateSQLCommandWithNamedArgs)
		sb.generic = &temp
	}
	return *sb.generic
}
