package sqlite_test

import (
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/sqlite"
	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/schema"
)

func Test_DBSchema_TableAndColumnExists_ShouldReturnTrue(t *testing.T) {
	lDatabase, lError := ormshift.OpenDatabase(sqlite.SQLiteDriver{}, ormshift.ConnectionParams{InMemory: true})
	if lError != nil {
		t.Errorf("ormshift.OpenDatabase failed: %v", lError)
		return
	}
	defer func() { _ = lDatabase.Close() }()

	lProductAttributeTable := testutils.FakeProductAttributeTable(t)
	if lProductAttributeTable == nil {
		return
	}

	_, lError = lDatabase.DB().Exec(sqlite.SQLiteDriver{}.SQLBuilder().CreateTable(*lProductAttributeTable))
	if !testutils.AssertNilError(t, lError, "DB.Exec") {
		return
	}

	lDBSchema := lDatabase.DBSchema()
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

func Test_DBSchema_TableExists_ShouldReturnFalse_WhenDatabaseIsInvalid(t *testing.T) {
	lDatabase, lError := ormshift.OpenDatabase(sqlite.SQLiteDriver{}, ormshift.ConnectionParams{InMemory: true})
	if lError != nil {
		t.Errorf("ormshift.OpenDatabase failed: %v", lError)
		return
	}
	defer func() { _ = lDatabase.Close() }()

	lProductAttributeTable := testutils.FakeProductAttributeTable(t)
	if lProductAttributeTable == nil {
		_ = lDatabase.Close()
		return
	}

	_, lError = lDatabase.DB().Exec(sqlite.SQLiteDriver{}.SQLBuilder().CreateTable(*lProductAttributeTable))
	if !testutils.AssertNilError(t, lError, "DB.Exec") {
		_ = lDatabase.Close()
		return
	}
	_ = lDatabase.Close()
	lDBSchema := lDatabase.DBSchema()
	testutils.AssertEqualWithLabel(t, false, lDBSchema.ExistsTable(lProductAttributeTable.Name()), "DBSchema.ExistsTable")
}
