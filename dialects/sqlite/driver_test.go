package sqlite_test

import (
	"database/sql"
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/sqlite"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestName(t *testing.T) {
	lDriver := sqlite.SQLiteDriver{}
	testutils.AssertEqualWithLabel(t, "sqlite", lDriver.Name(), "lDriver.Name")
}

func TestConnectionString(t *testing.T) {
	lDriver := sqlite.SQLiteDriver{}
	lReturnedConnectionString := lDriver.ConnectionString(ormshift.ConnectionParams{
		User:     "user",
		Password: "123456",
		Database: "my-db",
	})
	lExpectedConnectionString := "file:my-db.db?_auth&_auth_user=user&_auth_pass=123456&_locking=NORMAL"
	testutils.AssertEqualWithLabel(t, lExpectedConnectionString, lReturnedConnectionString, "lDriver.ConnectionString")
}

func TestConnectionStringInMemory(t *testing.T) {
	lDriver := sqlite.SQLiteDriver{}
	lReturnedConnectionString := lDriver.ConnectionString(ormshift.ConnectionParams{InMemory: true})
	lExpectedConnectionString := ":memory:"
	testutils.AssertEqualWithLabel(t, lExpectedConnectionString, lReturnedConnectionString, "lDriver.ConnectionString")
}

func TestDBSchema(t *testing.T) {
	lDriver := sqlite.SQLiteDriver{}
	lDB, lError := sql.Open(lDriver.Name(), lDriver.ConnectionString(ormshift.ConnectionParams{InMemory: true}))
	if !testutils.AssertNotNilResultAndNilError(t, lDB, lError, "sql.Open") {
		return
	}
	defer func() { _ = lDB.Close() }()

	lSchema, lError := lDriver.DBSchema(lDB)
	if !testutils.AssertNotNilResultAndNilError(t, lSchema, lError, "lDriver.DBSchema") {
		return
	}
}

func TestDBSchemaFailsWhenDBIsNil(t *testing.T) {
	lDriver := sqlite.SQLiteDriver{}
	lSchema, lError := lDriver.DBSchema(nil)
	if !testutils.AssertNilResultAndNotNilError(t, lSchema, lError, "lDriver.DBSchema") {
		return
	}
	testutils.AssertErrorMessage(t, "sql.DB cannot be nil", lError, "lDriver.DBSchema")
}
