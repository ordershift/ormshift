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
	}
}

func (t *Table) Name() string {
	return t.name
}

func (t *Table) Columns() []Column {
	return t.columns
}

func (t *Table) AddColumns(pParams ...NewColumnParams) error {
	for _, colParams := range pParams {
		column := NewColumn(colParams)
		columnAlreadyExists := slices.ContainsFunc(t.columns, func(c Column) bool {
			return strings.EqualFold(column.Name(), c.Name())
		})
		if columnAlreadyExists {
			return fmt.Errorf("column %q already exists in table %q", column.Name(), t.Name())
		}
		t.columns = append(t.columns, column)
	}
	return nil
}
