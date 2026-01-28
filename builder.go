package ormshift

import (
	"database/sql"
	"slices"

	"github.com/ordershift/ormshift/schema"
)

// DDSQLBuilder creates DDL (Data Definition Language) SQL commands for defining schema in DBMS.
type DDLSQLBuilder interface {
	CreateTable(pTable schema.Table) string
	DropTable(pTableName string) string
	AlterTableAddColumn(pTableName string, pColumn schema.Column) string
	AlterTableDropColumn(pTableName, pColumnName string) string
	ColumnTypeAsString(pColumnType schema.ColumnType) string
}

// ColumnsValues represents a mapping between column names and their corresponding values.
type ColumnsValues map[string]any

// ToNamedArgs transforms ColumnsValues to a sql.NamedArg array ordered by name, e.g.:
//
//	lColumnsValues := ColumnsValues{"id": 5, "sku": "ZTX-9000", "is_simple": true}
//	lNamedArgs := lColumnsValues.ToNamedArgs()
//	//lNamedArgs == []sql.NamedArg{{Name: "id", Value: 5},{Name: "is_simple", Value: true},{Name: "sku", Value: "ZTX-9000"}}
func (cv *ColumnsValues) ToNamedArgs() []sql.NamedArg {
	lNamedArgs := []sql.NamedArg{}
	for c, v := range *cv {
		lNamedArgs = append(lNamedArgs, sql.Named(c, v))
	}
	slices.SortFunc(lNamedArgs, func(a, b sql.NamedArg) int {
		if a.Name < b.Name {
			return -1
		}
		return 1
	})
	return lNamedArgs
}

// ToColumns returns the column names from ColumnsValues as a string array ordered by name, e.g.:
func (cv *ColumnsValues) ToColumns() []string {
	lColumns := []string{}
	for c := range *cv {
		lColumns = append(lColumns, c)
	}
	slices.Sort(lColumns)
	return lColumns
}

// DMLSQLBuilder creates DML (Data Manipulation Language) SQL commands for manipulating data in DBMS.
type DMLSQLBuilder interface {
	Insert(pTableName string, pColumns []string) string
	InsertWithValues(pTableName string, pColumnsValues ColumnsValues) (string, []any)
	Update(pTableName string, pColumns, pColumnsWhere []string) string
	UpdateWithValues(pTableName string, pColumns, pColumnsWhere []string, pValues ColumnsValues) (string, []any)
	Delete(pTableName string, pColumnsWhere []string) string
	DeleteWithValues(pTableName string, pWhereColumnsValues ColumnsValues) (string, []any)
	Select(pTableName string, pColumns, pColumnsWhere []string) string
	SelectWithValues(pTableName string, pColumns []string, pWhereColumnsValues ColumnsValues) (string, []any)
	SelectWithPagination(pSQLSelectCommand string, pRowsPerPage, pPageNumber uint) string

	// InteroperateSQLCommandWithNamedArgs acts as a SQL command translator that standardizes SQL commands according to the database driver being used e.g.,
	//
	//	pSQLCommand := "select * from user where id = @id"
	//	pNamedArg := sql.Named("id", 123)
	//
	// PostgreSQL:
	//	q, p := sqlbuilder.InteroperateSQLCommandWithNamedArgs(pSQLCommand, pNamedArg)
	//	//q == "select * from user where id = $1"
	//	//p == 123
	//
	// SQLite:
	//	q, p = sqlbuilder.InteroperateSQLCommandWithNamedArgs(pSQLCommand, pNamedArg)
	//	//q == "select * from user where id = @id"
	//	//p == sql.Named("id", 123)
	//
	// SQL Server:
	//	q, p = sqlbuilder.InteroperateSQLCommandWithNamedArgs(pSQLCommand, pNamedArg)
	//	//q == "select * from user where id = @id"
	//	//p == sql.Named("id", 123)
	//
	// MySQL (not yet supported, expects question marks in parameters):
	//
	//	q, p = sqlbuilder.InteroperateSQLCommandWithNamedArgs(pSQLCommand, pNamedArg)
	//	//q == "select * from user where id = ?"
	//	//p == 123
	InteroperateSQLCommandWithNamedArgs(pSQLCommand string, pNamedArgs ...sql.NamedArg) (string, []any)
}

type SQLBuilder interface {
	DDLSQLBuilder
	DMLSQLBuilder
	QuoteIdentifier(pIdentifier string) string
}
