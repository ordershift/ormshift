package postgresql_test

import (
	"database/sql"
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/postgresql"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestName(t *testing.T) {
	lDriver := postgresql.PostgreSQLDriver{}
	testutils.AssertEqualWithLabel(t, "postgres", lDriver.Name(), "lDriver.Name")
}

func TestConnectionString(t *testing.T) {
	lDriver := postgresql.PostgreSQLDriver{}
	lReturnedConnectionString := lDriver.ConnectionString(ormshift.ConnectionParams{
		User:     "pg",
		Password: "123456",
		Database: "my-db",
	})
	lExpectedConnectionString := "host=localhost port=5432 user=pg password=123456 dbname=my-db sslmode=disable"
	testutils.AssertEqualWithLabel(t, lExpectedConnectionString, lReturnedConnectionString, "lDriver.ConnectionString")
}

func TestDBSchema(t *testing.T) {
	lDriver := postgresql.PostgreSQLDriver{}
	lDB, lError := sql.Open(lDriver.Name(), "host=localhost port=5432 user=pg password=123456 dbname=my-db sslmode=disable")
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
	lDriver := postgresql.PostgreSQLDriver{}
	lSchema, lError := lDriver.DBSchema(nil)
	if !testutils.AssertNilResultAndNotNilError(t, lSchema, lError, "lDriver.DBSchema") {
		return
	}
	testutils.AssertErrorMessage(t, "sql.DB cannot be nil", lError, "lDriver.DBSchema")
}
