package sqlserver_test

import (
	"database/sql"
	"testing"

	"github.com/ordershift/ormshift/dialects/sqlserver"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestInteroperateSQLCommandWithNamedArgs(t *testing.T) {
	lDriver := sqlserver.SQLServerDriver{}
	lReturnedSQL, lReturnedValue := lDriver.SQLBuilder().InteroperateSQLCommandWithNamedArgs("select * from table where id = @id", sql.Named("id", 1))
	lExpectedSQL := "select * from table where id = @id"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, lDriver.Name()+".InteroperateSQLCommandWithNamedArgs.SQL")
	testutils.AssertEqualWithLabel(t, "id", lReturnedValue[0].(sql.NamedArg).Name, lDriver.Name()+".InteroperateSQLCommandWithNamedArgs.Name")
	testutils.AssertEqualWithLabel(t, 1, lReturnedValue[0].(sql.NamedArg).Value.(int), lDriver.Name()+".InteroperateSQLCommandWithNamedArgs.Value")
}

func TestCreateTable(t *testing.T) {
	lSQLBuilder := sqlserver.SQLServerDriver{}.SQLBuilder()

	lUserTable := testutils.FakeUserTable(t)
	lExpectedSQL := "CREATE TABLE user (id BIGINT NOT NULL IDENTITY (1, 1),email VARCHAR(80) NOT NULL,name VARCHAR(50) NOT NULL," +
		"password_hash VARCHAR(256),active BIT,created_at DATETIME2(6),user_master BIGINT,master_user_id BIGINT," +
		"licence_price MONEY,relevance FLOAT,photo VARBINARY(MAX),any VARCHAR,CONSTRAINT PK_user PRIMARY KEY (id,email));"
	lReturnedSQL := lSQLBuilder.CreateTable(*lUserTable)
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.CreateTable")

	lProductAttributeTable := testutils.FakeProductAttributeTable(t)
	lExpectedSQL = "CREATE TABLE product_attribute (product_id BIGINT NOT NULL,attribute_id BIGINT NOT NULL,value VARCHAR(75),position BIGINT,CONSTRAINT PK_product_attribute PRIMARY KEY (product_id,attribute_id));"
	lReturnedSQL = lSQLBuilder.CreateTable(*lProductAttributeTable)
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.CreateTable")
}
