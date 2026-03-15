package sqlite

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/internal"
	"github.com/ordershift/ormshift/schema"
)

type sqliteBuilder struct {
	generic ormshift.SQLBuilder
}

func newSQLiteBuilder() ormshift.SQLBuilder {
	sb := sqliteBuilder{}
	sb.generic = internal.NewGenericSQLBuilder(sb.columnDefinition, nil, nil)
	return &sb
}

func (sb *sqliteBuilder) CreateTable(table schema.Table) string {
	pk := table.PK()
	useInlinePK := pk != nil && len(pk.Columns()) == 1
	var pkColName string
	if useInlinePK {
		pkColName = pk.Columns()[0]
		for _, col := range table.Columns() {
			if strings.EqualFold(col.Name(), pkColName) {
				if col.Type() != schema.Integer || !col.AutoIncrement() {
					useInlinePK = false
				}
				break
			}
		}
	}

	var parts []string
	for _, column := range table.Columns() {
		if useInlinePK && strings.EqualFold(column.Name(), pkColName) {
			parts = append(parts, fmt.Sprintf("%s INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT", sb.QuoteIdentifier(column.Name())))
		} else {
			parts = append(parts, sb.columnDefinition(column))
		}
	}

	if pk != nil && !useInlinePK {
		pkColumns := ""
		for i, col := range pk.Columns() {
			if i > 0 {
				pkColumns += ","
			}
			pkColumns += sb.QuoteIdentifier(col)
		}
		parts = append(parts, fmt.Sprintf("CONSTRAINT %s PRIMARY KEY (%s)", sb.QuoteIdentifier(pk.Name()), pkColumns))
	}

	return fmt.Sprintf("CREATE TABLE %s (%s);", sb.QuoteIdentifier(table.Name()), strings.Join(parts, ","))
}

func (sb *sqliteBuilder) DropTable(table string) string {
	return sb.generic.DropTable(table)
}

func (sb *sqliteBuilder) AlterTableAddColumn(table string, column schema.Column) string {
	return sb.generic.AlterTableAddColumn(table, column)
}

func (sb *sqliteBuilder) AlterTableDropColumn(table, column string) string {
	return sb.generic.AlterTableDropColumn(table, column)
}

func (sb *sqliteBuilder) ColumnTypeAsString(columnType schema.ColumnType) string {
	switch columnType {
	case schema.Varchar:
		return "TEXT"
	case schema.Boolean:
		return "INTEGER"
	case schema.Integer:
		return "INTEGER"
	case schema.DateTime:
		return "DATETIME"
	case schema.Monetary:
		return "REAL"
	case schema.Decimal:
		return "REAL"
	case schema.Binary:
		return "BLOB"
	case schema.DateTimeOffSet:
		return "DATETIME"
	default:
		return "TEXT"
	}
}

func (sb *sqliteBuilder) columnDefinition(column schema.Column) string {
	columnDef := fmt.Sprintf("%s %s", sb.QuoteIdentifier(column.Name()), sb.ColumnTypeAsString(column.Type()))
	if column.NotNull() {
		columnDef += " NOT NULL"
	}
	if column.AutoIncrement() {
		columnDef += " AUTOINCREMENT"
	}
	return columnDef
}

func (sb *sqliteBuilder) Insert(table string, columns []string) string {
	return sb.generic.Insert(table, columns)
}

func (sb *sqliteBuilder) InsertWithValues(table string, values ormshift.ColumnsValues) (string, []any) {
	return sb.generic.InsertWithValues(table, values)
}

func (sb *sqliteBuilder) Update(table string, columns, where []string) string {
	return sb.generic.Update(table, columns, where)
}

func (sb *sqliteBuilder) UpdateWithValues(table string, columns, where []string, values ormshift.ColumnsValues) (string, []any) {
	return sb.generic.UpdateWithValues(table, columns, where, values)
}

func (sb *sqliteBuilder) Delete(table string, where []string) string {
	return sb.generic.Delete(table, where)
}

func (sb *sqliteBuilder) DeleteWithValues(table string, where ormshift.ColumnsValues) (string, []any) {
	return sb.generic.DeleteWithValues(table, where)
}

func (sb *sqliteBuilder) Select(table string, columns, where []string) string {
	return sb.generic.Select(table, columns, where)
}

func (sb *sqliteBuilder) SelectWithValues(table string, columns []string, where ormshift.ColumnsValues) (string, []any) {
	return sb.generic.SelectWithValues(table, columns, where)
}

func (sb *sqliteBuilder) SelectWithPagination(sql string, size, number uint) string {
	return sb.generic.SelectWithPagination(sql, size, number)
}

func (sb *sqliteBuilder) QuoteIdentifier(identifier string) string {
	return sb.generic.QuoteIdentifier(identifier)
}

func (sb *sqliteBuilder) InteroperateSQLCommandWithNamedArgs(sql string, args ...sql.NamedArg) (string, []any) {
	return sb.generic.InteroperateSQLCommandWithNamedArgs(sql, args...)
}
