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
	sb.generic = internal.NewGenericSQLBuilder(sb.columnDefinition, nil, sb.InteroperateSQLCommandWithNamedArgs)
	return &sb
}

func (sb *postgresqlBuilder) CreateTable(pTable schema.Table) string {
	return sb.generic.CreateTable(pTable)
}

func (sb *postgresqlBuilder) DropTable(pTableName string) string {
	return sb.generic.DropTable(pTableName)
}

func (sb *postgresqlBuilder) AlterTableAddColumn(pTableName string, pColumn schema.Column) string {
	return sb.generic.AlterTableAddColumn(pTableName, pColumn)
}

func (sb *postgresqlBuilder) AlterTableDropColumn(pTableName, pColumnName string) string {
	return sb.generic.AlterTableDropColumn(pTableName, pColumnName)
}

func (sb *postgresqlBuilder) ColumnTypeAsString(pColumnType schema.ColumnType) string {
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

func (sb *postgresqlBuilder) columnDefinition(pColumn schema.Column) string {
	columnDef := sb.QuoteIdentifier(pColumn.Name())
	if pColumn.AutoIncrement() {
		columnDef += " BIGSERIAL"
	} else {
		if pColumn.Type() == schema.Varchar {
			columnDef += fmt.Sprintf(" %s(%d)", sb.ColumnTypeAsString(pColumn.Type()), pColumn.Size())
		} else {
			columnDef += fmt.Sprintf(" %s", sb.ColumnTypeAsString(pColumn.Type()))
		}
	}
	if pColumn.NotNull() {
		columnDef += " NOT NULL"
	}
	return columnDef
}

func (sb *postgresqlBuilder) Insert(pTableName string, pColumns []string) string {
	return sb.generic.Insert(pTableName, pColumns)
}

func (sb *postgresqlBuilder) InsertWithValues(pTableName string, pColumnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.InsertWithValues(pTableName, pColumnsValues)
}

func (sb *postgresqlBuilder) Update(pTableName string, pColumns, pColumnsWhere []string) string {
	return sb.generic.Update(pTableName, pColumns, pColumnsWhere)
}

func (sb *postgresqlBuilder) UpdateWithValues(pTableName string, pColumns, pColumnsWhere []string, pValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.UpdateWithValues(pTableName, pColumns, pColumnsWhere, pValues)
}

func (sb *postgresqlBuilder) Delete(pTableName string, pColumnsWhere []string) string {
	return sb.generic.Delete(pTableName, pColumnsWhere)
}

func (sb *postgresqlBuilder) DeleteWithValues(pTableName string, pWhereColumnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.DeleteWithValues(pTableName, pWhereColumnsValues)
}

func (sb *postgresqlBuilder) Select(pTableName string, pColumns, pColumnsWhere []string) string {
	return sb.generic.Select(pTableName, pColumns, pColumnsWhere)
}

func (sb *postgresqlBuilder) SelectWithValues(pTableName string, pColumns []string, pWhereColumnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.SelectWithValues(pTableName, pColumns, pWhereColumnsValues)
}

func (sb *postgresqlBuilder) SelectWithPagination(pSQLSelectCommand string, pRowsPerPage, pPageNumber uint) string {
	return sb.generic.SelectWithPagination(pSQLSelectCommand, pRowsPerPage, pPageNumber)
}

func (sb *postgresqlBuilder) QuoteIdentifier(pIdentifier string) string {
	return sb.generic.QuoteIdentifier(pIdentifier)
}

func (sb *postgresqlBuilder) InteroperateSQLCommandWithNamedArgs(pSQLCommand string, pNamedArgs ...sql.NamedArg) (string, []any) {
	sqlCommand := pSQLCommand
	args := []any{}
	indexes := map[string]int{}
	for i, param := range pNamedArgs {
		indexes[strings.ToLower(param.Name)] = i + 1
		booleanValue, isBoolean := param.Value.(bool)
		if isBoolean {
			if booleanValue {
				args = append(args, int(1))
			} else {
				args = append(args, int(0))
			}
		} else {
			args = append(args, param.Value)
		}
	}
	regex := regexp.MustCompile(`@([a-zA-Z_][a-zA-Z0-9_]*)`)
	sqlCommand = regex.ReplaceAllStringFunc(sqlCommand, func(m string) string {
		name := m[1:]
		if idx, ok := indexes[strings.ToLower(name)]; ok {
			return fmt.Sprintf("$%d", idx)
		}
		return m
	})
	return sqlCommand, args
}
