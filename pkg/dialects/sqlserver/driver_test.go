package sqlserver_test

import (
	"testing"

	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/pkg/core"
	"github.com/ordershift/ormshift/pkg/dialects/sqlserver"
)

func Test_DriverSQLServer_ConnectionString_ShouldBeValid(t *testing.T) {
	lReturnedConnectionString := sqlserver.ConnectionString(core.ConnectionParams{
		Host:     "my-server",
		Port:     1433,
		Instance: "sqlexpress",
		User:     "sa",
		Password: "123456",
		Database: "my-db",
	})
	lExpectedConnectionString := "server=my-server\\sqlexpress;port=1433;user id=sa;password=123456;database=my-db" //NOSONAR go:S2068
	testutils.AssertEqualWithLabel(t, lExpectedConnectionString, lReturnedConnectionString, "DriverSQLServer.ConnectionString")
}
