package schema

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

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
	lColumnAlreadyExists := slices.ContainsFunc(t.columns, func(c Column) bool {
		return strings.EqualFold(lColumn.Name().String(), c.Name().String())
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
