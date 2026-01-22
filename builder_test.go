package ormshift_test

import (
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/internal/testutils"
)

func Test_ColumnsValues_ToNamedArgs_ShouldReturnDefault(t *testing.T) {
	lColumnsValues := ormshift.ColumnsValues{"id": 1, "sku": "ABC1234", "active": true}
	lNamedArgs := lColumnsValues.ToNamedArgs()
	testutils.AssertEqualWithLabel(t, 3, len(lNamedArgs), "ColumnsValues.ToNamedArgs")
	testutils.AssertEqualWithLabel(t, lNamedArgs[0].Name, "active", "ColumnsValues.ToNamedArgs[0].Name")
	testutils.AssertEqualWithLabel(t, lNamedArgs[0].Value, true, "ColumnsValues.ToNamedArgs[0].Value")
	testutils.AssertEqualWithLabel(t, lNamedArgs[1].Name, "id", "ColumnsValues.ToNamedArgs[1].Name")
	testutils.AssertEqualWithLabel(t, lNamedArgs[1].Value, 1, "ColumnsValues.ToNamedArgs[1].Value")
	testutils.AssertEqualWithLabel(t, lNamedArgs[2].Name, "sku", "ColumnsValues.ToNamedArgs[2].Name")
	testutils.AssertEqualWithLabel(t, lNamedArgs[2].Value, "ABC1234", "ColumnsValues.ToNamedArgs[2].Value")
}

func Test_ColumnsValues_ToColumns_ShouldReturnDefault(t *testing.T) {
	lColumnsValues := ormshift.ColumnsValues{"id": 1, "sku": "ABC1234"}
	lColumns := lColumnsValues.ToColumns()
	testutils.AssertEqualWithLabel(t, 2, len(lColumns), "ColumnsValues.ToColumns")
	testutils.AssertEqualWithLabel(t, lColumns[0], "id", "ColumnsValues.ToColumns[0]")
	testutils.AssertEqualWithLabel(t, lColumns[1], "sku", "ColumnsValues.ToColumns[1]")
}

// func Test_DriverDB_SQLBuilder_InteroperateSQLCommandWithNamedArgs_ShouldReturnDefault(t *testing.T) {
// 	lDriversDB := []ormshift.DriverDB{&sqlserver.DriverDB{}, &sqlite.DriverDB{}}
// 	for _, lDriverDB := range lDriversDB {
// 		lReturnedSQL, lReturnedValue := lDriverDB.SQLBuilder().InteroperateSQLCommandWithNamedArgs("select * from table where id = @id", sql.Named("id", 1))
// 		lExpectedSQL := "select * from table where id = @id"
// 		testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, lDriverDB.Name()+".InteroperateSQLCommandWithNamedArgs.SQL")
// 		testutils.AssertEqualWithLabel(t, "id", lReturnedValue[0].(sql.NamedArg).Name, lDriverDB.Name()+".InteroperateSQLCommandWithNamedArgs.Name")
// 		testutils.AssertEqualWithLabel(t, 1, lReturnedValue[0].(sql.NamedArg).Value.(int), lDriverDB.Name()+".InteroperateSQLCommandWithNamedArgs.Value")
// 	}
// }

