package core

import (
	"database/sql"
	"slices"
)

// DDLSQLBuilder cria comandos DDL (Data Definition Language) em SQL para definição de dados em SGBD.
// Os principais comandos DDL são CREATE, ALTER e DROP
type DDLSQLBuilder interface {
	CreateTable(pTable Table) string
	DropTable(pTableName TableName) string
	AlterTableAddColumn(pTableName TableName, pColumn Column) string
	AlterTableDropColumn(pTableName TableName, pColumnName ColumnName) string
	ColumnTypeAsString(pColumnType ColumnType) string
}

type ColumnsValues map[string]any

// ToNamedArgs transforms ColumnsValues to an ordened by name sql.NamedArg array. Eg.:
//
//	lColumnsValues := ColumnsValues{"id": 5, "sku": "ZTX-9000", "is_simple": true}
//	lNamedArgs := lColumnsValues.ToNamedArgs()
//	//lNamedArgs == []sql.NamedArg{{Name: "id", Value: 5},{Name: "is_simple", Value: true},{Name: "sku", Value: "ZTX-9000"}}
func (cv ColumnsValues) ToNamedArgs() []sql.NamedArg {
	lNamedArgs := []sql.NamedArg{}
	for c, v := range cv {
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

func (cv ColumnsValues) ToColumns() []string {
	lColumns := []string{}
	for c := range cv {
		lColumns = append(lColumns, c)
	}
	slices.Sort(lColumns)
	return lColumns
}

// DMLSQLBuilder cria comandos DML (Data Manipulation Language) em SQL para manipulação de dados em SGBD.
// Os principais comandos DML são INSERT, UPDATE, DELETE e SELECT
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

	// InteroperateSQLCommandWithNamedArgs serve para padronizar os comandos SQL
	// com parâmetros de acordo com o driver de conexão com banco de dados.
	//
	// Atualmente está preparado para os drivers Postgresql, SQLite e SQLServer.
	//
	// Os drivers para SQLite e SQLServer aceitam NamedArgs como parâmetros do
	// comando SQL, mas o driver para PostgreSQL aceita apenas parâmetros numerados
	// com prefixo $ (exemplo $1, $2, $3). Esta função transforma o comando SQL
	// com NamedArgs para o formato suportado pelo respectivo driver. Exemplo:
	//
	//	pSQLCommand := "select * from user where id = @id"
	//	pNamedArg := sql.Named("id", 123)
	//
	//	//POSTGRESQL:
	//	q, p := sqlbuilder.InteroperateSQLCommandWithNamedArgs(pSQLCommand, pNamedArg)
	//	//q == "select * from user where id = $1"
	//	//p == 123
	//
	//	//SQLITE:
	//	q, p = sqlbuilder.InteroperateSQLCommandWithNamedArgs(pSQLCommand, pNamedArg)
	//	//q == "select * from user where id = @id"
	//	//p == sql.Named("id", 123)
	//
	//	//SQLSERVER:
	//	q, p = sqlbuilder.InteroperateSQLCommandWithNamedArgs(pSQLCommand, pNamedArg)
	//	//q == "select * from user where id = @id"
	//	//p == sql.Named("id", 123)
	//
	// Futuramente, para suportar o driver para MySQL, que espera interrogação nos parâmetros
	//
	//	//MYSQL:
	//	q, p = sqlbuilder.InteroperateSQLCommandWithNamedArgs(pSQLCommand, pNamedArg)
	//	//q == "select * from user where id = ?"
	//	//p == 123
	InteroperateSQLCommandWithNamedArgs(pSQLCommand string, pNamedArgs ...sql.NamedArg) (string, []any)
}

type SQLBuilder interface {
	DDLSQLBuilder
	DMLSQLBuilder
}
