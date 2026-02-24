package internal_test

import (
	"database/sql"
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/internal"
	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/schema"
)

func TestCreateTable(t *testing.T) {
	sqlBuilder := internal.NewGenericSQLBuilder(nil, nil, nil)

	userTable := testutils.FakeUserTable(t)
	expectedSQL := "CREATE TABLE \"user\" (\"id\" <<TYPE_0>>,\"email\" <<TYPE_1>>,\"name\" <<TYPE_1>>,\"password_hash\" <<TYPE_1>>," +
		"\"active\" <<TYPE_5>>,\"created_at\" <<TYPE_3>>,\"updated_at\" <<TYPE_7>>,\"user_master\" <<TYPE_0>>,\"master_user_id\" <<TYPE_0>>," +
		"\"licence_price\" <<TYPE_2>>,\"relevance\" <<TYPE_4>>,\"photo\" <<TYPE_6>>,\"any\" <<TYPE_-1>>,PRIMARY KEY (\"id\",\"email\"));"
	returnedSQL := sqlBuilder.CreateTable(userTable)
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.CreateTable")

	productAttributeTable := testutils.FakeProductAttributeTable(t)
	expectedSQL = "CREATE TABLE \"product_attribute\" (\"product_id\" <<TYPE_0>>,\"attribute_id\" <<TYPE_0>>,\"value\" <<TYPE_1>>,\"position\" <<TYPE_0>>,PRIMARY KEY (\"product_id\",\"attribute_id\"));"
	returnedSQL = sqlBuilder.CreateTable(productAttributeTable)
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.CreateTable")
}

func TestDropTable(t *testing.T) {
	sqlBuilder := internal.NewGenericSQLBuilder(nil, nil, nil)

	userTableName := testutils.FakeUserTableName(t)
	expectedSQL := "DROP TABLE \"user\";"
	returnedSQL := sqlBuilder.DropTable(userTableName)
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.DropTable")
}

func TestAlterTableAddColumn(t *testing.T) {
	sqlBuilder := internal.NewGenericSQLBuilder(nil, nil, nil)

	userTableName := testutils.FakeUserTableName(t)

	updatedAtColumn := testutils.FakeUpdatedAtColumn(t)
	expectedSQL := "ALTER TABLE \"user\" ADD COLUMN \"updated_at\" <<TYPE_3>>;"
	returnedSQL := sqlBuilder.AlterTableAddColumn(userTableName, updatedAtColumn)
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.AlterTableAddColumn")

	createdAtColumn := testutils.FakeCreatedAtColumn(t)
	expectedSQL = "ALTER TABLE \"user\" ADD COLUMN \"created_at\" <<TYPE_7>> DEFAULT '1900-01-01 00:00:00.000000 +00:00';"
	returnedSQL = sqlBuilder.AlterTableAddColumn(userTableName, createdAtColumn)
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.AlterTableAddColumn")

	activatedAtColumn := testutils.FakeActivatedAtColumn(t)
	expectedSQL = "ALTER TABLE \"user\" ADD COLUMN \"activated_at\" <<TYPE_3>> DEFAULT '1900-01-01 00:00:00.000000';"
	returnedSQL = sqlBuilder.AlterTableAddColumn(userTableName, activatedAtColumn)
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.AlterTableAddColumn")

	priceColumn := testutils.FakePriceColumn(t)
	expectedSQL = "ALTER TABLE \"user\" ADD COLUMN \"price\" <<TYPE_2>> DEFAULT 0.0;"
	returnedSQL = sqlBuilder.AlterTableAddColumn(userTableName, priceColumn)
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.AlterTableAddColumn")

	scoreColumn := testutils.FakeScoreColumn(t)
	expectedSQL = "ALTER TABLE \"user\" ADD COLUMN \"score\" <<TYPE_0>> DEFAULT 0;"
	returnedSQL = sqlBuilder.AlterTableAddColumn(userTableName, scoreColumn)
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.AlterTableAddColumn")

	nameColumn := testutils.FakeNameColumn(t)
	expectedSQL = "ALTER TABLE \"user\" ADD COLUMN \"name\" <<TYPE_1>> DEFAULT '';"
	returnedSQL = sqlBuilder.AlterTableAddColumn(userTableName, nameColumn)
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.AlterTableAddColumn")
}