// func Test_DriverPostgresql_SQLBuilder_InteroperateSQLCommandWithNamedArgs_ShouldReturnValid(t *testing.T) {
// 	lReturnedSQL, lReturnedValue := postgresql.DriverPostgresql.SQLBuilder().InteroperateSQLCommandWithNamedArgs(
// 		"select * from user where role = @role and active = @active and master = @master",
// 		sql.Named("role", "admin"),
// 		sql.Named("active", true),
// 		sql.Named("master", false),
// 	)
// 	lExpectedSQL := "select * from user where role = $1 and active = $2 and master = $3"
// 	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "DriverPostgresql.InteroperateSQLCommandWithNamedArgs.SQL")
// 	testutils.AssertEqualWithLabel(t, "admin", lReturnedValue[0].(string), "DriverPostgresql.InteroperateSQLCommandWithNamedArgs.Value.1")
// 	testutils.AssertEqualWithLabel(t, 1, lReturnedValue[1].(int), "DriverPostgresql.InteroperateSQLCommandWithNamedArgs.Value.2")
// 	testutils.AssertEqualWithLabel(t, 0, lReturnedValue[2].(int), "DriverPostgresql.InteroperateSQLCommandWithNamedArgs.Value.3")
// 	lReturnedSQL, lReturnedValue = postgresql.DriverPostgresql.SQLBuilder().InteroperateSQLCommandWithNamedArgs(
// 		"update user set role = @role where role = @role",
// 		sql.Named("role", "admin"),
// 	)
// 	lExpectedSQL = "update user set role = $1 where role = $1"
// 	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "DriverPostgresql.InteroperateSQLCommandWithNamedArgs.SQL")
// 	testutils.AssertEqualWithLabel(t, "admin", lReturnedValue[0].(string), "DriverPostgresql.InteroperateSQLCommandWithNamedArgs.Value.1")
// 	lReturnedSQL, lReturnedValue = postgresql.DriverPostgresql.SQLBuilder().InteroperateSQLCommandWithNamedArgs(
// 		"delete from user where id = @id",
// 		sql.Named("role", "admin"),
// 	)
// 	lExpectedSQL = "delete from user where id = @id"
// 	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, "DriverPostgresql.InteroperateSQLCommandWithNamedArgs.SQL")
// 	testutils.AssertEqualWithLabel(t, "admin", lReturnedValue[0].(string), "DriverPostgresql.InteroperateSQLCommandWithNamedArgs.Value.1")
// }

// func Test_DriverDB_SQLBuilder_CreateTable_ShouldReturnCorrectSQL(t *testing.T) {
// 	lUserTable := testutils.FakeUserTable(t)
// 	if lUserTable == nil {
// 		return
// 	}
// 	lDriversDBWithExpectedSQL := map[ormshift.DriverDB]string{
// 		postgresql.DriverPostgresql: "CREATE TABLE user (id BIGSERIAL NOT NULL,email VARCHAR(80) NOT NULL,name VARCHAR(50) NOT NULL," +
// 			"password_hash VARCHAR(256),active SMALLINT,created_at TIMESTAMP(6),user_master BIGINT,master_user_id BIGINT," +
// 			"licence_price NUMERIC(17,2),relevance DOUBLE PRECISION,photo BYTEA,any VARCHAR,PRIMARY KEY (id,email));",
// 		sqlserver.DriverSQLServer: "CREATE TABLE user (id BIGINT NOT NULL IDENTITY (1, 1),email VARCHAR(80) NOT NULL,name VARCHAR(50) NOT NULL," +
// 			"password_hash VARCHAR(256),active BIT,created_at DATETIME2(6),user_master BIGINT,master_user_id BIGINT," +
// 			"licence_price MONEY,relevance FLOAT,photo VARBINARY(MAX),any VARCHAR,CONSTRAINT PK_user PRIMARY KEY (id,email));",
// 		sqlite.DriverSQLite: "CREATE TABLE user (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,email TEXT NOT NULL,name TEXT NOT NULL," +
// 			"password_hash TEXT,active INTEGER,created_at DATETIME,user_master INTEGER,master_user_id INTEGER,licence_price REAL,relevance REAL,photo BLOB,any TEXT);",
// 		ormshift.DriverDB(-1): "CREATE TABLE user (id <<TYPE_0>>,email <<TYPE_1>>,name <<TYPE_1>>,password_hash <<TYPE_1>>," +
// 			"active <<TYPE_5>>,created_at <<TYPE_3>>,user_master <<TYPE_0>>,master_user_id <<TYPE_0>>," +
// 			"licence_price <<TYPE_2>>,relevance <<TYPE_4>>,photo <<TYPE_6>>,any <<TYPE_-1>>,PRIMARY KEY (id,email));",
// 	}
// 	for lDriverDB, lExpectedSQL := range lDriversDBWithExpectedSQL {
// 		lReturnedSQL := lDriverDB.SQLBuilder().CreateTable(*lUserTable)
// 		testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, lDriverDB.Name()+".SQLBuilder.CreateTable")
// 	}

