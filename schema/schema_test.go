package schema_test

import (
	"fmt"
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/sqlite"
	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/schema"
)

func testColumnTypesQueryFunc(pTableName string) string {
	return fmt.Sprintf("SELECT * FROM %s WHERE 1=0", pTableName)
}

func TestNewDBSchema(t *testing.T) {
	db, err := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNotNilResultAndNilError(t, db, err, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = db.Close() }()

	dbSchema, err := schema.NewDBSchema(db.DB(), "query", testColumnTypesQueryFunc)
	if !testutils.AssertNotNilResultAndNilError(t, dbSchema, err, "schema.NewDBSchema") {
		return
	}
}

func TestNewDBSchemaFailsWhenDBIsNil(t *testing.T) {
	dbSchema, err := schema.NewDBSchema(nil, "query", testColumnTypesQueryFunc)
	if !testutils.AssertNilResultAndNotNilError(t, dbSchema, err, "schema.NewDBSchema") {
		return
	}
	testutils.AssertErrorMessage(t, "sql.DB cannot be nil", err, "schema.NewDBSchema")
}

func TestHasColumn(t *testing.T) {
	db, err := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNotNilResultAndNilError(t, db, err, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = db.Close() }()

	productAttributeTable := testutils.FakeProductAttributeTable(t)

	_, err = db.SQLExecutor().Exec(sqlite.Driver().SQLBuilder().CreateTable(productAttributeTable))
	if !testutils.AssertNilError(t, err, "DB.Exec") {
		return
	}

	dbSchema := db.DBSchema()
	testutils.AssertEqualWithLabel(t, true, dbSchema.HasTable(productAttributeTable.Name()), "DBSchema.HasTable")
	for _, column := range productAttributeTable.Columns() {
		testutils.AssertEqualWithLabel(t, true, dbSchema.HasColumn(productAttributeTable.Name(), column.Name()), "DBSchema.HasColumn")
	}
	anyTableName := "any_table"
	anyColumnName := "any_col"
	testutils.AssertEqualWithLabel(t, false, dbSchema.HasColumn(productAttributeTable.Name(), anyColumnName), "DBSchema.HasColumn")
	testutils.AssertEqualWithLabel(t, false, dbSchema.HasColumn(anyTableName, anyColumnName), "DBSchema.HasColumn")
}

func TestHasTableReturnsFalseWhenDatabaseIsInvalid(t *testing.T) {
	db, err := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNotNilResultAndNilError(t, db, err, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = db.Close() }()

	productAttributeTable := testutils.FakeProductAttributeTable(t)

	_, err = db.SQLExecutor().Exec(sqlite.Driver().SQLBuilder().CreateTable(productAttributeTable))
	if !testutils.AssertNilError(t, err, "DB.Exec") {
		return
	}

	err = db.Close()
	if !testutils.AssertNilError(t, err, "DB.Close") {
		return
	}
	dbSchema := db.DBSchema()
	testutils.AssertEqualWithLabel(t, false, dbSchema.HasTable(productAttributeTable.Name()), "DBSchema.HasTable")
}
