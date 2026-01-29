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

func (sb *postgresqlBuilder) CreateTable(table schema.Table) string {
	return sb.generic.CreateTable(table)
}

func (sb *postgresqlBuilder) DropTable(tableName string) string {
	return sb.generic.DropTable(tableName)
}

func (sb *postgresqlBuilder) AlterTableAddColumn(tableName string, column schema.Column) string {
	return sb.generic.AlterTableAddColumn(tableName, column)
}

func (sb *postgresqlBuilder) AlterTableDropColumn(tableName, columnName string) string {
	return sb.generic.AlterTableDropColumn(tableName, columnName)
}

func (sb *postgresqlBuilder) ColumnTypeAsString(columnType schema.ColumnType) string {
	switch columnType {
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

func (sb *postgresqlBuilder) columnDefinition(column schema.Column) string {
	columnDef := sb.QuoteIdentifier(column.Name())
	if column.AutoIncrement() {
		columnDef += " BIGSERIAL"
	} else {
		if column.Type() == schema.Varchar {
			columnDef += fmt.Sprintf(" %s(%d)", sb.ColumnTypeAsString(column.Type()), column.Size())
		} else {
			columnDef += fmt.Sprintf(" %s", sb.ColumnTypeAsString(column.Type()))
		}
	}
	if column.NotNull() {
		columnDef += " NOT NULL"
	}
	return columnDef
}

func (sb *postgresqlBuilder) Insert(tableName string, columns []string) string {
	return sb.generic.Insert(tableName, columns)
}

func (sb *postgresqlBuilder) InsertWithValues(tableName string, columnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.InsertWithValues(tableName, columnsValues)
}

func (sb *postgresqlBuilder) Update(tableName string, columns, columnsWhere []string) string {
	return sb.generic.Update(tableName, columns, columnsWhere)
}

func (sb *postgresqlBuilder) UpdateWithValues(tableName string, columns, columnsWhere []string, values ormshift.ColumnsValues) (string, []any) {
	return sb.generic.UpdateWithValues(tableName, columns, columnsWhere, values)
}

func (sb *postgresqlBuilder) Delete(tableName string, columnsWhere []string) string {
	return sb.generic.Delete(tableName, columnsWhere)
}

func (sb *postgresqlBuilder) DeleteWithValues(tableName string, whereColumnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.DeleteWithValues(tableName, whereColumnsValues)
}

func (sb *postgresqlBuilder) Select(tableName string, columns, columnsWhere []string) string {
	return sb.generic.Select(tableName, columns, columnsWhere)
}

func (sb *postgresqlBuilder) SelectWithValues(tableName string, columns []string, whereColumnsValues ormshift.ColumnsValues) (string, []any) {
	return sb.generic.SelectWithValues(tableName, columns, whereColumnsValues)
}

func (sb *postgresqlBuilder) SelectWithPagination(sqlSelectCommand string, rowsPerPage, pageNumber uint) string {
	return sb.generic.SelectWithPagination(sqlSelectCommand, rowsPerPage, pageNumber)
}

func (sb *postgresqlBuilder) QuoteIdentifier(identifier string) string {
	return sb.generic.QuoteIdentifier(identifier)
}

func (sb *postgresqlBuilder) InteroperateSQLCommandWithNamedArgs(sqlCommand string, namedArgs ...sql.NamedArg) (string, []any) {
	args := []any{}
	indexes := map[string]int{}
	for i, param := range namedArgs {
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
