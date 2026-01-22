package sqlite_test

import (
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/sqlite"
	"github.com/ordershift/ormshift/internal/testutils"
)

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
