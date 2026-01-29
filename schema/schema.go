package schema

import (
	"database/sql"
	"errors"
	"slices"
	"strings"
)

type DBSchema struct {
	db                   *sql.DB
	tableNamesQuery      string
	columnTypesQueryFunc ColumnTypesQueryFunc
}

type ColumnTypesQueryFunc func(pTableName string) string

func NewDBSchema(
	pDB *sql.DB,
	pTableNamesQuery string,
	pColumnTypesQueryFunc ColumnTypesQueryFunc,
) (*DBSchema, error) {
	if pDB == nil {
		return nil, errors.New("sql.DB cannot be nil")
	}
	return &DBSchema{
		db:                   pDB,
		tableNamesQuery:      pTableNamesQuery,
		columnTypesQueryFunc: pColumnTypesQueryFunc,
	}, nil
}

func (s *DBSchema) HasTable(pTableName string) bool {
	tables, err := s.fetchTableNames()
	if err != nil {
		return false
	}
	return slices.ContainsFunc(tables, func(t string) bool {
		return strings.EqualFold(t, pTableName)
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
	tableName := ""
	for rows.Next() {
		err = rows.Scan(&tableName)
		if err != nil {
			return
		}
		tableNames = append(tableNames, tableName)
	}
	return
}

func (s *DBSchema) HasColumn(pTableName, pColumnName string) bool {
	columnTypes, err := s.fetchColumnTypes(pTableName)
	if err != nil {
		return false
	}
	return slices.ContainsFunc(columnTypes, func(ct *sql.ColumnType) bool {
		return strings.EqualFold(ct.Name(), pColumnName)
	})
}

func (s *DBSchema) fetchColumnTypes(pTableName string) (columnTypes []*sql.ColumnType, err error) {
	rows, err := s.db.Query(s.columnTypesQueryFunc(pTableName))
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
