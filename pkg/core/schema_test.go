package core_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/pkg/core"
	"github.com/ordershift/ormshift/pkg/dialects/sqlite"
)

func Test_DBSchema_NewDBSchema_ShouldFail_WhenDBIsNil(t *testing.T) {
	lDBSchema, lError := core.NewDBSchema(nil, "query")
	if !testutils.AssertNilResultAndNotNilError(t, lDBSchema, lError, "ormshift.NewDBSchema") {
		return
	}
	testutils.AssertErrorMessage(t, "sql.DB cannot be nil", lError, "ormshift.NewDBSchema")
}

func Test_DBSchema_TableAndColumnExists_ShouldReturnTrue(t *testing.T) {
	lDB, lError := sql.Open(sqlite.DriverName(), sqlite.ConnectionString(core.ConnectionParams{InMemory: true}))
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
	lAnyTableName, lError := core.NewTableName("any_table")
	if !testutils.AssertNotNilResultAndNilError(t, lAnyTableName, lError, "ormshift.NewTableName") {
		return
	}
	testutils.AssertEqualWithLabel(t, false, lDBSchema.ExistsTable(*lAnyTableName), "DBSchema.ExistsTable")
	lAnyColumnName, lError := core.NewColumnName("any_col")
	if !testutils.AssertNotNilResultAndNilError(t, lAnyColumnName, lError, "ormshift.NewTableName") {
		return
	}
	testutils.AssertEqualWithLabel(t, false, lDBSchema.ExistsTableColumn(lProductAttributeTable.Name(), *lAnyColumnName), "DBSchema.ExistsTableColumn")
	testutils.AssertEqualWithLabel(t, false, lDBSchema.ExistsTableColumn(*lAnyTableName, *lAnyColumnName), "DBSchema.ExistsTableColumn")
}

func Test_DBSchema_TableExists_ShouldReturnFalse_WhenDBIsClosed(t *testing.T) {
	lDB, lError := sql.Open(sqlite.DriverName(), sqlite.ConnectionString(core.ConnectionParams{InMemory: true}))
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

func Test_DBSchema_NewTable_ShouldFail_WhenHasInvalidName(t *testing.T) {
	lInvalidTableName := "123456-table"
	lTable, lError := core.NewTable(lInvalidTableName)
	if !testutils.AssertNilResultAndNotNilError(t, lTable, lError, "ormshift.NewTable") {
		return
	}
	testutils.AssertErrorMessage(t, fmt.Sprintf("invalid table name: %q", lInvalidTableName), lError, "ormshift.NewTable")
}

func Test_DBSchema_Table_AddColumn_ShouldFail_WhenHasInvalidName(t *testing.T) {
	lProductAttributeTable := testutils.FakeProductAttributeTable(t)
	if lProductAttributeTable == nil {
		return
	}
	lInvalidColumnName := "123456-column"
	lError := lProductAttributeTable.AddColumn(core.NewColumnParams{
		Name: lInvalidColumnName,
		Type: core.Integer,
	})
	if !testutils.AssertNotNilError(t, lError, "Table.AddColumn") {
		return
	}
	testutils.AssertErrorMessage(t, fmt.Sprintf("invalid column name: %q", lInvalidColumnName), lError, "Table.AddColumn")
}

func Test_DBSchema_Table_AddColumn_ShouldFail_WhenAlreadyExists(t *testing.T) {
	lProductAttributeTable := testutils.FakeProductAttributeTable(t)
	if lProductAttributeTable == nil {
		return
	}
	lError := lProductAttributeTable.AddColumn(core.NewColumnParams{
		Name: "value",
		Type: core.Integer,
	})
	if !testutils.AssertNotNilError(t, lError, "Table.AddColumn") {
		return
	}
	testutils.AssertErrorMessage(t, fmt.Sprintf("column %q already exists in table %q", "value", "product_attribute"), lError, "Table.AddColumn")
}