func TestAlterTableDropColumn(t *testing.T) {
	sqlBuilder := internal.NewGenericSQLBuilder(nil, nil, nil)

	userTableName := testutils.FakeUserTableName(t)
	updatedAtColumnName := testutils.FakeUpdatedAtColumnName(t)
	expectedSQL := "ALTER TABLE \"user\" DROP COLUMN \"updated_at\";"
	returnedSQL := sqlBuilder.AlterTableDropColumn(userTableName, updatedAtColumnName)
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.AlterTableDropColumn")
}

func TestInsert(t *testing.T) {
	sqlBuilder := internal.NewGenericSQLBuilder(nil, nil, nil)

	returnedSQL := sqlBuilder.Insert("product", []string{"id", "sku", "name", "description"})
	expectedSQL := "insert into \"product\" (\"id\",\"sku\",\"name\",\"description\") values (@id,@sku,@name,@description)"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.Insert")
}

func TestInsertWithValues(t *testing.T) {
	sqlBuilder := internal.NewGenericSQLBuilder(nil, nil, nil)

	returnedSQL, returnedValues := sqlBuilder.InsertWithValues("product", ormshift.ColumnsValues{"id": 1, "sku": "1.005.12.9", "name": "Trufa Sabor Amarula 30g Cacaushow"})
	expectedSQL := "insert into \"product\" (\"id\",\"name\",\"sku\") values (@id,@name,@sku)"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.InsertWithValues.SQL")
	testutils.AssertEqualWithLabel(t, 3, len(returnedValues), "SQLBuilder.InsertWithValues.Values")
	testutils.AssertNamedArgEqualWithLabel(t, returnedValues[0], sql.NamedArg{Name: "id", Value: 1}, "SQLBuilder.InsertWithValues.Values[0]")
	testutils.AssertNamedArgEqualWithLabel(t, returnedValues[1], sql.NamedArg{Name: "name", Value: "Trufa Sabor Amarula 30g Cacaushow"}, "SQLBuilder.InsertWithValues.Values[1]")
	testutils.AssertNamedArgEqualWithLabel(t, returnedValues[2], sql.NamedArg{Name: "sku", Value: "1.005.12.9"}, "SQLBuilder.InsertWithValues.Values[2]")
}

func TestUpdate(t *testing.T) {
	sqlBuilder := internal.NewGenericSQLBuilder(nil, nil, nil)

	returnedSQL := sqlBuilder.Update("product", []string{"sku", "name", "description"}, []string{"id"})
	expectedSQL := "update \"product\" set \"sku\" = @sku,\"name\" = @name,\"description\" = @description where \"id\" = @id"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.Update")
}

func TestUpdateWithValues(t *testing.T) {
	sqlBuilder := internal.NewGenericSQLBuilder(nil, nil, nil)

	returnedSQL, returnedValues := sqlBuilder.UpdateWithValues("product", []string{"sku", "name"}, []string{"id"}, ormshift.ColumnsValues{"id": 1, "sku": "1.005.12.5", "name": "Trufa Sabor Amarula 18g Cacaushow"})
	expectedSQL := "update \"product\" set \"sku\" = @sku,\"name\" = @name where \"id\" = @id"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.UpdateWithValues.SQL")
	testutils.AssertEqualWithLabel(t, 3, len(returnedValues), "SQLBuilder.UpdateWithValues.Values")
	testutils.AssertNamedArgEqualWithLabel(t, returnedValues[0], sql.NamedArg{Name: "id", Value: 1}, "SQLBuilder.UpdateWithValues.Values[0]")
	testutils.AssertNamedArgEqualWithLabel(t, returnedValues[1], sql.NamedArg{Name: "name", Value: "Trufa Sabor Amarula 18g Cacaushow"}, "SQLBuilder.UpdateWithValues.Values[1]")
	testutils.AssertNamedArgEqualWithLabel(t, returnedValues[2], sql.NamedArg{Name: "sku", Value: "1.005.12.5"}, "SQLBuilder.UpdateWithValues.Values[2]")
}

func TestDelete(t *testing.T) {
	sqlBuilder := internal.NewGenericSQLBuilder(nil, nil, nil)

	returnedSQL := sqlBuilder.Delete("product", []string{"id"})
	expectedSQL := "delete from \"product\" where \"id\" = @id"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.Delete")
}

