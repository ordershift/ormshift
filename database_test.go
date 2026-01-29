package ormshift_test

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/sqlite"
	"github.com/ordershift/ormshift/dialects/sqlserver"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestOpenDatabase(t *testing.T) {
	db, err := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNotNilResultAndNilError(t, db, err, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = db.Close() }()

	testutils.AssertNotNilResultAndNilError(t, db.DB(), nil, "Database.DB")
	testutils.AssertEqualWithLabel(t, ":memory:", db.ConnectionString(), "Database.ConnectionString")
	testutils.AssertEqualWithLabel(t, "sqlite", db.DriverName(), "Database.ConnectionString")

	underlyingDB := db.DB()
	if !testutils.AssertNotNilResultAndNilError(t, underlyingDB, nil, "Database.DB") {
		return
	}
}

func TestOpenDatabaseWithNilDriver(t *testing.T) {
	db, err := ormshift.OpenDatabase(nil, ormshift.ConnectionParams{})
	if !testutils.AssertNilResultAndNotNilError(t, db, err, "ormshift.OpenDatabase") {
		return
	}
	testutils.AssertErrorMessage(t, "DatabaseDriver cannot be nil", err, "ormshift.OpenDatabase")
}

func TestOpenDatabaseWithBadDriver(t *testing.T) {
	driver := testutils.NewFakeDriverBadName(sqlite.Driver())
	db, err := ormshift.OpenDatabase(driver, ormshift.ConnectionParams{})
	if !testutils.AssertNilResultAndNotNilError(t, db, err, "ormshift.OpenDatabase") {
		return
	}
	testutils.AssertErrorMessage(t, "sql.Open failed: sql: unknown driver \"bad-driver-name\" (forgotten import?)", err, "ormshift.OpenDatabase")
}

func TestOpenDatabaseWithBadSchema(t *testing.T) {
	driver := testutils.NewFakeDriverBadSchema(sqlite.Driver())
	db, err := ormshift.OpenDatabase(driver, ormshift.ConnectionParams{})
	if !testutils.AssertNilResultAndNotNilError(t, db, err, "ormshift.OpenDatabase") {
		return
	}
	testutils.AssertErrorMessage(t, "failed to get DB schema: intentionally bad schema", err, "ormshift.OpenDatabase")
}

func TestClose(t *testing.T) {
	db, err := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNotNilResultAndNilError(t, db, err, "ormshift.OpenDatabase") {
		return
	}

	err = db.DB().Ping()
	if !testutils.AssertNilError(t, err, "Database.DB.Ping") {
		_ = db.Close()
		return
	}

	err = db.Close()
	testutils.AssertNilError(t, err, "Database.Close")

	err = db.DB().Ping()
	testutils.AssertErrorMessage(t, "sql: database is closed", err, "Database.DB.Ping")
}

func TestConnectionStringWithNoParams(t *testing.T) {
	connectionParams := ormshift.ConnectionParams{InMemory: true}
	db, err := ormshift.OpenDatabase(sqlite.Driver(), connectionParams)
	if !testutils.AssertNotNilResultAndNilError(t, db, err, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = db.Close() }()

	expectedConnectionString := ":memory:"
	returnedConnectionString := db.ConnectionString()
	testutils.AssertEqualWithLabel(t, expectedConnectionString, returnedConnectionString, "Database.ConnectionString")

	// Connection string should not be modified if the connection params is changed
	connectionParams.InMemory = false
	returnedConnectionString = db.ConnectionString()
	testutils.AssertEqualWithLabel(t, expectedConnectionString, returnedConnectionString, "Database.ConnectionString")
}

func TestDriverName(t *testing.T) {
	driver := testutils.NewFakeDriver(sqlite.Driver())
	testutils.AssertEqualWithLabel(t, "sqlite", driver.Name(), "FakeDriver.Name")
}

func TestDriverConnectionString(t *testing.T) {
	driver := testutils.NewFakeDriver(sqlserver.Driver())
	connectionParams := ormshift.ConnectionParams{
		Host:     "localhost",
		Port:     1433,
		User:     "sa",
		Password: "your_password",
		Database: "testdb",
	}

	returnedConnectionString := driver.ConnectionString(connectionParams)
	expectedConnectionString := "server=localhost;port=1433;user id=sa;password=your_password;database=testdb"
	testutils.AssertEqualWithLabel(t, expectedConnectionString, returnedConnectionString, "FakeDriver.ConnectionString")

	connectionParams = ormshift.ConnectionParams{
		Instance: "myServerName\\myInstanceName",
		User:     "sa",
		Password: "your_password",
		Database: "testdb",
	}
	returnedConnectionString = driver.ConnectionString(connectionParams)
	expectedConnectionString = "server=\\myServerName\\myInstanceName;user id=sa;password=your_password;database=testdb"
	testutils.AssertEqualWithLabel(t, expectedConnectionString, returnedConnectionString, "FakeDriver.ConnectionString")
}

func TestDriverSQLBuilder(t *testing.T) {
	driver := testutils.NewFakeDriver(sqlite.Driver())
	sqlBuilder := driver.SQLBuilder()
	testutils.AssertEqualWithLabel(t, "sqliteBuilder", reflect.TypeOf(sqlBuilder).Elem().Name(), "FakeDriver.SQLBuilder")
}

func TestDriverDBSchema(t *testing.T) {
	driver := testutils.NewFakeDriver(sqlite.Driver())
	db, err := sql.Open(driver.Name(), driver.ConnectionString(ormshift.ConnectionParams{InMemory: true}))
	if !testutils.AssertNotNilResultAndNilError(t, db, err, "sql.Open") {
		return
	}
	defer func() { _ = db.Close() }()

	schema, err := driver.DBSchema(db)
	if !testutils.AssertNotNilResultAndNilError(t, schema, err, "FakeDriver.DBSchema") {
		return
	}
}
