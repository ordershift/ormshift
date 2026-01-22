package schema_test

import (
	"fmt"
	"testing"

	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/schema"
)

func TestNewTableFailsWithInvalidName(t *testing.T) {
	lInvalidTableName := "123456-table"
	lTable, lError := schema.NewTable(lInvalidTableName)
	if !testutils.AssertNilResultAndNotNilError(t, lTable, lError, "schema.NewTable") {
		return
	}
	testutils.AssertErrorMessage(t, fmt.Sprintf("invalid table name: %q", lInvalidTableName), lError, "schema.NewTable")
}

func TestAddColumnFailsWithInvalidName(t *testing.T) {
	lProductAttributeTable := testutils.FakeProductAttributeTable(t)
	if lProductAttributeTable == nil {
		return
	}
	lInvalidColumnName := "123456-column"
	lError := lProductAttributeTable.AddColumn(schema.NewColumnParams{
		Name: lInvalidColumnName,
		Type: schema.Integer,
	})
	if !testutils.AssertNotNilError(t, lError, "Table.AddColumn") {
		return
	}
	testutils.AssertErrorMessage(t, fmt.Sprintf("invalid column name: %q", lInvalidColumnName), lError, "Table.AddColumn")
}

func TestAddColumnFailsWhenAlreadyExists(t *testing.T) {
	lProductAttributeTable := testutils.FakeProductAttributeTable(t)
	if lProductAttributeTable == nil {
		return
	}
	lError := lProductAttributeTable.AddColumn(schema.NewColumnParams{
		Name: "value",
		Type: schema.Integer,
	})
	if !testutils.AssertNotNilError(t, lError, "Table.AddColumn") {
		return
	}
	testutils.AssertErrorMessage(t, fmt.Sprintf("column %q already exists in table %q", "value", "product_attribute"), lError, "Table.AddColumn")
}
