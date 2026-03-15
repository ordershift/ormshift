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
	pk      *PrimaryKey
	fks     []ForeignKey
	ucs     []UniqueConstraint
}

func NewTable(name string) Table {
	return Table{
		name:    name,
		columns: []Column{},
		pk:      nil,
		fks:     []ForeignKey{},
		ucs:     []UniqueConstraint{},
	}
}

func (t *Table) Name() string {
	return t.name
}

func (t *Table) Columns() []Column {
	return t.columns
}

func (t *Table) PK() *PrimaryKey {
	return t.pk
}

func (t *Table) FKs() []ForeignKey {
	return t.fks
}

func (t *Table) UCs() []UniqueConstraint {
	return t.ucs
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

func (t *Table) PrimaryKey(columns ...string) error {
	if t.pk != nil {
		return fmt.Errorf("primary key already set for table %q", t.Name())
	}

	for _, col := range columns {
		found := false
		for _, tcol := range t.columns {
			if strings.EqualFold(col, tcol.Name()) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("primary key column %q does not exist in table %q", col, t.Name())
		}
	}
	pk := newPrimaryKey(t.Name(), columns)
	t.pk = &pk
	return nil
}

func (t *Table) AddForeignKey(fromColumns []string, toTable string, toColumns []string) error {
	for _, col := range fromColumns {
		found := false
		for _, tcol := range t.columns {
			if strings.EqualFold(col, tcol.Name()) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("foreign key column %q does not exist in table %q", col, t.Name())
		}
	}
	fk := newForeignKey(t.Name(), fromColumns, toTable, toColumns)
	t.fks = append(t.fks, fk)
	return nil
}

func (t *Table) AddUniqueConstraint(columns ...string) error {
	for _, col := range columns {
		found := false
		for _, tcol := range t.columns {
			if strings.EqualFold(col, tcol.Name()) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unique constraint column %q does not exist in table %q", col, t.Name())
		}
	}
	uc := newUniqueConstraint(t.Name(), columns)
	t.ucs = append(t.ucs, uc)
	return nil
}
