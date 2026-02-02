package sqlite_test

import (
	"database/sql"
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/sqlite"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestName(t *testing.T) {
	driver := sqlite.Driver()
	testutils.AssertEqualWithLabel(t, "sqlite", driver.Name(), "driver.Name")
}

func TestConnectionString(t *testing.T) {
	driver := sqlite.Driver()
	returnedConnectionString := driver.ConnectionString(ormshift.ConnectionParams{
		User:     "user",
		Password: "123456",
		Database: "my-db",
	})
	expectedConnectionString := "file:my-db.db?_auth&_auth_user=user&_auth_pass=123456&_locking=NORMAL"
	testutils.AssertEqualWithLabel(t, expectedConnectionString, returnedConnectionString, "driver.ConnectionString")
}

func TestConnectionStringInMemory(t *testing.T) {
	driver := sqlite.Driver()
	returnedConnectionString := driver.ConnectionString(ormshift.ConnectionParams{InMemory: true})
	expectedConnectionString := ":memory:"
	testutils.AssertEqualWithLabel(t, expectedConnectionString, returnedConnectionString, "driver.ConnectionString")
}

func TestDBSchema(t *testing.T) {
	driver := sqlite.Driver()
	db, err := sql.Open(driver.Name(), driver.ConnectionString(ormshift.ConnectionParams{InMemory: true}))
	if !testutils.AssertNotNilResultAndNilError(t, db, err, "sql.Open") {
		return
	}
	defer func() { _ = db.Close() }()

	schema, err := driver.DBSchema(db)
	if !testutils.AssertNotNilResultAndNilError(t, schema, err, "driver.DBSchema") {
		return
	}
}

func TestDBSchemaFailsWhenDBIsNil(t *testing.T) {
	driver := sqlite.Driver()
	schema, err := driver.DBSchema(nil)
	if !testutils.AssertNilResultAndNotNilError(t, schema, err, "driver.DBSchema") {
		return
	}
	testutils.AssertErrorMessage(t, "failed to get db schema: db cannot be nil", err, "driver.DBSchema")
}
