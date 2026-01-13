package ormshift_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/ordershift/ormshift"
)

func Test_DBSchema_NewDBSchema_ShouldFail_WhenDBIsNil(t *testing.T) {
	lDBSchema, lError := ormshift.NewDBSchema(nil, ormshift.DriverSQLite)
	if !assertNilResultAndNotNilError(t, lDBSchema, lError, "ormshift.NewDBSchema") {
		return
	}
	assertErrorMessage(t, "sql.DB cannot be nil", lError, "ormshift.NewDBSchema")
}

func Test_DBSchema_TableAndColumnExists_ShouldReturnTrue(t *testing.T) {
	lDB, lError := sql.Open(ormshift.DriverSQLite.Name(), ormshift.DriverSQLite.ConnectionString(ormshift.ConnectionParams{InMemory: true}))
	if !assertNotNilResultAndNilError(t, lDB, lError, "sql.Open") {
		return
	}
	defer lDB.Close()
	lProductAttributeTable := fakeProductAttributeTable(t)
	if lProductAttributeTable == nil {
		return
	}
	_, lError = lDB.Exec(ormshift.DriverSQLite.SQLBuilder().CreateTable(*lProductAttributeTable))
	if !assertNilError(t, lError, "DB.Exec") {
		return
	}
	lDBSchema, lError := ormshift.NewDBSchema(lDB, ormshift.DriverSQLite)
	if !assertNotNilResultAndNilError(t, lDBSchema, lError, "ormshift.NewDBSchema") {
		return
	}
	assertEqualWithLabel(t, true, lDBSchema.ExistsTable(lProductAttributeTable.Name()), "DBSchema.ExistsTable")
	for _, lColumn := range lProductAttributeTable.Columns() {
		assertEqualWithLabel(t, true, lDBSchema.ExistsTableColumn(lProductAttributeTable.Name(), lColumn.Name()), "DBSchema.ExistsTableColumn")
	}
	lAnyTableName, lError := ormshift.NewTableName("any_table")
	if !assertNotNilResultAndNilError(t, lAnyTableName, lError, "ormshift.NewTableName") {
		return
	}
	assertEqualWithLabel(t, false, lDBSchema.ExistsTable(*lAnyTableName), "DBSchema.ExistsTable")
	lAnyColumnName, lError := ormshift.NewColumnName("any_col")
	if !assertNotNilResultAndNilError(t, lAnyColumnName, lError, "ormshift.NewTableName") {
		return
	}
	assertEqualWithLabel(t, false, lDBSchema.ExistsTableColumn(lProductAttributeTable.Name(), *lAnyColumnName), "DBSchema.ExistsTableColumn")
	assertEqualWithLabel(t, false, lDBSchema.ExistsTableColumn(*lAnyTableName, *lAnyColumnName), "DBSchema.ExistsTableColumn")
}

func Test_DBSchema_TableExists_ShouldReturnFalse_WhenDBIsClosed(t *testing.T) {
	lDB, lError := sql.Open(ormshift.DriverSQLite.Name(), ormshift.DriverSQLite.ConnectionString(ormshift.ConnectionParams{InMemory: true}))
	if !assertNotNilResultAndNilError(t, lDB, lError, "sql.Open") {
		return
	}
	lProductAttributeTable := fakeProductAttributeTable(t)
	if lProductAttributeTable == nil {
		lDB.Close()
		return
	}
	_, lError = lDB.Exec(ormshift.DriverSQLite.SQLBuilder().CreateTable(*lProductAttributeTable))
	if !assertNilError(t, lError, "DB.Exec") {
		lDB.Close()
		return
	}
	lDB.Close()
	lDBSchema, lError := ormshift.NewDBSchema(lDB, ormshift.DriverSQLite)
	if !assertNotNilResultAndNilError(t, lDBSchema, lError, "ormshift.NewDBSchema") {
		return
	}
	assertEqualWithLabel(t, false, lDBSchema.ExistsTable(lProductAttributeTable.Name()), "DBSchema.ExistsTable")
}

func Test_DBSchema_NewTable_ShouldFail_WhenHasInvalidName(t *testing.T) {
	lInvalidTableName := "123456-table"
	lTable, lError := ormshift.NewTable(lInvalidTableName)
	if !assertNilResultAndNotNilError(t, lTable, lError, "ormshift.NewTable") {
		return
	}
	assertErrorMessage(t, fmt.Sprintf("invalid table name: %q", lInvalidTableName), lError, "ormshift.NewTable")
}

func Test_DBSchema_Table_AddColumn_ShouldFail_WhenHasInvalidName(t *testing.T) {
	lProductAttributeTable := fakeProductAttributeTable(t)
	if lProductAttributeTable == nil {
		return
	}
	lInvalidColumnName := "123456-column"
	lError := lProductAttributeTable.AddColumn(ormshift.NewColumnParams{
		Name: lInvalidColumnName,
		Type: ormshift.Integer,
	})
	if !assertNotNilError(t, lError, "Table.AddColumn") {
		return
	}
	assertErrorMessage(t, fmt.Sprintf("invalid column name: %q", lInvalidColumnName), lError, "Table.AddColumn")
}

func Test_DBSchema_Table_AddColumn_ShouldFail_WhenAlreadyExists(t *testing.T) {
	lProductAttributeTable := fakeProductAttributeTable(t)
	if lProductAttributeTable == nil {
		return
	}
	lError := lProductAttributeTable.AddColumn(ormshift.NewColumnParams{
		Name: "value",
		Type: ormshift.Integer,
	})
	if !assertNotNilError(t, lError, "Table.AddColumn") {
		return
	}
	assertErrorMessage(t, fmt.Sprintf("column %q already exists in table %q", "value", "product_attribute"), lError, "Table.AddColumn")
}
