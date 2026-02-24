package schema_test

import (
	"fmt"
	"testing"

	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/schema"
)

func TestAddColumnFailsWhenAlreadyExists(t *testing.T) {
	productAttributeTable := testutils.FakeProductAttributeTable(t)
	err := productAttributeTable.AddColumns(schema.NewColumnParams{
		Name: "value",
		Type: schema.Integer,
	})
	if !testutils.AssertNotNilError(t, err, "Table.AddColumns") {
		return
	}
	testutils.AssertErrorMessage(t, fmt.Sprintf("failed to add column %q in table %q: column already exists", "value", "product_attribute"), err, "Table.AddColumns")
}
