package sqlite_test

import (
	"database/sql"
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/sqlite"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestInteroperateSQLCommandWithNamedArgs(t *testing.T) {
	driver := sqlite.Driver()
	returnedSQL, returnedValue := driver.SQLBuilder().InteroperateSQLCommandWithNamedArgs("select * from table where id = @id", sql.Named("id", 1))
	expectedSQL := "select * from table where id = @id"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, driver.Name()+".InteroperateSQLCommandWithNamedArgs.SQL")
	testutils.AssertEqualWithLabel(t, "id", returnedValue[0].(sql.NamedArg).Name, driver.Name()+".InteroperateSQLCommandWithNamedArgs.Name")
	testutils.AssertEqualWithLabel(t, 1, returnedValue[0].(sql.NamedArg).Value.(int), driver.Name()+".InteroperateSQLCommandWithNamedArgs.Value")
}

func TestCreateTable(t *testing.T) {
	sqlBuilder := sqlite.Driver().SQLBuilder()

	userTable := testutils.FakeUserTable(t)
	expectedSQL := "CREATE TABLE \"user\" (\"id\" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,\"email\" TEXT NOT NULL,\"name\" TEXT NOT NULL," +
		"\"password_hash\" TEXT,\"active\" INTEGER,\"created_at\" DATETIME,\"user_master\" INTEGER,\"master_user_id\" INTEGER,\"licence_price\" REAL,\"relevance\" REAL,\"photo\" BLOB,\"any\" TEXT);"
	returnedSQL := sqlBuilder.CreateTable(userTable)
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.CreateTable")

	productAttributeTable := testutils.FakeProductAttributeTable(t)
	expectedSQL = "CREATE TABLE \"product_attribute\" (\"product_id\" INTEGER NOT NULL,\"attribute_id\" INTEGER NOT NULL,\"value\" TEXT,\"position\" INTEGER,CONSTRAINT \"PK_product_attribute\" PRIMARY KEY (\"product_id\",\"attribute_id\"));"
	returnedSQL = sqlBuilder.CreateTable(productAttributeTable)
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.CreateTable")
}

func TestDropTable(t *testing.T) {
	sqlBuilder := sqlite.Driver().SQLBuilder()

	userTableName := testutils.FakeUserTableName(t)
	expectedSQL := "DROP TABLE \"user\";"
	returnedSQL := sqlBuilder.DropTable(userTableName)
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.DropTable")
}

func TestAlterTableAddColumn(t *testing.T) {
	sqlBuilder := sqlite.Driver().SQLBuilder()

	userTableName := testutils.FakeUserTableName(t)
	updatedAtColumn := testutils.FakeUpdatedAtColumn(t)
	expectedSQL := "ALTER TABLE \"user\" ADD COLUMN \"updated_at\" DATETIME;"
	returnedSQL := sqlBuilder.AlterTableAddColumn(userTableName, updatedAtColumn)
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.AlterTableAddColumn")
}

func TestAlterTableDropColumn(t *testing.T) {
	sqlBuilder := sqlite.Driver().SQLBuilder()

	userTableName := testutils.FakeUserTableName(t)
	updatedAtColumnName := testutils.FakeUpdatedAtColumnName(t)
	expectedSQL := "ALTER TABLE \"user\" DROP COLUMN \"updated_at\";"
	returnedSQL := sqlBuilder.AlterTableDropColumn(userTableName, updatedAtColumnName)
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.AlterTableDropColumn")
}

func TestInsert(t *testing.T) {
	sqlBuilder := sqlite.Driver().SQLBuilder()

	returnedSQL := sqlBuilder.Insert("product", []string{"id", "sku", "name", "description"})
	expectedSQL := "insert into \"product\" (\"id\",\"sku\",\"name\",\"description\") values (@id,@sku,@name,@description)"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.Insert")
}

func TestInsertWithValues(t *testing.T) {
	sqlBuilder := sqlite.Driver().SQLBuilder()

	returnedSQL, returnedValues := sqlBuilder.InsertWithValues("product", ormshift.ColumnsValues{"id": 1, "sku": "1.005.12.9", "name": "Trufa Sabor Amarula 30g Cacaushow"})
	expectedSQL := "insert into \"product\" (\"id\",\"name\",\"sku\") values (@id,@name,@sku)"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.InsertWithValues.SQL")
	testutils.AssertEqualWithLabel(t, 3, len(returnedValues), "SQLBuilder.InsertWithValues.Values")
	testutils.AssertNamedArgEqualWithLabel(t, returnedValues[0], sql.NamedArg{Name: "id", Value: 1}, "SQLBuilder.InsertWithValues.Values[0]")
	testutils.AssertNamedArgEqualWithLabel(t, returnedValues[1], sql.NamedArg{Name: "name", Value: "Trufa Sabor Amarula 30g Cacaushow"}, "SQLBuilder.InsertWithValues.Values[1]")
	testutils.AssertNamedArgEqualWithLabel(t, returnedValues[2], sql.NamedArg{Name: "sku", Value: "1.005.12.9"}, "SQLBuilder.InsertWithValues.Values[2]")
}

