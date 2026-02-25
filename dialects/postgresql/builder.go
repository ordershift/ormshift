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

func (sb *postgresqlBuilder) DropTable(table string) string {
	return sb.generic.DropTable(table)
}

func (sb *postgresqlBuilder) AlterTableAddColumn(table string, column schema.Column) string {
	return sb.generic.AlterTableAddColumn(table, column)
}

func (sb *postgresqlBuilder) AlterTableDropColumn(table, column string) string {
	return sb.generic.AlterTableDropColumn(table, column)
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
	case schema.DateTimeOffSet:
		return "TIMESTAMPTZ(6)"
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

func (sb *postgresqlBuilder) Insert(table string, columns []string) string {
	return sb.generic.Insert(table, columns)
}

func (sb *postgresqlBuilder) InsertWithValues(table string, values ormshift.ColumnsValues) (string, []any) {
	return sb.generic.InsertWithValues(table, values)
}

func (sb *postgresqlBuilder) Update(table string, columns, where []string) string {
	return sb.generic.Update(table, columns, where)
}

func (sb *postgresqlBuilder) UpdateWithValues(table string, columns, where []string, values ormshift.ColumnsValues) (string, []any) {
	return sb.generic.UpdateWithValues(table, columns, where, values)
}

func (sb *postgresqlBuilder) Delete(table string, where []string) string {
	return sb.generic.Delete(table, where)
}

func (sb *postgresqlBuilder) DeleteWithValues(table string, where ormshift.ColumnsValues) (string, []any) {
	return sb.generic.DeleteWithValues(table, where)
}

func (sb *postgresqlBuilder) Select(table string, columns, where []string) string {
	return sb.generic.Select(table, columns, where)
}

func (sb *postgresqlBuilder) SelectWithValues(table string, columns []string, where ormshift.ColumnsValues) (string, []any) {
	return sb.generic.SelectWithValues(table, columns, where)
}

func (sb *postgresqlBuilder) SelectWithPagination(sql string, size, number uint) string {
	return sb.generic.SelectWithPagination(sql, size, number)
}

func (sb *postgresqlBuilder) QuoteIdentifier(identifier string) string {
	return sb.generic.QuoteIdentifier(identifier)
}

func (sb *postgresqlBuilder) InteroperateSQLCommandWithNamedArgs(sql string, args ...sql.NamedArg) (string, []any) {
	a := []any{}
	indexes := map[string]int{}
	for i, param := range args {
		indexes[strings.ToLower(param.Name)] = i + 1
		booleanValue, isBoolean := param.Value.(bool)
		if isBoolean {
			if booleanValue {
				a = append(a, int(1))
			} else {
				a = append(a, int(0))
			}
		} else {
			a = append(a, param.Value)
		}
	}
	regex := regexp.MustCompile(`@([a-zA-Z_][a-zA-Z0-9_]*)`)
	sql = regex.ReplaceAllStringFunc(sql, func(m string) string {
		name := m[1:]
		if idx, ok := indexes[strings.ToLower(name)]; ok {
			return fmt.Sprintf("$%d", idx)
		}
		return m
	})
	return sql, a
}
