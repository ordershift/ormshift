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

func TestExistsTableColumn(t *testing.T) {
	lDB, lError := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if lError != nil {
		t.Errorf("ormshift.OpenDatabase failed: %v", lError)
		return
	}
	defer func() { _ = lDB.Close() }()

	lProductAttributeTable := testutils.FakeProductAttributeTable(t)
	if lProductAttributeTable == nil {
		return
	}

	_, lError = lDB.SQLExecutor().Exec(sqlite.Driver().SQLBuilder().CreateTable(*lProductAttributeTable))
	if !testutils.AssertNilError(t, lError, "DB.Exec") {
		return
	}

	lDBSchema := lDB.DBSchema()
	testutils.AssertEqualWithLabel(t, true, lDBSchema.ExistsTable(lProductAttributeTable.Name()), "DBSchema.ExistsTable")
	for _, lColumn := range lProductAttributeTable.Columns() {
		testutils.AssertEqualWithLabel(t, true, lDBSchema.ExistsTableColumn(lProductAttributeTable.Name(), lColumn.Name()), "DBSchema.ExistsTableColumn")
	}
	lAnyTableName, lError := schema.NewTableName("any_table")
	if !testutils.AssertNotNilResultAndNilError(t, lAnyTableName, lError, "ormshift.NewTableName") {
		return
	}
	testutils.AssertEqualWithLabel(t, false, lDBSchema.ExistsTable(*lAnyTableName), "DBSchema.ExistsTable")
	lAnyColumnName, lError := schema.NewColumnName("any_col")
	if !testutils.AssertNotNilResultAndNilError(t, lAnyColumnName, lError, "ormshift.NewTableName") {
		return
	}
	testutils.AssertEqualWithLabel(t, false, lDBSchema.ExistsTableColumn(lProductAttributeTable.Name(), *lAnyColumnName), "DBSchema.ExistsTableColumn")
	testutils.AssertEqualWithLabel(t, false, lDBSchema.ExistsTableColumn(*lAnyTableName, *lAnyColumnName), "DBSchema.ExistsTableColumn")
}

func TestExistsTableReturnsFalseWhenDatabaseIsInvalid(t *testing.T) {
	lDB, lError := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if lError != nil {
		t.Errorf("ormshift.OpenDatabase failed: %v", lError)
		return
	}
	defer func() { _ = lDB.Close() }()

	lProductAttributeTable := testutils.FakeProductAttributeTable(t)
	if lProductAttributeTable == nil {
		_ = lDB.Close()
		return
	}

	_, lError = lDB.SQLExecutor().Exec(sqlite.Driver().SQLBuilder().CreateTable(*lProductAttributeTable))
	if !testutils.AssertNilError(t, lError, "DB.Exec") {
		_ = lDB.Close()
		return
	}
	_ = lDB.Close()
	lDBSchema := lDB.DBSchema()
	testutils.AssertEqualWithLabel(t, false, lDBSchema.ExistsTable(lProductAttributeTable.Name()), "DBSchema.ExistsTable")
}
