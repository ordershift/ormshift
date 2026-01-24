package schema_test

import (
	"testing"

	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/schema"
)

func TestColumn(t *testing.T) {
	lColumn, lError := schema.NewColumn(schema.NewColumnParams{Name: "id", Type: schema.Integer, NotNull: true, PrimaryKey: true, AutoIncrement: true})
	if !testutils.AssertNilError(t, lError, "schema.NewColumn") {
		return
	}
	testutils.AssertEqualWithLabel(t, "id", lColumn.Name().String(), "Column.Name")
	testutils.AssertEqualWithLabel(t, schema.Integer, lColumn.Type(), "Column.Type")
	testutils.AssertEqualWithLabel(t, uint(0), lColumn.Size(), "Column.Size")
	testutils.AssertEqualWithLabel(t, true, lColumn.PrimaryKey(), "Column.IsPrimaryKey")
	testutils.AssertEqualWithLabel(t, true, lColumn.NotNull(), "Column.IsNotNull")
	testutils.AssertEqualWithLabel(t, true, lColumn.AutoIncrement(), "Column.IsAutoIncrement")
}
