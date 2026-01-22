package ormshift_test

import (
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/postgresql"
	"github.com/ordershift/ormshift/dialects/sqlite"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestOpenDatabase(t *testing.T) {
	lDB, lError := ormshift.OpenDatabase(&sqlite.SQLiteDriver{}, ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNotNilResultAndNilError(t, lDB, lError, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = lDB.Close() }()

	testutils.AssertNotNilResultAndNilError(t, lDB.DB(), nil, "Database.DB")
	testutils.AssertEqualWithLabel(t, ":memory:", lDB.ConnectionString(), "Database.ConnectionString")
	testutils.AssertEqualWithLabel(t, "sqlite", lDB.DriverName(), "Database.ConnectionString")
}

func TestOpenDatabaseWithNilDriver(t *testing.T) {
	lDB, lError := ormshift.OpenDatabase(nil, ormshift.ConnectionParams{})
	if !testutils.AssertNilResultAndNotNilError(t, lDB, lError, "ormshift.OpenDatabase") {
		return
	}
	testutils.AssertErrorMessage(t, "DatabaseDriver cannot be nil", lError, "ormshift.OpenDatabase")
}

func TestOpenDatabaseWithBadSchema(t *testing.T) {
	lDriver := testutils.NewFakeDriverBadSchema(&sqlite.SQLiteDriver{})
	lDB, lError := ormshift.OpenDatabase(lDriver, ormshift.ConnectionParams{})
	if !testutils.AssertNilResultAndNotNilError(t, lDB, lError, "ormshift.OpenDatabase") {
		return
	}
	testutils.AssertErrorMessage(t, "failed to get DB schema: intentionally bad schema", lError, "ormshift.OpenDatabase")
}

func TestValidateFailsWithInvalidConnectionString(t *testing.T) {
	lDriver := testutils.NewFakeDriverInvalidConnectionString(&postgresql.PostgreSQLDriver{})
	lDB, lError := ormshift.OpenDatabase(lDriver, ormshift.ConnectionParams{})
	if !testutils.AssertNotNilResultAndNilError(t, lDB, lError, "ormshift.OpenDatabase") {
		return
	}
	lError = lDB.Validate()
	testutils.AssertErrorMessage(t, "database ping failed: missing \"=\" after \"invalid-connection-string\" in connection info string\"", lError, "ormshift.OpenDatabase")
}
