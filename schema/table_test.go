package schema_test

import (
	"fmt"
	"testing"

	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/schema"
)

func TestAddColumnFailsWhenAlreadyExists(t *testing.T) {
	lProductAttributeTable := testutils.FakeProductAttributeTable(t)
	lError := lProductAttributeTable.AddColumns(schema.NewColumnParams{
		Name: "value",
		Type: schema.Integer,
	})
	if !testutils.AssertNotNilError(t, lError, "Table.AddColumn") {
		return
	}
	testutils.AssertErrorMessage(t, fmt.Sprintf("column %q already exists in table %q", "value", "product_attribute"), lError, "Table.AddColumn")
}
