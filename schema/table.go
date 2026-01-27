package schema

import (
	"fmt"
	"slices"
	"strings"
)

type Table struct {
	name    string
	columns []Column
}

func NewTable(pName string) Table {
	return Table{
		name:    pName,
		columns: []Column{},

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

func (t Table) Name() string {
	return t.name
}

func (t Table) Columns() []Column {
	return t.columns
}
