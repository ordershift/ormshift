package schema

import "fmt"

type PrimaryKey struct {
	name    string
	columns []string
}

func newPrimaryKey(table string, columns []string) PrimaryKey {
	return PrimaryKey{
		name:    fmt.Sprintf("PK_%s", table),
		columns: columns,
	}
}

func (pk *PrimaryKey) Name() string {
	return pk.name
}

func (pk *PrimaryKey) Columns() []string {
	return pk.columns
}
