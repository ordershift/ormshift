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

type postgresqlBuilder struct {
	generic ormshift.SQLBuilder
}

func newPostgreSQLBuilder() ormshift.SQLBuilder {
	sb := postgresqlBuilder{}
	sb.generic = internal.NewGenericSQLBuilder(sb.columnDefinition, QuoteIdentifier, sb.InteroperateSQLCommandWithNamedArgs)
	return sb
}

func (sb postgresqlBuilder) CreateTable(pTable schema.Table) string {
	return sb.generic.CreateTable(pTable)
}

func (sb postgresqlBuilder) DropTable(pTableName string) string {
	return sb.generic.DropTable(pTableName)
}

func (sb postgresqlBuilder) AlterTableAddColumn(pTableName string, pColumn schema.Column) string {
	return sb.generic.AlterTableAddColumn(pTableName, pColumn)
}

func (sb postgresqlBuilder) AlterTableDropColumn(pTableName string, pColumnName string) string {
	return sb.generic.AlterTableDropColumn(pTableName, pColumnName)
}

func (sb postgresqlBuilder) ColumnTypeAsString(pColumnType schema.ColumnType) string {
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

func (sb postgresqlBuilder) columnDefinition(pColumn schema.Column) string {
	lColumnDef := QuoteIdentifier(pColumn.Name())
	if pColumn.AutoIncrement() {
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

func (sb postgresqlBuilder) Insert(pTableName string, pColumns []string) string {
	return sb.generic.Insert(pTableName, pColumns)
}

func (sb postgresqlBuilder) InsertWithValues(pTableName string, pColumnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.InsertWithValues(pTableName, pColumnsValues)
}

func (sb postgresqlBuilder) Update(pTableName string, pColumns, pColumnsWhere []string) string {
	return sb.generic.Update(pTableName, pColumns, pColumnsWhere)
}

func (sb postgresqlBuilder) UpdateWithValues(pTableName string, pColumns, pColumnsWhere []string, pValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.UpdateWithValues(pTableName, pColumns, pColumnsWhere, pValues)
}

func (sb postgresqlBuilder) Delete(pTableName string, pColumnsWhere []string) string {
	return sb.generic.Delete(pTableName, pColumnsWhere)
}

func (sb postgresqlBuilder) DeleteWithValues(pTableName string, pWhereColumnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.DeleteWithValues(pTableName, pWhereColumnsValues)
}

func (sb postgresqlBuilder) Select(pTableName string, pColumns, pColumnsWhere []string) string {
	return sb.generic.Select(pTableName, pColumns, pColumnsWhere)
}

func (sb postgresqlBuilder) SelectWithValues(pTableName string, pColumns []string, pWhereColumnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.SelectWithValues(pTableName, pColumns, pWhereColumnsValues)
}

func (sb postgresqlBuilder) SelectWithPagination(pSQLSelectCommand string, pRowsPerPage, pPageNumber uint) string {
	return sb.generic.SelectWithPagination(pSQLSelectCommand, pRowsPerPage, pPageNumber)
}

func (sb postgresqlBuilder) InteroperateSQLCommandWithNamedArgs(pSQLCommand string, pNamedArgs ...sql.NamedArg) (string, []any) {
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

func QuoteIdentifier(pIdentifier string) string {
	// PostgreSQL uses double quotes: "identifier"
	// Escape rule: double quote becomes two double quotes
	// Example: users -> "users", table"name -> "table""name"
	pIdentifier = strings.ReplaceAll(pIdentifier, `"`, `""`)
	return fmt.Sprintf(`"%s"`, pIdentifier)
}
