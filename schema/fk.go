package schema

import "fmt"

type ForeignKey struct {
	name        string
	fromColumns []string
	toTable     string
	toColumns   []string
}

func newForeignKey(fromTable string, fromColumns []string, toTable string, toColumns []string) ForeignKey {
	return ForeignKey{
		name:        fmt.Sprintf("FK_%s_%s", fromTable, toTable),
		fromColumns: fromColumns,
		toTable:     toTable,
		toColumns:   toColumns,
	}
}

func (fk *ForeignKey) Name() string {
	return fk.name
}

func (fk *ForeignKey) FromColumns() []string {
	return fk.fromColumns
}

func (fk *ForeignKey) ToTable() string {
	return fk.toTable
}

func (fk *ForeignKey) ToColumns() []string {
	return fk.toColumns
}
