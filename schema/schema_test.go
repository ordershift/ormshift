package schema_test

import (
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/sqlite"
	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/schema"
)

func TestNewDBSchema(t *testing.T) {
	lDB, lError := ormshift.OpenDatabase(sqlite.SQLiteDriver{}, ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNilError(t, lError, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = lDB.Close() }()

	lDBSchema, lError := schema.NewDBSchema(lDB.DB(), "query")
	if !testutils.AssertNotNilResultAndNilError(t, lDBSchema, lError, "schema.NewDBSchema") {
		return
	}
}

func TestNewDBSchemaFailsWhenDBIsNil(t *testing.T) {
	lDBSchema, lError := schema.NewDBSchema(nil, "query")
	if !testutils.AssertNilResultAndNotNilError(t, lDBSchema, lError, "schema.NewDBSchema") {
		return
	}
	testutils.AssertErrorMessage(t, "sql.DB cannot be nil", lError, "schema.NewDBSchema")
}
