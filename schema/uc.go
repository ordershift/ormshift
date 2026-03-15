package schema

import (
	"fmt"
	"strings"
)

type UniqueConstraint struct {
	name    string
	columns []string
}

func newUniqueConstraint(table string, columns []string) UniqueConstraint {
	return UniqueConstraint{
		name:    fmt.Sprintf("UC_%s_%s", table, strings.Join(columns, "_")),
		columns: columns,
	}
}
