package sqlite_test

import (
	"database/sql"
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/sqlite"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestInteroperateSQLCommandWithNamedArgs(t *testing.T) {
	lDriver := sqlite.Driver()
	lReturnedSQL, lReturnedValue := lDriver.SQLBuilder().InteroperateSQLCommandWithNamedArgs("select * from table where id = @id", sql.Named("id", 1))
	lExpectedSQL := "select * from table where id = @id"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, lDriver.Name()+".InteroperateSQLCommandWithNamedArgs.SQL")
	testutils.AssertEqualWithLabel(t, "id", lReturnedValue[0].(sql.NamedArg).Name, lDriver.Name()+".InteroperateSQLCommandWithNamedArgs.Name")
	testutils.AssertEqualWithLabel(t, 1, lReturnedValue[0].(sql.NamedArg).Value.(int), lDriver.Name()+".InteroperateSQLCommandWithNamedArgs.Value")
}

func TestCreateTable(t *testing.T) {
	lSQLBuilder := sqlite.Driver().SQLBuilder()

	lUserTable := testutils.FakeUserTable(t)
	lExpectedSQL := "CREATE TABLE \"user\" (\"id\" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,\"email\" TEXT NOT NULL,\"name\" TEXT NOT NULL," +
		"\"password_hash\" TEXT,\"active\" INTEGER,\"created_at\" DATETIME,\"user_master\" INTEGER,\"master_user_id\" INTEGER,\"licence_price\" REAL,\"relevance\" REAL,\"photo\" BLOB,\"any\" TEXT);"
	lReturnedSQL := lSQLBuilder.CreateTable(lUserTable)
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.CreateTable")

	lProductAttributeTable := testutils.FakeProductAttributeTable(t)
	lExpectedSQL = "CREATE TABLE \"product_attribute\" (\"product_id\" INTEGER NOT NULL,\"attribute_id\" INTEGER NOT NULL,\"value\" TEXT,\"position\" INTEGER,CONSTRAINT \"PK_product_attribute\" PRIMARY KEY (\"product_id\",\"attribute_id\"));"
	lReturnedSQL = lSQLBuilder.CreateTable(lProductAttributeTable)
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.CreateTable")
}

func TestDropTable(t *testing.T) {
	lSQLBuilder := sqlite.Driver().SQLBuilder()

	lUserTableName := testutils.FakeUserTableName(t)
	lExpectedSQL := "DROP TABLE \"user\";"
	lReturnedSQL := lSQLBuilder.DropTable(lUserTableName)
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.DropTable")
}

func TestAlterTableAddColumn(t *testing.T) {
	lSQLBuilder := sqlite.Driver().SQLBuilder()

	lUserTableName := testutils.FakeUserTableName(t)
	lUpdatedAtColumn := testutils.FakeUpdatedAtColumn(t)
	lExpectedSQL := "ALTER TABLE \"user\" ADD COLUMN \"updated_at\" DATETIME;"
	lReturnedSQL := lSQLBuilder.AlterTableAddColumn(lUserTableName, lUpdatedAtColumn)
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.AlterTableAddColumn")
}

func TestAlterTableDropColumn(t *testing.T) {
	lSQLBuilder := sqlite.Driver().SQLBuilder()

	lUserTableName := testutils.FakeUserTableName(t)
	lUpdatedAtColumnName := testutils.FakeUpdatedAtColumnName(t)
	lExpectedSQL := "ALTER TABLE \"user\" DROP COLUMN \"updated_at\";"
	lReturnedSQL := lSQLBuilder.AlterTableDropColumn(lUserTableName, lUpdatedAtColumnName)
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.AlterTableDropColumn")
}

func TestInsert(t *testing.T) {
	lSQLBuilder := sqlite.Driver().SQLBuilder()

	lReturnedSQL := lSQLBuilder.Insert("product", []string{"id", "sku", "name", "description"})
	lExpectedSQL := "insert into \"product\" (\"id\",\"sku\",\"name\",\"description\") values (@id,@sku,@name,@description)"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.Insert")
}

func TestInsertWithValues(t *testing.T) {
	lSQLBuilder := sqlite.Driver().SQLBuilder()

	lReturnedSQL, lReturnedValues := lSQLBuilder.InsertWithValues("product", ormshift.ColumnsValues{"id": 1, "sku": "1.005.12.9", "name": "Trufa Sabor Amarula 30g Cacaushow"})
	lExpectedSQL := "insert into \"product\" (\"id\",\"name\",\"sku\") values (@id,@name,@sku)"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.InsertWithValues.SQL")
	testutils.AssertEqualWithLabel(t, 3, len(lReturnedValues), "SQLBuilder.InsertWithValues.Values")
	testutils.AssertNamedArgEqualWithLabel(t, lReturnedValues[0], sql.NamedArg{Name: "id", Value: 1}, "SQLBuilder.InsertWithValues.Values[0]")
	testutils.AssertNamedArgEqualWithLabel(t, lReturnedValues[1], sql.NamedArg{Name: "name", Value: "Trufa Sabor Amarula 30g Cacaushow"}, "SQLBuilder.InsertWithValues.Values[1]")
	testutils.AssertNamedArgEqualWithLabel(t, lReturnedValues[2], sql.NamedArg{Name: "sku", Value: "1.005.12.9"}, "SQLBuilder.InsertWithValues.Values[2]")
}