// 	lProductAttributeTable := testutils.FakeProductAttributeTable(t)
// 	if lProductAttributeTable == nil {
// 		return
// 	}
// 	lDriversDBWithExpectedSQL = map[ormshift.DriverDB]string{
// 		postgresql.DriverPostgresql: "CREATE TABLE product_attribute (product_id BIGINT NOT NULL,attribute_id BIGINT NOT NULL,value VARCHAR(75),position BIGINT,PRIMARY KEY (product_id,attribute_id));",
// 		sqlserver.DriverSQLServer:   "CREATE TABLE product_attribute (product_id BIGINT NOT NULL,attribute_id BIGINT NOT NULL,value VARCHAR(75),position BIGINT,CONSTRAINT PK_product_attribute PRIMARY KEY (product_id,attribute_id));",
// 		sqlite.DriverSQLite:         "CREATE TABLE product_attribute (product_id INTEGER NOT NULL,attribute_id INTEGER NOT NULL,value TEXT,position INTEGER,CONSTRAINT PK_product_attribute PRIMARY KEY (product_id,attribute_id));",
// 		ormshift.DriverDB(-1):           "CREATE TABLE product_attribute (product_id <<TYPE_0>>,attribute_id <<TYPE_0>>,value <<TYPE_1>>,position <<TYPE_0>>,PRIMARY KEY (product_id,attribute_id));",
// 	}
// 	for lDriverDB, lExpectedSQL := range lDriversDBWithExpectedSQL {
// 		lReturnedSQL := lDriverDB.SQLBuilder().CreateTable(*lProductAttributeTable)
// 		testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, lDriverDB.Name()+".SQLBuilder.CreateTable")
// 	}
// }

// func Test_DriverDB_SQLBuilder_DropTable_ShouldReturnCorrectSQL(t *testing.T) {
// 	lUserTableName := testutils.FakeUserTableName(t)
// 	if lUserTableName == nil {
// 		return
// 	}
// 	lDriversDBWithExpectedSQL := map[ormshift.DriverDB]string{
// 		postgresql.DriverPostgresql: "DROP TABLE user;",
// 		sqlserver.DriverSQLServer:   "DROP TABLE user;",
// 		sqlite.DriverSQLite:         "DROP TABLE user;",
// 		ormshift.DriverDB(-1):           "DROP TABLE user;",
// 	}
// 	for lDriverDB, lExpectedSQL := range lDriversDBWithExpectedSQL {
// 		lReturnedSQL := lDriverDB.SQLBuilder().DropTable(*lUserTableName)
// 		testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, lDriverDB.Name()+".SQLBuilder.DropTable")
// 	}
// }

// func Test_DriverDB_SQLBuilder_AlterTableAddColumn_ShouldReturnCorrectSQL(t *testing.T) {
// 	lUserTableName := testutils.FakeUserTableName(t)
// 	if lUserTableName == nil {
// 		return
// 	}
// 	lUpdatedAtColumn := testutils.FakeUpdatedAtColumn(t)
// 	if lUpdatedAtColumn == nil {
// 		return
// 	}
// 	lDriversDBWithExpectedSQL := map[ormshift.DriverDB]string{
// 		postgresql.DriverPostgresql: "ALTER TABLE user ADD COLUMN updated_at TIMESTAMP(6);",
// 		sqlserver.DriverSQLServer:   "ALTER TABLE user ADD COLUMN updated_at DATETIME2(6);",
// 		sqlite.DriverSQLite:         "ALTER TABLE user ADD COLUMN updated_at DATETIME;",
// 		ormshift.DriverDB(-1):           "ALTER TABLE user ADD COLUMN updated_at <<TYPE_3>>;",
// 	}
// 	for lDriverDB, lExpectedSQL := range lDriversDBWithExpectedSQL {
// 		lReturnedSQL := lDriverDB.SQLBuilder().AlterTableAddColumn(*lUserTableName, *lUpdatedAtColumn)
// 		testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, lDriverDB.Name()+".SQLBuilder.AlterTableAddColumn")
// 	}
// }

