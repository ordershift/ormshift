package schema

import (
	"database/sql"
	"slices"
	"strings"

	"github.com/ordershift/ormshift/errs"
)

type DBSchema struct {
	db                   *sql.DB
	tableNamesQuery      string
	columnTypesQueryFunc ColumnTypesQueryFunc
}

type ColumnTypesQueryFunc func(table string) string

func NewDBSchema(
	db *sql.DB,
	tableNamesQuery string,
	columnTypesQueryFunc ColumnTypesQueryFunc,
) (*DBSchema, error) {
	if db == nil {
		err := errs.Nil("db")
		return nil, failedToGetDBSchema(err)
	}
	return &DBSchema{
		db:                   db,
		tableNamesQuery:      tableNamesQuery,
		columnTypesQueryFunc: columnTypesQueryFunc,
	}, nil
}

func failedToGetDBSchema(err error) error {
	return errs.FailedTo("get db schema", err)
}

func (s *DBSchema) HasTable(table string) bool {
	tables, err := s.fetchTableNames()
	if err != nil {
		return false
	}
	return slices.ContainsFunc(tables, func(t string) bool {
		return strings.EqualFold(t, table)
	})
}

func (s *DBSchema) fetchTableNames() (tableNames []string, err error) {
	rows, err := s.db.Query(s.tableNamesQuery)
	if err != nil {
		return
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	table := ""
	for rows.Next() {
		err = rows.Scan(&table)
		if err != nil {
			return
		}
		tableNames = append(tableNames, table)
	}
	return
}

func (s *DBSchema) HasColumn(table, column string) bool {
	columnTypes, err := s.fetchColumnTypes(table)
	if err != nil {
		return false
	}
	return slices.ContainsFunc(columnTypes, func(ct *sql.ColumnType) bool {
		return strings.EqualFold(ct.Name(), column)
	})
}

func (s *DBSchema) fetchColumnTypes(table string) (columnTypes []*sql.ColumnType, err error) {
	rows, err := s.db.Query(s.columnTypesQueryFunc(table))
	if err != nil {
		return
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	columnTypes, err = rows.ColumnTypes()
	return
}
