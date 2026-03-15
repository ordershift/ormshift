package schema_test

import (
	"testing"

	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/schema"
)

func TestColumn(t *testing.T) {
	column := schema.NewColumn(schema.NewColumnParams{Name: "id", Type: schema.Integer, NotNull: true, AutoIncrement: true})

	testutils.AssertEqualWithLabel(t, "id", column.Name(), "Column.Name")
	testutils.AssertEqualWithLabel(t, schema.Integer, column.Type(), "Column.Type")
	testutils.AssertEqualWithLabel(t, uint(0), column.Size(), "Column.Size")
	testutils.AssertEqualWithLabel(t, true, column.NotNull(), "Column.IsNotNull")
	testutils.AssertEqualWithLabel(t, true, column.AutoIncrement(), "Column.IsAutoIncrement")
	testutils.AssertEqualWithLabel(t, "", column.Default(), "Column.Default")
}

func TestColumnWithDefault(t *testing.T) {
	column := schema.NewColumn(schema.NewColumnParams{Name: "count", Type: schema.Integer, Default: "0"})
	testutils.AssertEqualWithLabel(t, "0", column.Default(), "Column.Default")
}

func TestColumnWithCheck(t *testing.T) {
	column := schema.NewColumn(schema.NewColumnParams{Name: "price", Type: schema.Monetary, Check: "price >= 0"})
	testutils.AssertEqualWithLabel(t, "price >= 0", column.Check(), "Column.Check")
}
