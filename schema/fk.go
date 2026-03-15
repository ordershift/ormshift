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
