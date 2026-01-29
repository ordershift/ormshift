package postgresql_test

import (
	"database/sql"
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/postgresql"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestInteroperateSQLCommandWithNamedArgs(t *testing.T) {
	driver := postgresql.Driver()
	returnedSQL, returnedValue := driver.SQLBuilder().InteroperateSQLCommandWithNamedArgs(
		"select * from user where role = @role and active = @active and master = @master",
		sql.Named("role", "admin"),
		sql.Named("active", true),
		sql.Named("master", false),
	)
	expectedSQL := "select * from user where role = $1 and active = $2 and master = $3"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "DriverPostgresql.InteroperateSQLCommandWithNamedArgs.SQL")
	testutils.AssertEqualWithLabel(t, "admin", returnedValue[0].(string), "DriverPostgresql.InteroperateSQLCommandWithNamedArgs.Value.1")
	testutils.AssertEqualWithLabel(t, 1, returnedValue[1].(int), "DriverPostgresql.InteroperateSQLCommandWithNamedArgs.Value.2")
	testutils.AssertEqualWithLabel(t, 0, returnedValue[2].(int), "DriverPostgresql.InteroperateSQLCommandWithNamedArgs.Value.3")
	returnedSQL, returnedValue = driver.SQLBuilder().InteroperateSQLCommandWithNamedArgs(
		"update user set role = @role where role = @role",
		sql.Named("role", "admin"),
	)
	expectedSQL = "update user set role = $1 where role = $1"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "DriverPostgresql.InteroperateSQLCommandWithNamedArgs.SQL")
	testutils.AssertEqualWithLabel(t, "admin", returnedValue[0].(string), "DriverPostgresql.InteroperateSQLCommandWithNamedArgs.Value.1")
	returnedSQL, returnedValue = driver.SQLBuilder().InteroperateSQLCommandWithNamedArgs(
		"delete from user where id = @id",
		sql.Named("role", "admin"),
	)
	expectedSQL = "delete from user where id = @id"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "DriverPostgresql.InteroperateSQLCommandWithNamedArgs.SQL")
	testutils.AssertEqualWithLabel(t, "admin", returnedValue[0].(string), "DriverPostgresql.InteroperateSQLCommandWithNamedArgs.Value.1")
}

func TestCreateTable(t *testing.T) {
	sqlBuilder := postgresql.Driver().SQLBuilder()

	userTable := testutils.FakeUserTable(t)
	expectedSQL := "CREATE TABLE \"user\" (\"id\" BIGSERIAL NOT NULL,\"email\" VARCHAR(80) NOT NULL,\"name\" VARCHAR(50) NOT NULL," +
		"\"password_hash\" VARCHAR(256),\"active\" SMALLINT,\"created_at\" TIMESTAMP(6),\"user_master\" BIGINT,\"master_user_id\" BIGINT," +
		"\"licence_price\" NUMERIC(17,2),\"relevance\" DOUBLE PRECISION,\"photo\" BYTEA,\"any\" VARCHAR,PRIMARY KEY (\"id\",\"email\"));"
	returnedSQL := sqlBuilder.CreateTable(userTable)
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.CreateTable")

	productAttributeTable := testutils.FakeProductAttributeTable(t)
	expectedSQL = "CREATE TABLE \"product_attribute\" (\"product_id\" BIGINT NOT NULL,\"attribute_id\" BIGINT NOT NULL,\"value\" VARCHAR(75),\"position\" BIGINT,PRIMARY KEY (\"product_id\",\"attribute_id\"));"
	returnedSQL = sqlBuilder.CreateTable(productAttributeTable)
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.CreateTable")
}

func TestDropTable(t *testing.T) {
	sqlBuilder := postgresql.Driver().SQLBuilder()

	userTableName := testutils.FakeUserTableName(t)
	expectedSQL := "DROP TABLE \"user\";"
	returnedSQL := sqlBuilder.DropTable(userTableName)
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.DropTable")
}

func TestAlterTableAddColumn(t *testing.T) {
	sqlBuilder := postgresql.Driver().SQLBuilder()

	userTableName := testutils.FakeUserTableName(t)
	updatedAtColumn := testutils.FakeUpdatedAtColumn(t)
	expectedSQL := "ALTER TABLE \"user\" ADD COLUMN \"updated_at\" TIMESTAMP(6);"
	returnedSQL := sqlBuilder.AlterTableAddColumn(userTableName, updatedAtColumn)
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.AlterTableAddColumn")
}

func TestAlterTableDropColumn(t *testing.T) {
	sqlBuilder := postgresql.Driver().SQLBuilder()

	userTableName := testutils.FakeUserTableName(t)
	updatedAtColumnName := testutils.FakeUpdatedAtColumnName(t)
	expectedSQL := "ALTER TABLE \"user\" DROP COLUMN \"updated_at\";"
	returnedSQL := sqlBuilder.AlterTableDropColumn(userTableName, updatedAtColumnName)
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.AlterTableDropColumn")
}

func TestInsert(t *testing.T) {
	sqlBuilder := postgresql.Driver().SQLBuilder()

	returnedSQL := sqlBuilder.Insert("product", []string{"id", "sku", "name", "description"})
	expectedSQL := "insert into \"product\" (\"id\",\"sku\",\"name\",\"description\") values (@id,@sku,@name,@description)"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.Insert")
}

