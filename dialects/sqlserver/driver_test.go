package sqlserver_test

import (
	"database/sql"
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/sqlserver"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestName(t *testing.T) {
	driver := sqlserver.Driver()
	testutils.AssertEqualWithLabel(t, "sqlserver", driver.Name(), "driver.Name")
}

func TestConnectionString(t *testing.T) {
	driver := sqlserver.Driver()
	returnedConnectionString := driver.ConnectionString(ormshift.ConnectionParams{
		Host:     "my-server",
		Port:     1433,
		Instance: "sqlexpress",
		User:     "sa",
		Password: "123456",
		Database: "my-db",
	})
	expectedConnectionString := "server=my-server\\sqlexpress;port=1433;user id=sa;password=123456;database=my-db"
	testutils.AssertEqualWithLabel(t, expectedConnectionString, returnedConnectionString, "driver.ConnectionString")
}

func TestDBSchema(t *testing.T) {
	driver := sqlserver.Driver()
	db, err := sql.Open(driver.Name(), "server=my-server\\sqlexpress;port=1433;user id=sa;password=123456;database=my-db")
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
	driver := sqlserver.Driver()
	schema, err := driver.DBSchema(nil)
	if !testutils.AssertNilResultAndNotNilError(t, schema, err, "driver.DBSchema") {
		return
	}
	testutils.AssertErrorMessage(t, "failed to get db schema: db cannot be nil", err, "driver.DBSchema")
}
