package postgresql_test

import (
	"database/sql"
	"testing"

	"github.com/ordershift/ormshift/dialects/postgresql"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestInteroperateSQLCommandWithNamedArgs(t *testing.T) {
	lDriver := postgresql.PostgreSQLDriver{}
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

func TestCreateTable(t *testing.T) {
	lSQLBuilder := postgresql.PostgreSQLDriver{}.SQLBuilder()

	lUserTable := testutils.FakeUserTable(t)
	lExpectedSQL := "CREATE TABLE user (id BIGSERIAL NOT NULL,email VARCHAR(80) NOT NULL,name VARCHAR(50) NOT NULL," +
		"password_hash VARCHAR(256),active SMALLINT,created_at TIMESTAMP(6),user_master BIGINT,master_user_id BIGINT," +
		"licence_price NUMERIC(17,2),relevance DOUBLE PRECISION,photo BYTEA,any VARCHAR,PRIMARY KEY (id,email));"
	lReturnedSQL := lSQLBuilder.CreateTable(*lUserTable)
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.CreateTable")

	lProductAttributeTable := testutils.FakeProductAttributeTable(t)
	lExpectedSQL = "CREATE TABLE product_attribute (product_id BIGINT NOT NULL,attribute_id BIGINT NOT NULL,value VARCHAR(75),position BIGINT,PRIMARY KEY (product_id,attribute_id));"
	lReturnedSQL = lSQLBuilder.CreateTable(*lProductAttributeTable)
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.CreateTable")
}
