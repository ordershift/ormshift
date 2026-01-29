package postgresql_test

import (
	"database/sql"
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/postgresql"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestName(t *testing.T) {
	driver := postgresql.Driver()
	testutils.AssertEqualWithLabel(t, "postgres", driver.Name(), "driver.Name")
}

func TestConnectionString(t *testing.T) {
	driver := postgresql.Driver()
	returnedConnectionString := driver.ConnectionString(ormshift.ConnectionParams{
		User:     "pg",
		Password: "123456",
		Database: "my-db",
	})
	expectedConnectionString := "host=localhost port=5432 user=pg password=123456 dbname=my-db sslmode=disable"
	testutils.AssertEqualWithLabel(t, expectedConnectionString, returnedConnectionString, "driver.ConnectionString")
}

func TestDBSchema(t *testing.T) {
	driver := postgresql.Driver()
	db, err := sql.Open(driver.Name(), "host=localhost port=5432 user=pg password=123456 dbname=my-db sslmode=disable")
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
	driver := postgresql.Driver()
	schema, err := driver.DBSchema(nil)
	if !testutils.AssertNilResultAndNotNilError(t, schema, err, "driver.DBSchema") {
		return
	}
	testutils.AssertErrorMessage(t, "sql.DB cannot be nil", err, "driver.DBSchema")
}