// func Test_DriverDB_SQLBuilder_AlterTableDropColumn_ShouldReturnCorrectSQL(t *testing.T) {
// 	lUserTableName := testutils.FakeUserTableName(t)
// 	if lUserTableName == nil {
// 		return
// 	}
// 	lUpdatedAtColumnName := testutils.FakeUpdatedAtColumnName(t)
// 	if lUpdatedAtColumnName == nil {
// 		return
// 	}
// 	lDriversDBWithExpectedSQL := map[ormshift.DriverDB]string{
// 		postgresql.DriverPostgresql: "ALTER TABLE user DROP COLUMN updated_at;",
// 		sqlserver.DriverSQLServer:   "ALTER TABLE user DROP COLUMN updated_at;",
// 		sqlite.DriverSQLite:         "ALTER TABLE user DROP COLUMN updated_at;",
// 		ormshift.DriverDB(-1):           "ALTER TABLE user DROP COLUMN updated_at;",
// 	}
// 	for lDriverDB, lExpectedSQL := range lDriversDBWithExpectedSQL {
// 		lReturnedSQL := lDriverDB.SQLBuilder().AlterTableDropColumn(*lUserTableName, *lUpdatedAtColumnName)
// 		testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, lDriverDB.Name()+".SQLBuilder.AlterTableDropColumn")
// 	}
// }

// func Test_DriverDB_SQLBuilder_Insert_ShouldReturnCorrectSQL(t *testing.T) {
// 	lDriversDB := []ormshift.DriverDB{
// 		postgresql.DriverPostgresql,
// 		sqlserver.DriverSQLServer,
// 		sqlite.DriverSQLite,
// 		ormshift.DriverDB(-1),
// 	}
// 	for _, lDriverDB := range lDriversDB {
// 		lReturnedSQL := lDriverDB.SQLBuilder().Insert("product", []string{"id", "sku", "name", "description"})
// 		lExpectedSQL := "insert into product (id,sku,name,description) values (@id,@sku,@name,@description)"
// 		testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, lDriverDB.Name()+".SQLBuilder.Insert")
// 	}
// }

// func Test_DriverDB_SQLBuilder_InsertWithValues_ShouldReturnCorrectSQL(t *testing.T) {
// 	lDriversDBWithExpectedSQL := map[ormshift.DriverDB]string{
// 		postgresql.DriverPostgresql: "insert into product (id,name,sku) values ($1,$2,$3)",
// 		sqlserver.DriverSQLServer:   "insert into product (id,name,sku) values (@id,@name,@sku)",
// 		sqlite.DriverSQLite:         "insert into product (id,name,sku) values (@id,@name,@sku)",
// 		ormshift.DriverDB(-1):           "insert into product (id,name,sku) values (@id,@name,@sku)",
// 	}
// 	lDriversDBWithExpectedValues := map[ormshift.DriverDB][]any{
// 		postgresql.DriverPostgresql: {1, "Trufa Sabor Amarula 30g Cacaushow", "1.005.12.9"},
// 		sqlserver.DriverSQLServer: {
// 			sql.NamedArg{Name: "id", Value: 1},
// 			sql.NamedArg{Name: "name", Value: "Trufa Sabor Amarula 30g Cacaushow"},
// 			sql.NamedArg{Name: "sku", Value: "1.005.12.9"},
// 		},
// 		sqlite.DriverSQLite: {
// 			sql.NamedArg{Name: "id", Value: 1},
// 			sql.NamedArg{Name: "name", Value: "Trufa Sabor Amarula 30g Cacaushow"},
// 			sql.NamedArg{Name: "sku", Value: "1.005.12.9"},
// 		},
// 		ormshift.DriverDB(-1): {
// 			sql.NamedArg{Name: "id", Value: 1},
// 			sql.NamedArg{Name: "name", Value: "Trufa Sabor Amarula 30g Cacaushow"},
// 			sql.NamedArg{Name: "sku", Value: "1.005.12.9"},
// 		},
// 	}
// 	for lDriverDB, lExpectedSQL := range lDriversDBWithExpectedSQL {
// 		lReturnedSQL, lReturnedValues := lDriverDB.SQLBuilder().InsertWithValues("product", ormshift.ColumnsValues{"id": 1, "sku": "1.005.12.9", "name": "Trufa Sabor Amarula 30g Cacaushow"})
// 		testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, lDriverDB.Name()+".SQLBuilder.InsertWithValues.SQL")
// 		testutils.AssertEqualWithLabel(t, 3, len(lReturnedValues), lDriverDB.Name()+".SQLBuilder.InsertWithValues.Values")
// 		testutils.AssertEqualWithLabel(t, lDriversDBWithExpectedValues[lDriverDB][0], lReturnedValues[0], lDriverDB.Name()+".SQLBuilder.InsertWithValues.Values[0]")
// 		testutils.AssertEqualWithLabel(t, lDriversDBWithExpectedValues[lDriverDB][1], lReturnedValues[1], lDriverDB.Name()+".SQLBuilder.InsertWithValues.Values[1]")
// 		testutils.AssertEqualWithLabel(t, lDriversDBWithExpectedValues[lDriverDB][2], lReturnedValues[2], lDriverDB.Name()+".SQLBuilder.InsertWithValues.Values[2]")
// 	}
// }

