package sqlserver_test

import (
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/sqlserver"
	"github.com/ordershift/ormshift/internal/testutils"
)

func Test_DriverSQLServer_ConnectionString_ShouldBeValid(t *testing.T) {
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
