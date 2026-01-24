package internal_test

import (
	"database/sql"
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/internal"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestCreateTable(t *testing.T) {
	lSQLBuilder := internal.NewGenericSQLBuilder(nil, nil)

	lUserTable := testutils.FakeUserTable(t)
	lExpectedSQL := "CREATE TABLE user (id <<TYPE_0>>,email <<TYPE_1>>,name <<TYPE_1>>,password_hash <<TYPE_1>>," +
		"active <<TYPE_5>>,created_at <<TYPE_3>>,user_master <<TYPE_0>>,master_user_id <<TYPE_0>>," +
		"licence_price <<TYPE_2>>,relevance <<TYPE_4>>,photo <<TYPE_6>>,any <<TYPE_-1>>,PRIMARY KEY (id,email));"
	lReturnedSQL := lSQLBuilder.CreateTable(*lUserTable)
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.CreateTable")

	lProductAttributeTable := testutils.FakeProductAttributeTable(t)
	lExpectedSQL = "CREATE TABLE product_attribute (product_id <<TYPE_0>>,attribute_id <<TYPE_0>>,value <<TYPE_1>>,position <<TYPE_0>>,PRIMARY KEY (product_id,attribute_id));"
	lReturnedSQL = lSQLBuilder.CreateTable(*lProductAttributeTable)
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.CreateTable")
}

func TestDropTable(t *testing.T) {
	lSQLBuilder := internal.NewGenericSQLBuilder(nil, nil)

	lUserTableName := testutils.FakeUserTableName(t)
	lExpectedSQL := "DROP TABLE user;"
	lReturnedSQL := lSQLBuilder.DropTable(*lUserTableName)
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.DropTable")
}

func TestAlterTableAddColumn(t *testing.T) {
	lSQLBuilder := internal.NewGenericSQLBuilder(nil, nil)

	lUserTableName := testutils.FakeUserTableName(t)
	lUpdatedAtColumn := testutils.FakeUpdatedAtColumn(t)
	lExpectedSQL := "ALTER TABLE user ADD COLUMN updated_at <<TYPE_3>>;"
	lReturnedSQL := lSQLBuilder.AlterTableAddColumn(*lUserTableName, *lUpdatedAtColumn)
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.AlterTableAddColumn")
}

func TestAlterTableDropColumn(t *testing.T) {
	lSQLBuilder := internal.NewGenericSQLBuilder(nil, nil)

	lUserTableName := testutils.FakeUserTableName(t)
	lUpdatedAtColumnName := testutils.FakeUpdatedAtColumnName(t)
	lExpectedSQL := "ALTER TABLE user DROP COLUMN updated_at;"
	lReturnedSQL := lSQLBuilder.AlterTableDropColumn(*lUserTableName, *lUpdatedAtColumnName)
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.AlterTableDropColumn")
}

func TestInsert(t *testing.T) {
	lSQLBuilder := internal.NewGenericSQLBuilder(nil, nil)

	lReturnedSQL := lSQLBuilder.Insert("product", []string{"id", "sku", "name", "description"})
	lExpectedSQL := "insert into product (id,sku,name,description) values (@id,@sku,@name,@description)"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.Insert")
}

func TestInsertWithValues(t *testing.T) {
	lSQLBuilder := internal.NewGenericSQLBuilder(nil, nil)

	lReturnedSQL, lReturnedValues := lSQLBuilder.InsertWithValues("product", ormshift.ColumnsValues{"id": 1, "sku": "1.005.12.9", "name": "Trufa Sabor Amarula 30g Cacaushow"})
	lExpectedSQL := "insert into product (id,name,sku) values (@id,@name,@sku)"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.InsertWithValues.SQL")
	testutils.AssertEqualWithLabel(t, 3, len(lReturnedValues), "SQLBuilder.InsertWithValues.Values")
	testutils.AssertNamedArgEqualWithLabel(t, lReturnedValues[0], sql.NamedArg{Name: "id", Value: 1}, "SQLBuilder.InsertWithValues.Values[0]")
	testutils.AssertNamedArgEqualWithLabel(t, lReturnedValues[1], sql.NamedArg{Name: "name", Value: "Trufa Sabor Amarula 30g Cacaushow"}, "SQLBuilder.InsertWithValues.Values[1]")
	testutils.AssertNamedArgEqualWithLabel(t, lReturnedValues[2], sql.NamedArg{Name: "sku", Value: "1.005.12.9"}, "SQLBuilder.InsertWithValues.Values[2]")
}

func TestUpdate(t *testing.T) {
	lSQLBuilder := internal.NewGenericSQLBuilder(nil, nil)

	lReturnedSQL := lSQLBuilder.Update("product", []string{"sku", "name", "description"}, []string{"id"})
	lExpectedSQL := "update product set sku = @sku,name = @name,description = @description where id = @id"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.Update")
}

