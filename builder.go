package ormshift

import (
	"database/sql"
	"slices"

	"github.com/ordershift/ormshift/schema"
)

// DDSQLBuilder creates DDL (Data Definition Language) SQL commands for defining schema in DBMS.
type DDLSQLBuilder interface {
	CreateTable(table schema.Table) string
	DropTable(tableName string) string
	AlterTableAddColumn(tableName string, column schema.Column) string
	AlterTableDropColumn(tableName, columnName string) string
	ColumnTypeAsString(columnType schema.ColumnType) string
}

// ColumnsValues represents a mapping between column names and their corresponding values.
type ColumnsValues map[string]any

// ToNamedArgs transforms ColumnsValues to a sql.NamedArg array ordered by name, e.g.:
//
//	columnsValues := ColumnsValues{"id": 5, "sku": "ZTX-9000", "is_simple": true}
//	namedArgs := columnsValues.ToNamedArgs()
//	//namedArgs == []sql.NamedArg{{Name: "id", Value: 5},{Name: "is_simple", Value: true},{Name: "sku", Value: "ZTX-9000"}}
func (cv *ColumnsValues) ToNamedArgs() []sql.NamedArg {
	namedArgs := []sql.NamedArg{}
	for c, v := range *cv {
		namedArgs = append(namedArgs, sql.Named(c, v))
	}
	slices.SortFunc(namedArgs, func(a, b sql.NamedArg) int {
		if a.Name < b.Name {
			return -1
		}
		return 1
	})
	return namedArgs
}

// ToColumns returns the column names from ColumnsValues as a string array ordered by name, e.g.:
func (cv *ColumnsValues) ToColumns() []string {
	columns := []string{}
	for c := range *cv {
		columns = append(columns, c)
	}
	slices.Sort(columns)
	return columns
}

// DMLSQLBuilder creates DML (Data Manipulation Language) SQL commands for manipulating data in DBMS.
type DMLSQLBuilder interface {
	Insert(tableName string, columns []string) string
	InsertWithValues(tableName string, columnsValues ColumnsValues) (string, []any)
	Update(tableName string, columns, columnsWhere []string) string
	UpdateWithValues(tableName string, columns, columnsWhere []string, values ColumnsValues) (string, []any)
	Delete(tableName string, columnsWhere []string) string
	DeleteWithValues(tableName string, whereColumnsValues ColumnsValues) (string, []any)
	Select(tableName string, columns, columnsWhere []string) string
	SelectWithValues(tableName string, columns []string, whereColumnsValues ColumnsValues) (string, []any)
	SelectWithPagination(sqlSelectCommand string, rowsPerPage, pageNumber uint) string

	// InteroperateSQLCommandWithNamedArgs acts as a SQL command translator that standardizes SQL commands according to the database driver being used e.g.,
	//
	//	sqlCommand := "select * from user where id = @id"
	//	namedArg := sql.Named("id", 123)
	//
	// PostgreSQL:
	//	q, p := sqlbuilder.InteroperateSQLCommandWithNamedArgs(sqlCommand, namedArg)
	//	//q == "select * from user where id = $1"
	//	//p == 123
	//
	// SQLite:
	//	q, p = sqlbuilder.InteroperateSQLCommandWithNamedArgs(sqlCommand, namedArg)
	//	//q == "select * from user where id = @id"
	//	//p == sql.Named("id", 123)
	//
	// SQL Server:
	//	q, p = sqlbuilder.InteroperateSQLCommandWithNamedArgs(sqlCommand, namedArg)
	//	//q == "select * from user where id = @id"
	//	//p == sql.Named("id", 123)
	//
	// MySQL (not yet supported, expects question marks in parameters):
	//
	//	q, p = sqlbuilder.InteroperateSQLCommandWithNamedArgs(sqlCommand, namedArg)
	//	//q == "select * from user where id = ?"
	//	//p == 123
	InteroperateSQLCommandWithNamedArgs(sqlCommand string, namedArgs ...sql.NamedArg) (string, []any)
}

type SQLBuilder interface {
	DDLSQLBuilder
	DMLSQLBuilder
	QuoteIdentifier(identifier string) string
}
