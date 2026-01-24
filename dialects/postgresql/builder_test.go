package postgresql_test

import (
	"database/sql"
	"testing"

	"github.com/ordershift/ormshift/dialects/postgresql"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestInteroperateSQLCommandWithNamedArgs(t *testing.T) {
	lDriver := &postgresql.PostgreSQLDriver{}
	lReturnedSQL, lReturnedValue := lDriver.SQLBuilder().InteroperateSQLCommandWithNamedArgs(
		"select * from user where role = @role and active = @active and master = @master",
		sql.Named("role", "admin"),
		sql.Named("active", true),
		sql.Named("master", false),
	)
	lExpectedSQL := "select * from user where role = $1 and active = $2 and master = $3"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "DriverPostgresql.InteroperateSQLCommandWithNamedArgs.SQL")
	testutils.AssertEqualWithLabel(t, "admin", lReturnedValue[0].(string), "DriverPostgresql.InteroperateSQLCommandWithNamedArgs.Value.1")
	testutils.AssertEqualWithLabel(t, 1, lReturnedValue[1].(int), "DriverPostgresql.InteroperateSQLCommandWithNamedArgs.Value.2")
	testutils.AssertEqualWithLabel(t, 0, lReturnedValue[2].(int), "DriverPostgresql.InteroperateSQLCommandWithNamedArgs.Value.3")
	lReturnedSQL, lReturnedValue = lDriver.SQLBuilder().InteroperateSQLCommandWithNamedArgs(
		"update user set role = @role where role = @role",
		sql.Named("role", "admin"),
	)
	lExpectedSQL = "update user set role = $1 where role = $1"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "DriverPostgresql.InteroperateSQLCommandWithNamedArgs.SQL")
	testutils.AssertEqualWithLabel(t, "admin", lReturnedValue[0].(string), "DriverPostgresql.InteroperateSQLCommandWithNamedArgs.Value.1")
	lReturnedSQL, lReturnedValue = lDriver.SQLBuilder().InteroperateSQLCommandWithNamedArgs(
		"delete from user where id = @id",
		sql.Named("role", "admin"),
	)
	lExpectedSQL = "delete from user where id = @id"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "DriverPostgresql.InteroperateSQLCommandWithNamedArgs.SQL")
	testutils.AssertEqualWithLabel(t, "admin", lReturnedValue[0].(string), "DriverPostgresql.InteroperateSQLCommandWithNamedArgs.Value.1")
}
