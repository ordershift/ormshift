package sqlserver_test

import (
	"database/sql"
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/sqlserver"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestName(t *testing.T) {
	lDriver := sqlserver.SQLServerDriver{}
	testutils.AssertEqualWithLabel(t, "sqlserver", lDriver.Name(), "lDriver.Name")
}

func TestConnectionString(t *testing.T) {
	lDriver := sqlserver.SQLServerDriver{}
	lReturnedConnectionString := lDriver.ConnectionString(ormshift.ConnectionParams{
		Host:     "my-server",
		Port:     1433,
		Instance: "sqlexpress",
		User:     "sa",
		Password: "123456",
		Database: "my-db",
	})
	lExpectedConnectionString := "server=my-server\\sqlexpress;port=1433;user id=sa;password=123456;database=my-db"
	testutils.AssertEqualWithLabel(t, lExpectedConnectionString, lReturnedConnectionString, "lDriver.ConnectionString")
}

func TestDBSchema(t *testing.T) {
	lDriver := sqlserver.SQLServerDriver{}
	lDB, lError := sql.Open(lDriver.Name(), "server=my-server\\sqlexpress;port=1433;user id=sa;password=123456;database=my-db")
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
	lDriver := sqlserver.SQLServerDriver{}
	lSchema, lError := lDriver.DBSchema(nil)
	if !testutils.AssertNilResultAndNotNilError(t, lSchema, lError, "lDriver.DBSchema") {
		return
	}
	testutils.AssertErrorMessage(t, "sql.DB cannot be nil", lError, "lDriver.DBSchema")
}
