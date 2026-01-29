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

func NewTable(name string) Table {
	return Table{
		name:    name,
		columns: []Column{},
	}
}

func (t *Table) Name() string {
	return t.name
}

func (t *Table) Columns() []Column {
	return t.columns
}

func (t *Table) AddColumns(params ...NewColumnParams) error {
	for _, colParams := range params {
		column := NewColumn(colParams)
		exists := slices.ContainsFunc(t.columns, func(c Column) bool {
			return strings.EqualFold(column.Name(), c.Name())
		})
		if exists {
			return fmt.Errorf("column %q already exists in table %q", column.Name(), t.Name())
		}
		t.columns = append(t.columns, column)
	}
	return nil
}
