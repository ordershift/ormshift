package schema_test

import (
	"testing"

	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/schema"
)

func TestNewDBSchemaFailsWhenDBIsNil(t *testing.T) {
	lDBSchema, lError := schema.NewDBSchema(nil, "query")
	if !testutils.AssertNilResultAndNotNilError(t, lDBSchema, lError, "schema.NewDBSchema") {
		return
	}
	testutils.AssertErrorMessage(t, "sql.DB cannot be nil", lError, "schema.NewDBSchema")
}
