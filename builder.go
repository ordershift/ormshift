package ormshift

import (
	"database/sql"
	"errors"
	"slices"

	"github.com/ordershift/ormshift/errs"
	"github.com/ordershift/ormshift/schema"
)

// DDLSQLBuilder creates DDL (Data Definition Language) SQL commands for defining schema in DBMS.
type DDLSQLBuilder interface {
	CreateTable(table schema.Table) string
	DropTable(table string) string
	AlterTableAddColumn(table string, column schema.Column) string
	AlterTableDropColumn(table, column string) string
	ColumnTypeAsString(columnType schema.ColumnType) string
}

// ColumnsValues represents a mapping between column names and their corresponding values.
type ColumnsValues map[string]any

func NewColumnsValues(columns []string, values []any) (*ColumnsValues, error) {
	if len(columns) != len(values) {
		return nil, failedToGetColumnsValues(errors.New("columns len must be equal to values len"))
	}
	cv := ColumnsValues{}
	for i, c := range columns {
		cv[c] = values[i]
	}
	return &cv, nil
}

func failedToGetColumnsValues(err error) error {
	return errs.FailedTo("get columns values", err)
}

// ToNamedArgs transforms ColumnsValues to a sql.NamedArg array ordered by name, e.g.:
//
//	values := ColumnsValues{"id": 5, "sku": "ZTX-9000", "is_simple": true}
//	args := values.ToNamedArgs()
//	//args == []sql.NamedArg{{Name: "id", Value: 5},{Name: "is_simple", Value: true},{Name: "sku", Value: "ZTX-9000"}}
func (cv *ColumnsValues) ToNamedArgs() []sql.NamedArg {
	args := []sql.NamedArg{}
	for c, v := range *cv {
		args = append(args, sql.Named(c, v))
	}
	slices.SortFunc(args, func(a, b sql.NamedArg) int {
		if a.Name < b.Name {
			return -1
		}
		return 1
	})
	return args
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
	Insert(table string, columns []string) string
	InsertWithValues(table string, values ColumnsValues) (string, []any)
	Update(table string, columns, where []string) string
	UpdateWithValues(table string, columns, where []string, values ColumnsValues) (string, []any)
	Delete(table string, where []string) string
	DeleteWithValues(table string, where ColumnsValues) (string, []any)
	Select(table string, columns, where []string) string
	SelectWithValues(table string, columns []string, where ColumnsValues) (string, []any)
	SelectWithPagination(sql string, size, number uint) string

	// InteroperateSQLCommandWithNamedArgs acts as a SQL command translator that standardizes SQL commands according to the database driver being used e.g.,
	//
	//	sql := "select * from user where id = @id"
	//	namedArg := sql.Named("id", 123)
	//
	// PostgreSQL:
	//	q, p := sqlbuilder.InteroperateSQLCommandWithNamedArgs(sql, namedArg)
	//	//q == "select * from user where id = $1"
	//	//p == 123
	//
	// SQLite:
	//	q, p = sqlbuilder.InteroperateSQLCommandWithNamedArgs(sql, namedArg)
	//	//q == "select * from user where id = @id"
	//	//p == sql.Named("id", 123)
	//
	// SQL Server:
	//	q, p = sqlbuilder.InteroperateSQLCommandWithNamedArgs(sql, namedArg)
	//	//q == "select * from user where id = @id"
	//	//p == sql.Named("id", 123)
	//
	// MySQL (not yet supported, expects question marks in parameters):
	//
	//	q, p = sqlbuilder.InteroperateSQLCommandWithNamedArgs(sql, namedArg)
	//	//q == "select * from user where id = ?"
	//	//p == 123
	InteroperateSQLCommandWithNamedArgs(sql string, args ...sql.NamedArg) (string, []any)
}

type SQLBuilder interface {
	DDLSQLBuilder
	DMLSQLBuilder
	QuoteIdentifier(identifier string) string
}
