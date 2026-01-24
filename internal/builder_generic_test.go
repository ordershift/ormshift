package internal_test

import (
	"testing"

	"github.com/ordershift/ormshift/internal"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestCreateTable(t *testing.T) {
	lSQLBuilder := internal.GenericSQLBuilder{}

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
