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
	lTables, lError := s.fetchTableNames()
	if lError != nil {
		return false
	}
	return slices.ContainsFunc(lTables, func(t string) bool {
		return strings.EqualFold(t, pTableName)
	})
}

func (s *DBSchema) fetchTableNames() (rTableNames []string, rError error) {
	lRows, rError := s.db.Query(s.tableNamesQuery)
	if rError != nil {
		return
	}
	defer func() {
		if err := lRows.Close(); err != nil && rError == nil {
			rError = err
		}
	}()
	lTableName := ""
	for lRows.Next() {
		rError = lRows.Scan(&lTableName)
		if rError != nil {
			return
		}
		rTableNames = append(rTableNames, lTableName)
	}
	return
}

func (s *DBSchema) HasColumn(pTableName, pColumnName string) bool {
	lColumnTypes, lError := s.fetchColumnTypes(pTableName)
	if lError != nil {
		return false
	}
	return slices.ContainsFunc(lColumnTypes, func(ct *sql.ColumnType) bool {
		return strings.EqualFold(ct.Name(), pColumnName)
	})
}

func (s *DBSchema) fetchColumnTypes(pTableName string) (rColumnTypes []*sql.ColumnType, rError error) {
	lRows, rError := s.db.Query(s.columnTypesQueryFunc(pTableName))
	if rError != nil {
		return
	}
	defer func() {
		if err := lRows.Close(); err != nil && rError == nil {
			rError = err
		}
	}()
	rColumnTypes, rError = lRows.ColumnTypes()
	return
}