func TestUpdateWithValues(t *testing.T) {
	lSQLBuilder := internal.NewGenericSQLBuilder(nil, nil)

	lReturnedSQL, lReturnedValues := lSQLBuilder.UpdateWithValues("product", []string{"sku", "name"}, []string{"id"}, ormshift.ColumnsValues{"id": 1, "sku": "1.005.12.5", "name": "Trufa Sabor Amarula 18g Cacaushow"})
	lExpectedSQL := "update product set sku = @sku,name = @name where id = @id"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.UpdateWithValues.SQL")
	testutils.AssertEqualWithLabel(t, 3, len(lReturnedValues), "SQLBuilder.UpdateWithValues.Values")
	testutils.AssertNamedArgEqualWithLabel(t, lReturnedValues[0], sql.NamedArg{Name: "id", Value: 1}, "SQLBuilder.UpdateWithValues.Values[0]")
	testutils.AssertNamedArgEqualWithLabel(t, lReturnedValues[1], sql.NamedArg{Name: "name", Value: "Trufa Sabor Amarula 18g Cacaushow"}, "SQLBuilder.UpdateWithValues.Values[1]")
	testutils.AssertNamedArgEqualWithLabel(t, lReturnedValues[2], sql.NamedArg{Name: "sku", Value: "1.005.12.5"}, "SQLBuilder.UpdateWithValues.Values[2]")
}

func TestDelete(t *testing.T) {
	lSQLBuilder := internal.NewGenericSQLBuilder(nil, nil)

	lReturnedSQL := lSQLBuilder.Delete("product", []string{"id"})
	lExpectedSQL := "delete from product where id = @id"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.Delete")
}

func TestDeleteWithValues(t *testing.T) {
	lSQLBuilder := internal.NewGenericSQLBuilder(nil, nil)

	lReturnedSQL, lReturnedValues := lSQLBuilder.DeleteWithValues("product", ormshift.ColumnsValues{"id": 1})
	lExpectedSQL := "delete from product where id = @id"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.DeleteWithValues.SQL")
	testutils.AssertEqualWithLabel(t, 1, len(lReturnedValues), "SQLBuilder.DeleteWithValues.Values")
	testutils.AssertNamedArgEqualWithLabel(t, lReturnedValues[0], sql.NamedArg{Name: "id", Value: 1}, "SQLBuilder.DeleteWithValues.Values[0]")
}

func TestSelect(t *testing.T) {
	lSQLBuilder := internal.NewGenericSQLBuilder(nil, nil)

	lReturnedSQL := lSQLBuilder.Select("product", []string{"id", "name", "description"}, []string{"sku", "active"})
	lExpectedSQL := "select id,name,description from product where sku = @sku and active = @active"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.Select")
}

func TestSelectWithValues(t *testing.T) {
	lSQLBuilder := internal.NewGenericSQLBuilder(nil, nil)

	lReturnedSQL, lReturnedValues := lSQLBuilder.SelectWithValues("product", []string{"id", "sku", "name", "description"}, ormshift.ColumnsValues{"category_id": 1, "active": true})
	lExpectedSQL := "select id,sku,name,description from product where active = @active and category_id = @category_id"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.SelectWithValues.SQL")
	testutils.AssertEqualWithLabel(t, 2, len(lReturnedValues), "SQLBuilder.SelectWithValues.Values")
	testutils.AssertNamedArgEqualWithLabel(t, lReturnedValues[0], sql.NamedArg{Name: "active", Value: true}, "SQLBuilder.SelectWithValues.Values[0]")
	testutils.AssertNamedArgEqualWithLabel(t, lReturnedValues[1], sql.NamedArg{Name: "category_id", Value: 1}, "SQLBuilder.SelectWithValues.Values[1]")
}

func TestSelectWithPagination(t *testing.T) {
	lSQLBuilder := internal.NewGenericSQLBuilder(nil, nil)

	lReturnedSQL := lSQLBuilder.SelectWithPagination("select * from product", 10, 5)
	lExpectedSQL := "select * from product LIMIT 10 OFFSET 40"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.SelectWithPagination")
}

func TestInteroperateSQLCommandWithNamedArgs(t *testing.T) {
	lSQLBuilder := internal.NewGenericSQLBuilder(nil, testutils.FakeInteroperateSQLCommandWithNamedArgsFunc)
	lReturnedSQL, lReturnedNamedArgs := lSQLBuilder.InteroperateSQLCommandWithNamedArgs("original command", sql.NamedArg{Name: "param1", Value: 1})
	lExpectedSQL := "command has been modified"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.InteroperateSQLCommandWithNamedArgs.SQL")
	testutils.AssertEqualWithLabel(t, 1, len(lReturnedNamedArgs), "SQLBuilder.InteroperateSQLCommandWithNamedArgs.NamedArgs")
	testutils.AssertNamedArgEqualWithLabel(t, lReturnedNamedArgs[0], sql.NamedArg{Name: "param1", Value: 1}, "SQLBuilder.InteroperateSQLCommandWithNamedArgs.NamedArgs[0]")
}
