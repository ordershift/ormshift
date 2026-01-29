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
	testutils.AssertErrorMessage(t, fmt.Sprintf("column %q already exists in table %q", "value", "product_attribute"), err, "Table.AddColumns")
}