func TestDeleteWithValues(t *testing.T) {
	sqlBuilder := internal.NewGenericSQLBuilder(nil, nil, nil)

	returnedSQL, returnedValues := sqlBuilder.DeleteWithValues("product", ormshift.ColumnsValues{"id": 1})
	expectedSQL := "delete from \"product\" where \"id\" = @id"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.DeleteWithValues.SQL")
	testutils.AssertEqualWithLabel(t, 1, len(returnedValues), "SQLBuilder.DeleteWithValues.Values")
	testutils.AssertNamedArgEqualWithLabel(t, returnedValues[0], sql.NamedArg{Name: "id", Value: 1}, "SQLBuilder.DeleteWithValues.Values[0]")
}

func TestSelect(t *testing.T) {
	sqlBuilder := internal.NewGenericSQLBuilder(nil, nil, nil)

	returnedSQL := sqlBuilder.Select("product", []string{"id", "name", "description"}, []string{"sku", "active"})
	expectedSQL := "select \"id\",\"name\",\"description\" from \"product\" where \"sku\" = @sku and \"active\" = @active"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.Select")
}

func TestSelectWithValues(t *testing.T) {
	sqlBuilder := internal.NewGenericSQLBuilder(nil, nil, nil)

	returnedSQL, returnedValues := sqlBuilder.SelectWithValues("product", []string{"id", "sku", "name", "description"}, ormshift.ColumnsValues{"category_id": 1, "active": true})
	expectedSQL := "select \"id\",\"sku\",\"name\",\"description\" from \"product\" where \"active\" = @active and \"category_id\" = @category_id"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.SelectWithValues.SQL")
	testutils.AssertEqualWithLabel(t, 2, len(returnedValues), "SQLBuilder.SelectWithValues.Values")
	testutils.AssertNamedArgEqualWithLabel(t, returnedValues[0], sql.NamedArg{Name: "active", Value: true}, "SQLBuilder.SelectWithValues.Values[0]")
	testutils.AssertNamedArgEqualWithLabel(t, returnedValues[1], sql.NamedArg{Name: "category_id", Value: 1}, "SQLBuilder.SelectWithValues.Values[1]")
}

func TestSelectWithPagination(t *testing.T) {
	sqlBuilder := internal.NewGenericSQLBuilder(nil, nil, nil)

	returnedSQL := sqlBuilder.SelectWithPagination("select * from product", 10, 5)
	expectedSQL := "select * from product LIMIT 10 OFFSET 40"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.SelectWithPagination")
}

func TestInteroperateSQLCommandWithNamedArgs(t *testing.T) {
	sqlBuilder := internal.NewGenericSQLBuilder(nil, nil, testutils.FakeInteroperateSQLCommandWithNamedArgsFunc)
	returnedSQL, returnedNamedArgs := sqlBuilder.InteroperateSQLCommandWithNamedArgs("original command", sql.NamedArg{Name: "param1", Value: 1})
	expectedSQL := "command has been modified"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.InteroperateSQLCommandWithNamedArgs.SQL")
	testutils.AssertEqualWithLabel(t, 1, len(returnedNamedArgs), "SQLBuilder.InteroperateSQLCommandWithNamedArgs.NamedArgs")
	testutils.AssertNamedArgEqualWithLabel(t, returnedNamedArgs[0], sql.NamedArg{Name: "param1", Value: 1}, "SQLBuilder.InteroperateSQLCommandWithNamedArgs.NamedArgs[0]")
}

func TestColumnDefinition(t *testing.T) {
	sqlBuilder := internal.NewGenericSQLBuilder(testutils.FakeColumnDefinitionFunc, nil, nil)
	column := schema.NewColumn(schema.NewColumnParams{Name: "column_name", Type: schema.Integer, Size: 0})
	table := "test_table"
	returnedSQL := sqlBuilder.AlterTableAddColumn(table, column)
	testutils.AssertEqualWithLabel(t, "ALTER TABLE \"test_table\" ADD COLUMN fake;", returnedSQL, "SQLBuilder.ColumnDefinition")
}

func TestQuoteIdentifier(t *testing.T) {
	sqlBuilder := internal.NewGenericSQLBuilder(nil, testutils.FakeQuoteIdentifierFunc, nil)
	column := schema.NewColumn(schema.NewColumnParams{Name: "column_name", Type: schema.Integer, Size: 0})
	table := "test_table"
	returnedSQL := sqlBuilder.AlterTableAddColumn(table, column)
	testutils.AssertEqualWithLabel(t, "ALTER TABLE quoted_test_table ADD COLUMN quoted_column_name <<TYPE_0>>;", returnedSQL, "SQLBuilder.QuoteIdentifier")
}
