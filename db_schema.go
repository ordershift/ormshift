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
	db       *sql.DB
	driverDB DriverDB
}

func NewDBSchema(pDB *sql.DB, pDriverDB DriverDB) (*DBSchema, error) {
	if pDB == nil {
		return nil, errors.New("sql.DB cannot be nil")
	}
	if !pDriverDB.IsValid() {
		return nil, errors.New("driver db should be valid")
	}
	return &DBSchema{db: pDB, driverDB: pDriverDB}, nil
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
	lSQLQuery := ""
	switch s.driverDB {
	case DriverSQLServer:
		lSQLQuery = `
			SELECT
				t.name
			FROM
				sys.tables t
			LEFT JOIN
				sys.extended_properties ep
			ON	ep.major_id = t.[object_id]
			WHERE
				t.is_ms_shipped = 0 AND
				(ep.class_desc IS NULL OR (ep.class_desc <> 'OBJECT_OR_COLUMN' AND
					ep.[name] <> 'microsoft_database_tools_support'))
			ORDER BY
				t.name
		`
	case DriverSQLite:
		lSQLQuery = `
			SELECT
				name 
			FROM
				sqlite_master
			WHERE
				type = 'table'
			ORDER BY
				name
		`
	case DriverPostgresql:
		lSQLQuery = `
			SELECT
				table_name
			FROM
				information_schema.tables
			WHERE
				table_type = 'BASE TABLE' AND
				table_schema NOT IN ('pg_catalog', 'information_schema')
			ORDER BY
				table_name
		`
	}
	lRows, lError := s.db.Query(lSQLQuery)
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
	lColumnTypes, lErro := s.fetchColumnTypes(pTableName.String())
	if lErro != nil {
		lColumnTypes, lErro = s.fetchColumnTypes(strings.ToLower(pTableName.String()))
	}
	if lErro != nil {
		lColumnTypes, lErro = s.fetchColumnTypes(strings.ToUpper(pTableName.String()))
	}
	if lErro != nil {
		return nil, lErro
	}
	for _, lColumnType := range lColumnTypes {
		lColumnName := strings.ToUpper(lColumnType.Name())
		if lColumnName == strings.ToUpper(pColumnName.String()) {
			return lColumnType, nil
		}
	}
	return nil, fmt.Errorf("column %q not found in table %q", pColumnName.String(), pTableName.String())
}

func (s DBSchema) ExistsTableColumn(pTableName TableName, pColumnName ColumnName) bool {
	_, lError := s.CheckTableColumnType(pTableName, pColumnName)
	return lError == nil
}

func (s DBSchema) fetchColumnTypes(pTableName string) ([]*sql.ColumnType, error) {
	lRows, lError := s.db.Query(fmt.Sprintf("SELECT * FROM %s WHERE 1=0", pTableName))
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
