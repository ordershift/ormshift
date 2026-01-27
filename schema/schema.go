package schema

import (
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
)

type DBSchema struct {
	db              *sql.DB
	tableNamesQuery string
}

func NewDBSchema(pDB *sql.DB, pTableNamesQuery string) (*DBSchema, error) {
	if pDB == nil {
		return nil, errors.New("sql.DB cannot be nil")
	}
	return &DBSchema{db: pDB, tableNamesQuery: pTableNamesQuery}, nil
}

func (s DBSchema) HasTable(pTableName TableName) bool {
	lTables, lError := s.fetchTableNames()
	if lError != nil {
		return false
	}
	return slices.ContainsFunc(lTables, func(t string) bool {
		return strings.EqualFold(t, pTableName.String())
	})
}

func (s DBSchema) fetchTableNames() (rTableNames []string, rError error) {
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

func (s DBSchema) HasColumn(pTableName TableName, pColumnName ColumnName) bool {
	lColumnTypes, lError := s.fetchColumnTypes(pTableName)
	if lError != nil {
		return false
	}
	return slices.ContainsFunc(lColumnTypes, func(ct *sql.ColumnType) bool {
		return strings.EqualFold(ct.Name(), pColumnName.String())
	})
}

func (s DBSchema) fetchColumnTypes(pTableName TableName) (rColumnTypes []*sql.ColumnType, rError error) {
	lRows, rError := s.db.Query(fmt.Sprintf("SELECT * FROM %s WHERE 1=0", pTableName.String())) // NOSONAR go:S2077 - Dynamic SQL is controlled and sanitized internally
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