// func Test_DriverDB_SQLBuilder_Update_ShouldReturnCorrectSQL(t *testing.T) {
// 	lDriversDB := []ormshift.DriverDB{
// 		postgresql.DriverPostgresql,
// 		sqlserver.DriverSQLServer,
// 		sqlite.DriverSQLite,
// 		ormshift.DriverDB(-1),
// 	}
// 	for _, lDriverDB := range lDriversDB {
// 		lReturnedSQL := lDriverDB.SQLBuilder().Update("product", []string{"sku", "name", "description"}, []string{"id"})
// 		lExpectedSQL := "update product set sku = @sku,name = @name,description = @description where id = @id"
// 		testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, lDriverDB.Name()+".SQLBuilder.Update")
// 	}
// }

// func Test_DriverDB_SQLBuilder_UpdateWithValues_ShouldReturnCorrectSQL(t *testing.T) {
// 	lDriversDBWithExpectedSQL := map[ormshift.DriverDB]string{
// 		postgresql.DriverPostgresql: "update product set sku = $3,name = $2 where id = $1",
// 		sqlserver.DriverSQLServer:   "update product set sku = @sku,name = @name where id = @id",
// 		sqlite.DriverSQLite:         "update product set sku = @sku,name = @name where id = @id",
// 		ormshift.DriverDB(-1):           "update product set sku = @sku,name = @name where id = @id",
// 	}
// 	lDriversDBWithExpectedValues := map[ormshift.DriverDB][]any{
// 		postgresql.DriverPostgresql: {1, "Trufa Sabor Amarula 18g Cacaushow", "1.005.12.5"},
// 		sqlserver.DriverSQLServer: {
// 			sql.NamedArg{Name: "id", Value: 1},
// 			sql.NamedArg{Name: "name", Value: "Trufa Sabor Amarula 18g Cacaushow"},
// 			sql.NamedArg{Name: "sku", Value: "1.005.12.5"},
// 		},
// 		sqlite.DriverSQLite: {
// 			sql.NamedArg{Name: "id", Value: 1},
// 			sql.NamedArg{Name: "name", Value: "Trufa Sabor Amarula 18g Cacaushow"},
// 			sql.NamedArg{Name: "sku", Value: "1.005.12.5"},
// 		},
// 		ormshift.DriverDB(-1): {
// 			sql.NamedArg{Name: "id", Value: 1},
// 			sql.NamedArg{Name: "name", Value: "Trufa Sabor Amarula 18g Cacaushow"},
// 			sql.NamedArg{Name: "sku", Value: "1.005.12.5"},
// 		},
// 	}
// 	for lDriverDB, lExpectedSQL := range lDriversDBWithExpectedSQL {
// 		lReturnedSQL, lReturnedValues := lDriverDB.SQLBuilder().UpdateWithValues(
// 			"product",
// 			[]string{"sku", "name"},
// 			[]string{"id"},
// 			ormshift.ColumnsValues{"id": 1, "sku": "1.005.12.5", "name": "Trufa Sabor Amarula 18g Cacaushow"},
// 		)
// 		testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, lDriverDB.Name()+".SQLBuilder.UpdateWithValues.SQL")
// 		testutils.AssertEqualWithLabel(t, 3, len(lReturnedValues), lDriverDB.Name()+".SQLBuilder.UpdateWithValues.Values")
// 		testutils.AssertEqualWithLabel(t, lDriversDBWithExpectedValues[lDriverDB][0], lReturnedValues[0], lDriverDB.Name()+".SQLBuilder.UpdateWithValues.Values[0]")
// 		testutils.AssertEqualWithLabel(t, lDriversDBWithExpectedValues[lDriverDB][1], lReturnedValues[1], lDriverDB.Name()+".SQLBuilder.UpdateWithValues.Values[1]")
// 		testutils.AssertEqualWithLabel(t, lDriversDBWithExpectedValues[lDriverDB][2], lReturnedValues[2], lDriverDB.Name()+".SQLBuilder.UpdateWithValues.Values[2]")
// 	}
// }

