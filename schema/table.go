package schema

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ordershift/ormshift/errs"
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
			return failedToAddColumnInTable(*t, column, errs.AlreadyExists("column"))
		}
		t.columns = append(t.columns, column)
	}
	return nil
}

func failedToAddColumnInTable(table Table, column Column, err error) error {
	msg := fmt.Sprintf("add column %q in table %q", column.Name(), table.Name())
	return errs.FailedTo(msg, err)
}