func TestUpdate(t *testing.T) {
	sqlBuilder := sqlite.Driver().SQLBuilder()

	returnedSQL := sqlBuilder.Update("product", []string{"sku", "name", "description"}, []string{"id"})
	expectedSQL := "update \"product\" set \"sku\" = @sku,\"name\" = @name,\"description\" = @description where \"id\" = @id"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.Update")
}

func TestUpdateWithValues(t *testing.T) {
	sqlBuilder := sqlite.Driver().SQLBuilder()

	returnedSQL, returnedValues := sqlBuilder.UpdateWithValues("product", []string{"sku", "name"}, []string{"id"}, ormshift.ColumnsValues{"id": 1, "sku": "1.005.12.5", "name": "Trufa Sabor Amarula 18g Cacaushow"})
	expectedSQL := "update \"product\" set \"sku\" = @sku,\"name\" = @name where \"id\" = @id"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.UpdateWithValues.SQL")
	testutils.AssertEqualWithLabel(t, 3, len(returnedValues), "SQLBuilder.UpdateWithValues.Values")
	testutils.AssertNamedArgEqualWithLabel(t, returnedValues[0], sql.NamedArg{Name: "id", Value: 1}, "SQLBuilder.UpdateWithValues.Values[0]")
	testutils.AssertNamedArgEqualWithLabel(t, returnedValues[1], sql.NamedArg{Name: "name", Value: "Trufa Sabor Amarula 18g Cacaushow"}, "SQLBuilder.UpdateWithValues.Values[1]")
	testutils.AssertNamedArgEqualWithLabel(t, returnedValues[2], sql.NamedArg{Name: "sku", Value: "1.005.12.5"}, "SQLBuilder.UpdateWithValues.Values[2]")
}

func TestDelete(t *testing.T) {
	sqlBuilder := sqlite.Driver().SQLBuilder()

	returnedSQL := sqlBuilder.Delete("product", []string{"id"})
	expectedSQL := "delete from \"product\" where \"id\" = @id"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.Delete")
}

func TestDeleteWithValues(t *testing.T) {
	sqlBuilder := sqlite.Driver().SQLBuilder()

	returnedSQL, returnedValues := sqlBuilder.DeleteWithValues("product", ormshift.ColumnsValues{"id": 1})
	expectedSQL := "delete from \"product\" where \"id\" = @id"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.DeleteWithValues.SQL")
	testutils.AssertEqualWithLabel(t, 1, len(returnedValues), "SQLBuilder.DeleteWithValues.Values")
	testutils.AssertNamedArgEqualWithLabel(t, returnedValues[0], sql.NamedArg{Name: "id", Value: 1}, "SQLBuilder.DeleteWithValues.Values[0]")
}

func TestSelect(t *testing.T) {
	sqlBuilder := sqlite.Driver().SQLBuilder()

	returnedSQL := sqlBuilder.Select("product", []string{"id", "name", "description"}, []string{"sku", "active"})
	expectedSQL := "select \"id\",\"name\",\"description\" from \"product\" where \"sku\" = @sku and \"active\" = @active"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.Select")
}

func TestSelectWithValues(t *testing.T) {
	sqlBuilder := sqlite.Driver().SQLBuilder()

	returnedSQL, returnedValues := sqlBuilder.SelectWithValues("product", []string{"id", "sku", "name", "description"}, ormshift.ColumnsValues{"category_id": 1, "active": true})
	expectedSQL := "select \"id\",\"sku\",\"name\",\"description\" from \"product\" where \"active\" = @active and \"category_id\" = @category_id"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.SelectWithValues.SQL")
	testutils.AssertEqualWithLabel(t, 2, len(returnedValues), "SQLBuilder.SelectWithValues.Values")
	testutils.AssertNamedArgEqualWithLabel(t, returnedValues[0], sql.NamedArg{Name: "active", Value: true}, "SQLBuilder.SelectWithValues.Values[0]")
	testutils.AssertNamedArgEqualWithLabel(t, returnedValues[1], sql.NamedArg{Name: "category_id", Value: 1}, "SQLBuilder.SelectWithValues.Values[1]")
}

func TestSelectWithPagination(t *testing.T) {
	sqlBuilder := sqlite.Driver().SQLBuilder()

	returnedSQL := sqlBuilder.SelectWithPagination("select * from product", 10, 5)
	expectedSQL := "select * from product LIMIT 10 OFFSET 40"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.SelectWithPagination")
}
