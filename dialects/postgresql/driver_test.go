package postgresql_test

import (
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/postgresql"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestConnectionString(t *testing.T) {
	lDriver := postgresql.PostgreSQLDriver{}
	lReturnedConnectionString := lDriver.ConnectionString(ormshift.ConnectionParams{
		User:     "pg",
		Password: "123456",
		Database: "my-db",
	})
	lExpectedConnectionString := "host=localhost port=5432 user=pg password=123456 dbname=my-db sslmode=disable"
	testutils.AssertEqualWithLabel(t, lExpectedConnectionString, lReturnedConnectionString, "lDriver.ConnectionString")
}
