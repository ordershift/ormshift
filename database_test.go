package ormshift_test

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/postgresql"
	"github.com/ordershift/ormshift/dialects/sqlite"
	"github.com/ordershift/ormshift/dialects/sqlserver"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestOpenDatabase(t *testing.T) {
	lDB, lError := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNotNilResultAndNilError(t, lDB, lError, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = lDB.Close() }()

	testutils.AssertNotNilResultAndNilError(t, lDB.DB(), nil, "Database.DB")
	testutils.AssertEqualWithLabel(t, ":memory:", lDB.ConnectionString(), "Database.ConnectionString")
	testutils.AssertEqualWithLabel(t, "sqlite", lDB.DriverName(), "Database.ConnectionString")

	lUnderlyingDB := lDB.DB()
	if !testutils.AssertNotNilResultAndNilError(t, lUnderlyingDB, nil, "Database.DB") {
		return
	}
}

func TestOpenDatabaseWithNilDriver(t *testing.T) {
	lDB, lError := ormshift.OpenDatabase(nil, ormshift.ConnectionParams{})
	if !testutils.AssertNilResultAndNotNilError(t, lDB, lError, "ormshift.OpenDatabase") {
		return
	}
	testutils.AssertErrorMessage(t, "DatabaseDriver cannot be nil", lError, "ormshift.OpenDatabase")
}

func TestOpenDatabaseWithBadDriver(t *testing.T) {
	lDriver := testutils.NewFakeDriverBadName(sqlite.Driver())
	lDB, lError := ormshift.OpenDatabase(lDriver, ormshift.ConnectionParams{})
	if !testutils.AssertNilResultAndNotNilError(t, lDB, lError, "ormshift.OpenDatabase") {
		return
	}
	testutils.AssertErrorMessage(t, "sql.Open failed: sql: unknown driver \"bad-driver-name\" (forgotten import?)", lError, "ormshift.OpenDatabase")
}

func TestOpenDatabaseWithBadSchema(t *testing.T) {
	lDriver := testutils.NewFakeDriverBadSchema(sqlite.Driver())
	lDB, lError := ormshift.OpenDatabase(lDriver, ormshift.ConnectionParams{})
	if !testutils.AssertNilResultAndNotNilError(t, lDB, lError, "ormshift.OpenDatabase") {
		return
	}
	testutils.AssertErrorMessage(t, "failed to get DB schema: intentionally bad schema", lError, "ormshift.OpenDatabase")
}

func TestClose(t *testing.T) {
	lDB, lError := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNotNilResultAndNilError(t, lDB, lError, "ormshift.OpenDatabase") {
		return
	}

	lError = lDB.DB().Ping()
	if !testutils.AssertNilError(t, lError, "Database.DB.Ping") {
		_ = lDB.Close()
		return
	}

	lError = lDB.Close()
	testutils.AssertNilError(t, lError, "Database.Close")

	lError = lDB.DB().Ping()
	testutils.AssertErrorMessage(t, "sql: database is closed", lError, "Database.DB.Ping")
}

func TestValidateFailsWithInvalidConnectionString(t *testing.T) {
	lDriver := testutils.NewFakeDriverInvalidConnectionString(postgresql.Driver())
	lDB, lError := ormshift.OpenDatabase(lDriver, ormshift.ConnectionParams{})
	if !testutils.AssertNotNilResultAndNilError(t, lDB, lError, "ormshift.OpenDatabase") {
		return
	}
	lError = lDB.Validate()
	testutils.AssertErrorMessage(t, "database ping failed: missing \"=\" after \"invalid-connection-string\" in connection info string\"", lError, "ormshift.OpenDatabase")
}

func TestConnectionStringWithNoParams(t *testing.T) {
	lConnectionParams := ormshift.ConnectionParams{InMemory: true}
	lDB, lError := ormshift.OpenDatabase(sqlite.Driver(), lConnectionParams)
	if !testutils.AssertNotNilResultAndNilError(t, lDB, lError, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = lDB.Close() }()

	lExpectedConnectionString := ":memory:"
	lReturnedConnectionString := lDB.ConnectionString()
	testutils.AssertEqualWithLabel(t, lExpectedConnectionString, lReturnedConnectionString, "Database.ConnectionString")

	// Connection string should not be modified if the connection params is changed
	lConnectionParams.InMemory = false
	lReturnedConnectionString = lDB.ConnectionString()
	testutils.AssertEqualWithLabel(t, lExpectedConnectionString, lReturnedConnectionString, "Database.ConnectionString")
}

func TestDriverName(t *testing.T) {
	lDriver := testutils.NewFakeDriver(sqlite.Driver())
	testutils.AssertEqualWithLabel(t, "sqlite", lDriver.Name(), "FakeDriver.Name")
}

func TestDriverConnectionString(t *testing.T) {
	lDriver := testutils.NewFakeDriver(sqlserver.Driver())
	lConnectionParams := ormshift.ConnectionParams{
		Host:     "localhost",
		Port:     1433,
		User:     "sa",
		Password: "your_password",
		Database: "testdb",
	}

	lReturnedConnectionString := lDriver.ConnectionString(lConnectionParams)
	lExpectedConnectionString := "server=localhost;port=1433;user id=sa;password=your_password;database=testdb"
	testutils.AssertEqualWithLabel(t, lExpectedConnectionString, lReturnedConnectionString, "FakeDriver.ConnectionString")

	lConnectionParams = ormshift.ConnectionParams{
		Instance: "myServerName\\myInstanceName",
		User:     "sa",
		Password: "your_password",
		Database: "testdb",
	}
	lReturnedConnectionString = lDriver.ConnectionString(lConnectionParams)
	lExpectedConnectionString = "server=\\myServerName\\myInstanceName;user id=sa;password=your_password;database=testdb"
	testutils.AssertEqualWithLabel(t, lExpectedConnectionString, lReturnedConnectionString, "FakeDriver.ConnectionString")
}

func TestDriverSQLBuilder(t *testing.T) {
	lDriver := testutils.NewFakeDriver(sqlite.Driver())
	lSQLBuilder := lDriver.SQLBuilder()
	testutils.AssertEqualWithLabel(t, "sqliteSQLBuilder", reflect.TypeOf(lSQLBuilder).Name(), "FakeDriver.SQLBuilder")
}

func TestDriverDBSchema(t *testing.T) {
	lDriver := testutils.NewFakeDriver(sqlite.Driver())
	lDB, lError := sql.Open(lDriver.Name(), lDriver.ConnectionString(ormshift.ConnectionParams{InMemory: true}))
	if !testutils.AssertNotNilResultAndNilError(t, lDB, lError, "sql.Open") {
		return
	}
	defer func() { _ = lDB.Close() }()

	lSchema, lError := lDriver.DBSchema(lDB)
	if !testutils.AssertNotNilResultAndNilError(t, lSchema, lError, "FakeDriver.DBSchema") {
		return
	}
}
