package schema

import (
	"database/sql"
	"errors"
	"fmt"
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

func (s DBSchema) ExistsTable(pTableName TableName) bool {
	lTables, lError := s.fetchTableNames()
	if lError != nil {
		return false
	}
	for _, lTable := range lTables {
		lUpperTableName := strings.ToUpper(lTable)
		if lUpperTableName == strings.ToUpper(pTableName.String()) {
			return true
		}
	}
	return false
}

func (s DBSchema) fetchTableNames() ([]string, error) {
	lRows, lError := s.db.Query(s.tableNamesQuery)
	if lError != nil {
		return nil, lError
	}
	defer lRows.Close()

	var lTableNames []string
	lTableName := ""
	for lRows.Next() {
		lError = lRows.Scan(&lTableName)
		if lError != nil {
			return nil, lError
		}
		lTableNames = append(lTableNames, lTableName)
	}
	return lTableNames, nil
}

func (s DBSchema) CheckTableColumnType(pTableName TableName, pColumnName ColumnName) (*sql.ColumnType, error) {
	// TODO: Maybe the table match should be case sensitive depending on the DBMS
	lColumnTypes, lError := s.fetchColumnTypes(pTableName)
	// if lErro != nil {
	// 	lColumnTypes, lErro = s.fetchColumnTypes(strings.ToLower(pTableName.String()))
	// }
	// if lErro != nil {
	// 	lColumnTypes, lErro = s.fetchColumnTypes(strings.ToUpper(pTableName.String()))
	// }
	if lError != nil {
		return nil, lError
	}
	for _, lColumnType := range lColumnTypes {
		// TODO: Maybe the table match should be case sensitive depending on the DBMS
		// lColumnName := strings.ToUpper(lColumnType.Name())
		// if lColumnName == strings.ToUpper(pColumnName.String()) {
		if lColumnType.Name() == pColumnName.String() {
			return lColumnType, nil
		}
	}
	return nil, fmt.Errorf("column %q not found in table %q", pColumnName.String(), pTableName.String())
}

func (s DBSchema) ExistsTableColumn(pTableName TableName, pColumnName ColumnName) bool {
	_, lError := s.CheckTableColumnType(pTableName, pColumnName)
	return lError == nil
}

func (s DBSchema) fetchColumnTypes(pTableName TableName) ([]*sql.ColumnType, error) {
	lTableName := pTableName.String()
	if !regexValidTableName.MatchString(lTableName) {
		return nil, fmt.Errorf("invalid table name: %q", lTableName)
	}

	lRows, lError := s.db.Query(fmt.Sprintf("SELECT * FROM %s WHERE 1=0", lTableName))
	if lError != nil {
		return nil, lError
	}
	defer lRows.Close()
	return lRows.ColumnTypes()
}
