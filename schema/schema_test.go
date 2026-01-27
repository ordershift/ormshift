package schema_test

import (
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/sqlite"
	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/schema"
)

func TestNewDBSchema(t *testing.T) {
	lDB, lError := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
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

func TestHasColumn(t *testing.T) {
	lDB, lError := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if lError != nil {
		t.Errorf("ormshift.OpenDatabase failed: %v", lError)
		return
	}
	defer func() { _ = lDB.Close() }()

	lProductAttributeTable := testutils.FakeProductAttributeTable(t)

	_, lError = lDB.SQLExecutor().Exec(sqlite.Driver().SQLBuilder().CreateTable(lProductAttributeTable))
	if !testutils.AssertNilError(t, lError, "DB.Exec") {
		return
	}

	lDBSchema := lDB.DBSchema()
	testutils.AssertEqualWithLabel(t, true, lDBSchema.HasTable(lProductAttributeTable.Name()), "DBSchema.HasTable")
	for _, lColumn := range lProductAttributeTable.Columns() {
		testutils.AssertEqualWithLabel(t, true, lDBSchema.HasColumn(lProductAttributeTable.Name(), lColumn.Name()), "DBSchema.HasColumn")
	}
	lAnyTableName := "any_table"
	lAnyColumnName := "any_col"
	testutils.AssertEqualWithLabel(t, false, lDBSchema.HasColumn(lProductAttributeTable.Name(), lAnyColumnName), "DBSchema.HasColumn")
	testutils.AssertEqualWithLabel(t, false, lDBSchema.HasColumn(lAnyTableName, lAnyColumnName), "DBSchema.HasColumn")
}

func TestHasTableReturnsFalseWhenDatabaseIsInvalid(t *testing.T) {
	lDB, lError := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if lError != nil {
		t.Errorf("ormshift.OpenDatabase failed: %v", lError)
		return
	}
	defer func() { _ = lDB.Close() }()

	lProductAttributeTable := testutils.FakeProductAttributeTable(t)

	_, lError = lDB.SQLExecutor().Exec(sqlite.Driver().SQLBuilder().CreateTable(lProductAttributeTable))
	if !testutils.AssertNilError(t, lError, "DB.Exec") {
		_ = lDB.Close()
		return
	}
	_ = lDB.Close()
	lDBSchema := lDB.DBSchema()
	testutils.AssertEqualWithLabel(t, false, lDBSchema.HasTable(lProductAttributeTable.Name()), "DBSchema.HasTable")
}
