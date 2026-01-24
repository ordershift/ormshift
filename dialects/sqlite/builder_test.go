package sqlite_test

import (
	"database/sql"
	"testing"

	"github.com/ordershift/ormshift/dialects/sqlite"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestInteroperateSQLCommandWithNamedArgs(t *testing.T) {
	lDriver := sqlite.SQLiteDriver{}
	lReturnedSQL, lReturnedValue := lDriver.SQLBuilder().InteroperateSQLCommandWithNamedArgs("select * from table where id = @id", sql.Named("id", 1))
	lExpectedSQL := "select * from table where id = @id"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, lDriver.Name()+".InteroperateSQLCommandWithNamedArgs.SQL")
	testutils.AssertEqualWithLabel(t, "id", lReturnedValue[0].(sql.NamedArg).Name, lDriver.Name()+".InteroperateSQLCommandWithNamedArgs.Name")
	testutils.AssertEqualWithLabel(t, 1, lReturnedValue[0].(sql.NamedArg).Value.(int), lDriver.Name()+".InteroperateSQLCommandWithNamedArgs.Value")
}

func TestCreateTable(t *testing.T) {
	lSQLBuilder := sqlite.SQLiteDriver{}.SQLBuilder()

	lUserTable := testutils.FakeUserTable(t)
	lExpectedSQL := "CREATE TABLE user (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,email TEXT NOT NULL,name TEXT NOT NULL," +
		"password_hash TEXT,active INTEGER,created_at DATETIME,user_master INTEGER,master_user_id INTEGER,licence_price REAL,relevance REAL,photo BLOB,any TEXT);"
	lReturnedSQL := lSQLBuilder.CreateTable(*lUserTable)
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.CreateTable")

	lProductAttributeTable := testutils.FakeProductAttributeTable(t)
	lExpectedSQL = "CREATE TABLE product_attribute (product_id INTEGER NOT NULL,attribute_id INTEGER NOT NULL,value TEXT,position INTEGER,CONSTRAINT PK_product_attribute PRIMARY KEY (product_id,attribute_id));"
	lReturnedSQL = lSQLBuilder.CreateTable(*lProductAttributeTable)
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.CreateTable")
}
