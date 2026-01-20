package ormshift

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
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

var regexValidTableName = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*\.)*[A-Za-z_][A-Za-z0-9_]*$`)

type TableName struct {
	tableName string
}

func NewTableName(pName string) (*TableName, error) {
	if !regexValidTableName.MatchString(pName) {
		return nil, fmt.Errorf("invalid table name: %q", pName)
	}
	return &TableName{pName}, nil
}

func (tn TableName) String() string {
	return tn.tableName
}

var regexValidColumnName = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)

type ColumnName struct {
	columnName string
}

func NewColumnName(pName string) (*ColumnName, error) {
	if !regexValidColumnName.MatchString(pName) {
		return nil, fmt.Errorf("invalid column name: %q", pName)
	}
	return &ColumnName{pName}, nil
}

func (tn ColumnName) String() string {
	return tn.columnName
}

type ColumnType int

const (
	Integer ColumnType = iota
	Varchar
	Monetary
	DateTime
	Decimal
	Boolean
	Binary
)

type Table struct {
	name    TableName
	columns []Column
}

func NewTable(pName string) (*Table, error) {
	lTableName, lError := NewTableName(pName)
	if lError != nil {
		return nil, lError
	}
	return &Table{
		name:    *lTableName,
		columns: []Column{},
	}, nil
}

func (t *Table) AddColumn(pParams NewColumnParams) error {
	lColumn, lError := NewColumn(pParams)
	if lError != nil {
		return lError
	}
	lLowerColumnName := strings.ToLower(lColumn.Name().String())
	lColumnAlreadyExists := slices.ContainsFunc(t.columns, func(c Column) bool {
		return lLowerColumnName == strings.ToLower(c.Name().String())
	})
	if lColumnAlreadyExists {
		return fmt.Errorf("column %q already exists in table %q", lColumn.Name().String(), t.name)
	}
	t.columns = append(t.columns, *lColumn)
	return nil
}

func (t Table) Name() TableName {
	return t.name
}

func (t Table) Columns() []Column {
	return t.columns
}

type NewColumnParams struct {
	Name          string
	Type          ColumnType
	Size          uint
	PrimaryKey    bool
	NotNull       bool
	Autoincrement bool
}

type Column struct {
	name       ColumnName
	columnType ColumnType
	size       uint
	pk         bool
	notNull    bool
	autoInc    bool
}

func NewColumn(pParams NewColumnParams) (*Column, error) {
	lColumnName, lError := NewColumnName(pParams.Name)
	if lError != nil {
		return nil, lError
	}
	return &Column{
		name:       *lColumnName,
		columnType: pParams.Type,
		size:       pParams.Size,
		pk:         pParams.PrimaryKey,
		notNull:    pParams.NotNull,
		autoInc:    pParams.Autoincrement,
	}, nil
}

func (c Column) Name() ColumnName {
	return c.name
}

func (c Column) Type() ColumnType {
	return c.columnType
}

func (c Column) Size() uint {
	return c.size
}

func (c Column) PrimaryKey() bool {
	return c.pk
}

func (c Column) NotNull() bool {
	return c.notNull
}

func (c Column) Autoincrement() bool {
	return c.autoInc
}
