package sqlite_test

import (
	"database/sql"
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/schema"
	"github.com/ordershift/ormshift/sqlite"
)

func Test_DBSchema_TableAndColumnExists_ShouldReturnTrue(t *testing.T) {
	lDB, lError := sql.Open(sqlite.DriverName(), sqlite.ConnectionString(ormshift.ConnectionParams{InMemory: true}))
	if !testutils.AssertNotNilResultAndNilError(t, lDB, lError, "sql.Open") {
		return
	}
	defer lDB.Close()
	lProductAttributeTable := testutils.FakeProductAttributeTable(t)
	if lProductAttributeTable == nil {
		return
	}
	_, lError = lDB.Exec(sqlite.SQLBuilder().CreateTable(*lProductAttributeTable))
	if !testutils.AssertNilError(t, lError, "DB.Exec") {
		return
	}
	lDBSchema, lError := sqlite.DBSchema(lDB)
	if !testutils.AssertNotNilResultAndNilError(t, lDBSchema, lError, "ormshift.NewDBSchema") {
		return
	}
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

func Test_DBSchema_TableExists_ShouldReturnFalse_WhenDBIsClosed(t *testing.T) {
	lDB, lError := sql.Open(sqlite.DriverName(), sqlite.ConnectionString(ormshift.ConnectionParams{InMemory: true}))
	if !testutils.AssertNotNilResultAndNilError(t, lDB, lError, "sql.Open") {
		return
	}
	lProductAttributeTable := testutils.FakeProductAttributeTable(t)
	if lProductAttributeTable == nil {
		lDB.Close()
		return
	}
	_, lError = lDB.Exec(sqlite.SQLBuilder().CreateTable(*lProductAttributeTable))
	if !testutils.AssertNilError(t, lError, "DB.Exec") {
		lDB.Close()
		return
	}
	lDB.Close()
	lDBSchema, lError := sqlite.DBSchema(lDB)
	if !testutils.AssertNotNilResultAndNilError(t, lDBSchema, lError, "ormshift.NewDBSchema") {
		return
	}
	testutils.AssertEqualWithLabel(t, false, lDBSchema.ExistsTable(lProductAttributeTable.Name()), "DBSchema.ExistsTable")
}