func TestUpdate(t *testing.T) {
	lSQLBuilder := sqlite.Driver().SQLBuilder()

	lReturnedSQL := lSQLBuilder.Update("product", []string{"sku", "name", "description"}, []string{"id"})
	lExpectedSQL := "update \"product\" set \"sku\" = @sku,\"name\" = @name,\"description\" = @description where \"id\" = @id"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.Update")
}

func TestUpdateWithValues(t *testing.T) {
	lSQLBuilder := sqlite.Driver().SQLBuilder()

	lReturnedSQL, lReturnedValues := lSQLBuilder.UpdateWithValues("product", []string{"sku", "name"}, []string{"id"}, ormshift.ColumnsValues{"id": 1, "sku": "1.005.12.5", "name": "Trufa Sabor Amarula 18g Cacaushow"})
	lExpectedSQL := "update \"product\" set \"sku\" = @sku,\"name\" = @name where \"id\" = @id"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.UpdateWithValues.SQL")
	testutils.AssertEqualWithLabel(t, 3, len(lReturnedValues), "SQLBuilder.UpdateWithValues.Values")
	testutils.AssertNamedArgEqualWithLabel(t, lReturnedValues[0], sql.NamedArg{Name: "id", Value: 1}, "SQLBuilder.UpdateWithValues.Values[0]")
	testutils.AssertNamedArgEqualWithLabel(t, lReturnedValues[1], sql.NamedArg{Name: "name", Value: "Trufa Sabor Amarula 18g Cacaushow"}, "SQLBuilder.UpdateWithValues.Values[1]")
	testutils.AssertNamedArgEqualWithLabel(t, lReturnedValues[2], sql.NamedArg{Name: "sku", Value: "1.005.12.5"}, "SQLBuilder.UpdateWithValues.Values[2]")
}

func TestDelete(t *testing.T) {
	lSQLBuilder := sqlite.Driver().SQLBuilder()

	lReturnedSQL := lSQLBuilder.Delete("product", []string{"id"})
	lExpectedSQL := "delete from \"product\" where \"id\" = @id"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.Delete")
}

func TestDeleteWithValues(t *testing.T) {
	lSQLBuilder := sqlite.Driver().SQLBuilder()

	lReturnedSQL, lReturnedValues := lSQLBuilder.DeleteWithValues("product", ormshift.ColumnsValues{"id": 1})
	lExpectedSQL := "delete from \"product\" where \"id\" = @id"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.DeleteWithValues.SQL")
	testutils.AssertEqualWithLabel(t, 1, len(lReturnedValues), "SQLBuilder.DeleteWithValues.Values")
	testutils.AssertNamedArgEqualWithLabel(t, lReturnedValues[0], sql.NamedArg{Name: "id", Value: 1}, "SQLBuilder.DeleteWithValues.Values[0]")
}

func TestSelect(t *testing.T) {
	lSQLBuilder := sqlite.Driver().SQLBuilder()

	lReturnedSQL := lSQLBuilder.Select("product", []string{"id", "name", "description"}, []string{"sku", "active"})
	lExpectedSQL := "select \"id\",\"name\",\"description\" from \"product\" where \"sku\" = @sku and \"active\" = @active"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.Select")
}

func TestSelectWithValues(t *testing.T) {
	lSQLBuilder := sqlite.Driver().SQLBuilder()

	lReturnedSQL, lReturnedValues := lSQLBuilder.SelectWithValues("product", []string{"id", "sku", "name", "description"}, ormshift.ColumnsValues{"category_id": 1, "active": true})
	lExpectedSQL := "select \"id\",\"sku\",\"name\",\"description\" from \"product\" where \"active\" = @active and \"category_id\" = @category_id"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.SelectWithValues.SQL")
	testutils.AssertEqualWithLabel(t, 2, len(lReturnedValues), "SQLBuilder.SelectWithValues.Values")
	testutils.AssertNamedArgEqualWithLabel(t, lReturnedValues[0], sql.NamedArg{Name: "active", Value: true}, "SQLBuilder.SelectWithValues.Values[0]")
	testutils.AssertNamedArgEqualWithLabel(t, lReturnedValues[1], sql.NamedArg{Name: "category_id", Value: 1}, "SQLBuilder.SelectWithValues.Values[1]")
}

func TestSelectWithPagination(t *testing.T) {
	lSQLBuilder := sqlite.Driver().SQLBuilder()

	lReturnedSQL := lSQLBuilder.SelectWithPagination("select * from product", 10, 5)
	lExpectedSQL := "select * from product LIMIT 10 OFFSET 40"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.SelectWithPagination")
}