// func Test_DriverDB_SQLBuilder_Delete_ShouldReturnCorrectSQL(t *testing.T) {
// 	lDriversDB := []ormshift.DriverDB{
// 		postgresql.DriverPostgresql,
// 		sqlserver.DriverSQLServer,
// 		sqlite.DriverSQLite,
// 		ormshift.DriverDB(-1),
// 	}
// 	for _, lDriverDB := range lDriversDB {
// 		lReturnedSQL := lDriverDB.SQLBuilder().Delete("product", []string{"id"})
// 		lExpectedSQL := "delete from product where id = @id"
// 		testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, lDriverDB.Name()+".SQLBuilder.Delete")
// 	}
// }

// func Test_DriverDB_SQLBuilder_DeleteWithValues_ShouldReturnCorrectSQL(t *testing.T) {
// 	lDriversDBWithExpectedSQL := map[ormshift.DriverDB]string{
// 		postgresql.DriverPostgresql: "delete from product where id = $1",
// 		sqlserver.DriverSQLServer:   "delete from product where id = @id",
// 		sqlite.DriverSQLite:         "delete from product where id = @id",
// 		ormshift.DriverDB(-1):           "delete from product where id = @id",
// 	}
// 	lDriversDBWithExpectedValues := map[ormshift.DriverDB][]any{
// 		postgresql.DriverPostgresql: {1},
// 		sqlserver.DriverSQLServer:   {sql.NamedArg{Name: "id", Value: 1}},
// 		sqlite.DriverSQLite:         {sql.NamedArg{Name: "id", Value: 1}},
// 		ormshift.DriverDB(-1):           {sql.NamedArg{Name: "id", Value: 1}},
// 	}
// 	for lDriverDB, lExpectedSQL := range lDriversDBWithExpectedSQL {
// 		lReturnedSQL, lReturnedValues := lDriverDB.SQLBuilder().DeleteWithValues("product", ormshift.ColumnsValues{"id": 1})
// 		testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, lDriverDB.Name()+".SQLBuilder.DeleteWithValues.SQL")
// 		testutils.AssertEqualWithLabel(t, 1, len(lReturnedValues), lDriverDB.Name()+".SQLBuilder.DeleteWithValues.Values")
// 		testutils.AssertEqualWithLabel(t, lDriversDBWithExpectedValues[lDriverDB][0], lReturnedValues[0], lDriverDB.Name()+".SQLBuilder.DeleteWithValues.Values[0]")
// 	}
// }

