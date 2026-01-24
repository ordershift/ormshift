package ormshift_test

import (
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/postgresql"
	"github.com/ordershift/ormshift/dialects/sqlite"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestOpenDatabase(t *testing.T) {
	lDB, lError := ormshift.OpenDatabase(sqlite.SQLiteDriver{}, ormshift.ConnectionParams{InMemory: true})
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
	lDriver := testutils.NewFakeDriverBadName(sqlite.SQLiteDriver{})
	lDB, lError := ormshift.OpenDatabase(lDriver, ormshift.ConnectionParams{})
	if !testutils.AssertNilResultAndNotNilError(t, lDB, lError, "ormshift.OpenDatabase") {
		return
	}
	testutils.AssertErrorMessage(t, "sql.Open failed: sql: unknown driver \"bad-driver-name\" (forgotten import?)", lError, "ormshift.OpenDatabase")
}

func TestOpenDatabaseWithBadSchema(t *testing.T) {
	lDriver := testutils.NewFakeDriverBadSchema(sqlite.SQLiteDriver{})
	lDB, lError := ormshift.OpenDatabase(lDriver, ormshift.ConnectionParams{})
	if !testutils.AssertNilResultAndNotNilError(t, lDB, lError, "ormshift.OpenDatabase") {
		return
	}
	testutils.AssertErrorMessage(t, "failed to get DB schema: intentionally bad schema", lError, "ormshift.OpenDatabase")
}

func TestClose(t *testing.T) {
	lDB, lError := ormshift.OpenDatabase(sqlite.SQLiteDriver{}, ormshift.ConnectionParams{InMemory: true})
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
	lDriver := testutils.NewFakeDriverInvalidConnectionString(postgresql.PostgreSQLDriver{})
	lDB, lError := ormshift.OpenDatabase(lDriver, ormshift.ConnectionParams{})
	if !testutils.AssertNotNilResultAndNilError(t, lDB, lError, "ormshift.OpenDatabase") {
		return
	}
	lError = lDB.Validate()
	testutils.AssertErrorMessage(t, "database ping failed: missing \"=\" after \"invalid-connection-string\" in connection info string\"", lError, "ormshift.OpenDatabase")
}

func TestConnectionStringWithNoParams(t *testing.T) {
	lConnectionParams := ormshift.ConnectionParams{InMemory: true}
	lDB, lError := ormshift.OpenDatabase(sqlite.SQLiteDriver{}, lConnectionParams)
	if !testutils.AssertNotNilResultAndNilError(t, lDB, lError, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = lDB.Close() }()

	testutils.AssertEqualWithLabel(t, ":memory:", lDB.ConnectionString(), "Database.ConnectionString")

	// Connection string should not be modified if the connection params is changed
	lConnectionParams.InMemory = false
	testutils.AssertEqualWithLabel(t, ":memory:", lDB.ConnectionString(), "Database.ConnectionString")
}
