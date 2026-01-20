package postgresql_test

import (
	"testing"

	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/pkg/core"
	"github.com/ordershift/ormshift/pkg/dialects/postgresql"
)

func Test_DriverPostgresql_ConnectionString_ShouldBeValid(t *testing.T) {
	lReturnedConnectionString := postgresql.ConnectionString(core.ConnectionParams{
		User:     "pg",
		Password: "123456",
		Database: "my-db",
	})
	lExpectedConnectionString := "host=localhost port=5432 user=pg password=123456 dbname=my-db sslmode=disable"
	testutils.AssertEqualWithLabel(t, lExpectedConnectionString, lReturnedConnectionString, "DriverPostgresql.ConnectionString")
}