// func Test_DriverDB_SQLBuilder_Select_ShouldReturnCorrectSQL(t *testing.T) {
// 	lDriversDB := []ormshift.DriverDB{
// 		postgresql.DriverPostgresql,
// 		sqlserver.DriverSQLServer,
// 		sqlite.DriverSQLite,
// 		ormshift.DriverDB(-1),
// 	}
// 	for _, lDriverDB := range lDriversDB {
// 		lReturnedSQL := lDriverDB.SQLBuilder().Select("product", []string{"id", "name", "description"}, []string{"sku", "active"})
// 		lExpectedSQL := "select id,name,description from product where sku = @sku and active = @active"
// 		testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, lDriverDB.Name()+".SQLBuilder.Select")
// 	}
// }

// func Test_DriverDB_SQLBuilder_SelectWithValues_ShouldReturnCorrectSQL(t *testing.T) {
// 	lDriversDBWithExpectedSQL := map[ormshift.DriverDB]string{
// 		postgresql.DriverPostgresql: "select id,sku,name,description from product where active = $1 and category_id = $2",
// 		sqlserver.DriverSQLServer:   "select id,sku,name,description from product where active = @active and category_id = @category_id",
// 		sqlite.DriverSQLite:         "select id,sku,name,description from product where active = @active and category_id = @category_id",
// 		ormshift.DriverDB(-1):           "select id,sku,name,description from product where active = @active and category_id = @category_id",
// 	}
// 	lDriversDBWithExpectedValues := map[ormshift.DriverDB][]any{
// 		postgresql.DriverPostgresql: {1, 1},
// 		sqlserver.DriverSQLServer:   {sql.NamedArg{Name: "active", Value: true}, sql.NamedArg{Name: "category_id", Value: 1}},
// 		sqlite.DriverSQLite:         {sql.NamedArg{Name: "active", Value: true}, sql.NamedArg{Name: "category_id", Value: 1}},
// 		ormshift.DriverDB(-1):           {sql.NamedArg{Name: "active", Value: true}, sql.NamedArg{Name: "category_id", Value: 1}},
// 	}
// 	for lDriverDB, lExpectedSQL := range lDriversDBWithExpectedSQL {
// 		lReturnedSQL, lReturnedValues := lDriverDB.SQLBuilder().SelectWithValues("product", []string{"id", "sku", "name", "description"}, ormshift.ColumnsValues{"category_id": 1, "active": true})
// 		testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, lDriverDB.Name()+".SQLBuilder.SelectWithValues.SQL")
// 		testutils.AssertEqualWithLabel(t, 2, len(lReturnedValues), lDriverDB.Name()+".SQLBuilder.SelectWithValues.Values")
// 		testutils.AssertEqualWithLabel(t, lDriversDBWithExpectedValues[lDriverDB][0], lReturnedValues[0], lDriverDB.Name()+".SQLBuilder.SelectWithValues.Values[0]")
// 		testutils.AssertEqualWithLabel(t, lDriversDBWithExpectedValues[lDriverDB][1], lReturnedValues[1], lDriverDB.Name()+".SQLBuilder.SelectWithValues.Values[1]")
// 	}
// }

// func Test_DriverDB_SQLBuilder_SelectWithPagination_ShouldReturnCorrectSQL(t *testing.T) {
// 	lDriversDBWithExpectedSQL := map[ormshift.DriverDB]string{
// 		postgresql.DriverPostgresql: "select * from product LIMIT 10 OFFSET 40",
// 		sqlserver.DriverSQLServer:   "select * from product OFFSET 40 ROWS FETCH NEXT 10 ROWS ONLY",
// 		sqlite.DriverSQLite:         "select * from product LIMIT 10 OFFSET 40",
// 		ormshift.DriverDB(-1):           "select * from product LIMIT 10 OFFSET 40",
// 	}
// 	for lDriverDB, lExpectedSQL := range lDriversDBWithExpectedSQL {
// 		lReturnedSQL := lDriverDB.SQLBuilder().SelectWithPagination("select * from product", 10, 5)
// 		testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, lDriverDB.Name()+".SQLBuilder.SelectWithPagination")
// 	}
// }
