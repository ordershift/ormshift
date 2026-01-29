package schema_test

import (
	"testing"

	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/schema"
)

func TestColumn(t *testing.T) {
	column := schema.NewColumn(schema.NewColumnParams{Name: "id", Type: schema.Integer, NotNull: true, PrimaryKey: true, AutoIncrement: true})

	testutils.AssertEqualWithLabel(t, "id", column.Name(), "Column.Name")
	testutils.AssertEqualWithLabel(t, schema.Integer, column.Type(), "Column.Type")
	testutils.AssertEqualWithLabel(t, uint(0), column.Size(), "Column.Size")
	testutils.AssertEqualWithLabel(t, true, column.PrimaryKey(), "Column.IsPrimaryKey")
	testutils.AssertEqualWithLabel(t, true, column.NotNull(), "Column.IsNotNull")
	testutils.AssertEqualWithLabel(t, true, column.AutoIncrement(), "Column.IsAutoIncrement")
}
