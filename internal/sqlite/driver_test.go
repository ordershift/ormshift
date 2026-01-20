package sqlite_test

import (
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/internal/sqlite"
	"github.com/ordershift/ormshift/internal/testutils"
)

func Test_DriverSQLite_ConnectionString_ShouldBeValid(t *testing.T) {
	lReturnedConnectionString := sqlite.ConnectionString(ormshift.ConnectionParams{
		User:     "user",
		Password: "123456",
		Database: "my-db",
	})
	lExpectedConnectionString := "file:my-db.db?_auth&_auth_user=user&_auth_pass=123456&_locking=NORMAL"
	testutils.AssertEqualWithLabel(t, lExpectedConnectionString, lReturnedConnectionString, "DriverSQLite.ConnectionString")
}

func Test_DriverSQLite_ConnectionString_ShouldBeValid_WhenInMemory(t *testing.T) {
	lReturnedConnectionString := sqlite.ConnectionString(ormshift.ConnectionParams{InMemory: true})
	lExpectedConnectionString := ":memory:"
	testutils.AssertEqualWithLabel(t, lExpectedConnectionString, lReturnedConnectionString, "DriverSQLite.ConnectionString")
}
