package postgresql_test

import (
	"database/sql"
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/postgresql"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestInteroperateSQLCommandWithNamedArgs(t *testing.T) {
	lDriver := postgresql.Driver()
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
	lSQLBuilder := postgresql.Driver().SQLBuilder()

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

func TestDropTable(t *testing.T) {
	lSQLBuilder := postgresql.Driver().SQLBuilder()

	lUserTableName := testutils.FakeUserTableName(t)
	lExpectedSQL := "DROP TABLE user;"
	lReturnedSQL := lSQLBuilder.DropTable(*lUserTableName)
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.DropTable")
}

func TestAlterTableAddColumn(t *testing.T) {
	lSQLBuilder := postgresql.Driver().SQLBuilder()

	lUserTableName := testutils.FakeUserTableName(t)
	lUpdatedAtColumn := testutils.FakeUpdatedAtColumn(t)
	lExpectedSQL := "ALTER TABLE user ADD COLUMN updated_at TIMESTAMP(6);"
	lReturnedSQL := lSQLBuilder.AlterTableAddColumn(*lUserTableName, *lUpdatedAtColumn)
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.AlterTableAddColumn")
}

func TestAlterTableDropColumn(t *testing.T) {
	lSQLBuilder := postgresql.Driver().SQLBuilder()

	lUserTableName := testutils.FakeUserTableName(t)
	lUpdatedAtColumnName := testutils.FakeUpdatedAtColumnName(t)
	lExpectedSQL := "ALTER TABLE user DROP COLUMN updated_at;"
	lReturnedSQL := lSQLBuilder.AlterTableDropColumn(*lUserTableName, *lUpdatedAtColumnName)
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.AlterTableDropColumn")
}

func TestInsert(t *testing.T) {
	lSQLBuilder := postgresql.Driver().SQLBuilder()

	lReturnedSQL := lSQLBuilder.Insert("product", []string{"id", "sku", "name", "description"})
	lExpectedSQL := "insert into product (id,sku,name,description) values (@id,@sku,@name,@description)"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.Insert")
}

func TestInsertWithValues(t *testing.T) {
	lSQLBuilder := postgresql.Driver().SQLBuilder()

	lReturnedSQL, lReturnedValues := lSQLBuilder.InsertWithValues("product", ormshift.ColumnsValues{"id": 1, "sku": "1.005.12.9", "name": "Trufa Sabor Amarula 30g Cacaushow"})
	lExpectedSQL := "insert into product (id,name,sku) values ($1,$2,$3)"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.InsertWithValues.SQL")
	testutils.AssertEqualWithLabel(t, 3, len(lReturnedValues), "SQLBuilder.InsertWithValues.Values")
	testutils.AssertEqualWithLabel(t, 1, lReturnedValues[0], "SQLBuilder.InsertWithValues.Values[0]")
	testutils.AssertEqualWithLabel(t, "Trufa Sabor Amarula 30g Cacaushow", lReturnedValues[1], "SQLBuilder.InsertWithValues.Values[1]")
	testutils.AssertEqualWithLabel(t, "1.005.12.9", lReturnedValues[2], "SQLBuilder.InsertWithValues.Values[2]")
}

func TestUpdate(t *testing.T) {
	lSQLBuilder := postgresql.Driver().SQLBuilder()

	lReturnedSQL := lSQLBuilder.Update("product", []string{"sku", "name", "description"}, []string{"id"})
	lExpectedSQL := "update product set sku = @sku,name = @name,description = @description where id = @id"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.Update")
}

func TestUpdateWithValues(t *testing.T) {
	lSQLBuilder := postgresql.Driver().SQLBuilder()

	lReturnedSQL, lReturnedValues := lSQLBuilder.UpdateWithValues("product", []string{"sku", "name"}, []string{"id"}, ormshift.ColumnsValues{"id": 1, "sku": "1.005.12.5", "name": "Trufa Sabor Amarula 18g Cacaushow"})
	lExpectedSQL := "update product set sku = $3,name = $2 where id = $1"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.UpdateWithValues.SQL")
	testutils.AssertEqualWithLabel(t, 3, len(lReturnedValues), "SQLBuilder.UpdateWithValues.Values")
	testutils.AssertEqualWithLabel(t, 1, lReturnedValues[0], "SQLBuilder.UpdateWithValues.Values[0]")
	testutils.AssertEqualWithLabel(t, "Trufa Sabor Amarula 18g Cacaushow", lReturnedValues[1], "SQLBuilder.UpdateWithValues.Values[1]")
	testutils.AssertEqualWithLabel(t, "1.005.12.5", lReturnedValues[2], "SQLBuilder.UpdateWithValues.Values[2]")
}

func TestDelete(t *testing.T) {
	lSQLBuilder := postgresql.Driver().SQLBuilder()

	lReturnedSQL := lSQLBuilder.Delete("product", []string{"id"})
	lExpectedSQL := "delete from product where id = @id"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.Delete")
}

func TestDeleteWithValues(t *testing.T) {
	lSQLBuilder := postgresql.Driver().SQLBuilder()

	lReturnedSQL, lReturnedValues := lSQLBuilder.DeleteWithValues("product", ormshift.ColumnsValues{"id": 1})
	lExpectedSQL := "delete from product where id = $1"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.DeleteWithValues.SQL")
	testutils.AssertEqualWithLabel(t, 1, len(lReturnedValues), "SQLBuilder.DeleteWithValues.Values")
	testutils.AssertEqualWithLabel(t, 1, lReturnedValues[0], "SQLBuilder.DeleteWithValues.Values[0]")
}

func TestSelect(t *testing.T) {
	lSQLBuilder := postgresql.Driver().SQLBuilder()

	lReturnedSQL := lSQLBuilder.Select("product", []string{"id", "name", "description"}, []string{"sku", "active"})
	lExpectedSQL := "select id,name,description from product where sku = @sku and active = @active"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.Select")
}

func TestSelectWithValues(t *testing.T) {
	lSQLBuilder := postgresql.Driver().SQLBuilder()

	lReturnedSQL, lReturnedValues := lSQLBuilder.SelectWithValues("product", []string{"id", "sku", "name", "description"}, ormshift.ColumnsValues{"category_id": 1, "active": true})
	lExpectedSQL := "select id,sku,name,description from product where active = $1 and category_id = $2"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.SelectWithValues.SQL")
	testutils.AssertEqualWithLabel(t, 2, len(lReturnedValues), "SQLBuilder.SelectWithValues.Values")
	testutils.AssertEqualWithLabel(t, 1, lReturnedValues[0], "SQLBuilder.SelectWithValues.Values[0]")
	testutils.AssertEqualWithLabel(t, 1, lReturnedValues[1], "SQLBuilder.SelectWithValues.Values[1]")
}

func TestSelectWithPagination(t *testing.T) {
	lSQLBuilder := postgresql.Driver().SQLBuilder()

	lReturnedSQL := lSQLBuilder.SelectWithPagination("select * from product", 10, 5)
	lExpectedSQL := "select * from product LIMIT 10 OFFSET 40"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "SQLBuilder.SelectWithPagination")
}