func TestInsertWithValues(t *testing.T) {
	sqlBuilder := postgresql.Driver().SQLBuilder()

	returnedSQL, returnedValues := sqlBuilder.InsertWithValues("product", ormshift.ColumnsValues{"id": 1, "sku": "1.005.12.9", "name": "Trufa Sabor Amarula 30g Cacaushow"})
	expectedSQL := "insert into \"product\" (\"id\",\"name\",\"sku\") values ($1,$2,$3)"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.InsertWithValues.SQL")
	testutils.AssertEqualWithLabel(t, 3, len(returnedValues), "SQLBuilder.InsertWithValues.Values")
	testutils.AssertEqualWithLabel(t, 1, returnedValues[0], "SQLBuilder.InsertWithValues.Values[0]")
	testutils.AssertEqualWithLabel(t, "Trufa Sabor Amarula 30g Cacaushow", returnedValues[1], "SQLBuilder.InsertWithValues.Values[1]")
	testutils.AssertEqualWithLabel(t, "1.005.12.9", returnedValues[2], "SQLBuilder.InsertWithValues.Values[2]")
}

func TestUpdate(t *testing.T) {
	sqlBuilder := postgresql.Driver().SQLBuilder()

	returnedSQL := sqlBuilder.Update("product", []string{"sku", "name", "description"}, []string{"id"})
	expectedSQL := "update \"product\" set \"sku\" = @sku,\"name\" = @name,\"description\" = @description where \"id\" = @id"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.Update")
}

func TestUpdateWithValues(t *testing.T) {
	sqlBuilder := postgresql.Driver().SQLBuilder()

	returnedSQL, returnedValues := sqlBuilder.UpdateWithValues("product", []string{"sku", "name"}, []string{"id"}, ormshift.ColumnsValues{"id": 1, "sku": "1.005.12.5", "name": "Trufa Sabor Amarula 18g Cacaushow"})
	expectedSQL := "update \"product\" set \"sku\" = $3,\"name\" = $2 where \"id\" = $1"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.UpdateWithValues.SQL")
	testutils.AssertEqualWithLabel(t, 3, len(returnedValues), "SQLBuilder.UpdateWithValues.Values")
	testutils.AssertEqualWithLabel(t, 1, returnedValues[0], "SQLBuilder.UpdateWithValues.Values[0]")
	testutils.AssertEqualWithLabel(t, "Trufa Sabor Amarula 18g Cacaushow", returnedValues[1], "SQLBuilder.UpdateWithValues.Values[1]")
	testutils.AssertEqualWithLabel(t, "1.005.12.5", returnedValues[2], "SQLBuilder.UpdateWithValues.Values[2]")
}

func TestDelete(t *testing.T) {
	sqlBuilder := postgresql.Driver().SQLBuilder()

	returnedSQL := sqlBuilder.Delete("product", []string{"id"})
	expectedSQL := "delete from \"product\" where \"id\" = @id"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.Delete")
}

func TestDeleteWithValues(t *testing.T) {
	sqlBuilder := postgresql.Driver().SQLBuilder()

	returnedSQL, returnedValues := sqlBuilder.DeleteWithValues("product", ormshift.ColumnsValues{"id": 1})
	expectedSQL := "delete from \"product\" where \"id\" = $1"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.DeleteWithValues.SQL")
	testutils.AssertEqualWithLabel(t, 1, len(returnedValues), "SQLBuilder.DeleteWithValues.Values")
	testutils.AssertEqualWithLabel(t, 1, returnedValues[0], "SQLBuilder.DeleteWithValues.Values[0]")
}

func TestSelect(t *testing.T) {
	sqlBuilder := postgresql.Driver().SQLBuilder()

	returnedSQL := sqlBuilder.Select("product", []string{"id", "name", "description"}, []string{"sku", "active"})
	expectedSQL := "select \"id\",\"name\",\"description\" from \"product\" where \"sku\" = @sku and \"active\" = @active"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.Select")
}

func TestSelectWithValues(t *testing.T) {
	sqlBuilder := postgresql.Driver().SQLBuilder()

	returnedSQL, returnedValues := sqlBuilder.SelectWithValues("product", []string{"id", "sku", "name", "description"}, ormshift.ColumnsValues{"category_id": 1, "active": true})
	expectedSQL := "select \"id\",\"sku\",\"name\",\"description\" from \"product\" where \"active\" = $1 and \"category_id\" = $2"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.SelectWithValues.SQL")
	testutils.AssertEqualWithLabel(t, 2, len(returnedValues), "SQLBuilder.SelectWithValues.Values")
	testutils.AssertEqualWithLabel(t, 1, returnedValues[0], "SQLBuilder.SelectWithValues.Values[0]")
	testutils.AssertEqualWithLabel(t, 1, returnedValues[1], "SQLBuilder.SelectWithValues.Values[1]")
}

func TestSelectWithPagination(t *testing.T) {
	sqlBuilder := postgresql.Driver().SQLBuilder()

	returnedSQL := sqlBuilder.SelectWithPagination("select * from product", 10, 5)
	expectedSQL := "select * from product LIMIT 10 OFFSET 40"
	testutils.AssertEqualWithLabel(t, expectedSQL, returnedSQL, "SQLBuilder.SelectWithPagination")
}
