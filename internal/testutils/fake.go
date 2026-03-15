package testutils

import (
	"database/sql"
	"testing"

	"github.com/ordershift/ormshift/schema"
)

func FakeProductAttributeTable(t *testing.T) schema.Table {
	productAttributeTable := schema.NewTable("product_attribute")
	err := productAttributeTable.AddColumns(
		schema.NewColumnParams{
			Name:    "product_id",
			Type:    schema.Integer,
			NotNull: true,
		},
		schema.NewColumnParams{
			Name:    "attribute_id",
			Type:    schema.Integer,
			NotNull: true,
		},
		schema.NewColumnParams{
			Name: "value",
			Type: schema.Varchar,
			Size: 75,
		},
		schema.NewColumnParams{
			Name: "position",
			Type: schema.Integer,
		},
	)
	if !AssertNilError(t, err, "ProductAttributeTable.AddColumns") {
		panic(err)
	}

	err = productAttributeTable.PrimaryKey("product_id", "attribute_id")
	if !AssertNilError(t, err, "ProductAttributeTable.PrimaryKey") {
		panic(err)
	}

	err = productAttributeTable.AddForeignKey([]string{"product_id"}, "product", []string{"id"})
	if !AssertNilError(t, err, "ProductAttributeTable.AddForeignKey product_id") {
		panic(err)
	}
	err = productAttributeTable.AddForeignKey([]string{"attribute_id"}, "attribute", []string{"id"})
	if !AssertNilError(t, err, "ProductAttributeTable.AddForeignKey attribute_id") {
		panic(err)
	}

	return productAttributeTable
}

// FakeTableWithCompositeFKAndUC returns a table "booking" with resource_id, slot_date, slot_hour, guest_id;
// PK (resource_id, slot_date, slot_hour); composite FK (resource_id, slot_date) -> resource_schedule(resource_id, schedule_date);
// composite UC (resource_id, slot_date, slot_hour).
func FakeTableWithCompositeFKAndUC(t *testing.T) schema.Table {
	tbl := schema.NewTable("booking")
	err := tbl.AddColumns(
		schema.NewColumnParams{Name: "resource_id", Type: schema.Integer, NotNull: true},
		schema.NewColumnParams{Name: "slot_date", Type: schema.Varchar, Size: 10, NotNull: true},
		schema.NewColumnParams{Name: "slot_hour", Type: schema.Integer, NotNull: true},
		schema.NewColumnParams{Name: "guest_id", Type: schema.Integer, NotNull: false},
	)
	if !AssertNilError(t, err, "TableWithCompositeFKAndUC.AddColumns") {
		panic(err)
	}
	err = tbl.PrimaryKey("resource_id", "slot_date", "slot_hour")
	if !AssertNilError(t, err, "TableWithCompositeFKAndUC.PrimaryKey") {
		panic(err)
	}
	err = tbl.AddForeignKey([]string{"resource_id", "slot_date"}, "resource_schedule", []string{"resource_id", "schedule_date"})
	if !AssertNilError(t, err, "TableWithCompositeFKAndUC.AddForeignKey") {
		panic(err)
	}
	err = tbl.AddUniqueConstraint("resource_id", "slot_date", "slot_hour")
	if !AssertNilError(t, err, "TableWithCompositeFKAndUC.AddUniqueConstraint") {
		panic(err)
	}
	return tbl
}

func FakeUserTable(t *testing.T) schema.Table {
	userTable := schema.NewTable("user")
	err := userTable.AddColumns(
		schema.NewColumnParams{
			Name:          "id",
			Type:          schema.Integer,
			NotNull:       true,
			AutoIncrement: true,
		},
		schema.NewColumnParams{
			Name:          "email",
			Type:          schema.Varchar,
			Size:          80,
			NotNull:       true,
			AutoIncrement: false,
		},
		schema.NewColumnParams{
			Name:          "name",
			Type:          schema.Varchar,
			Size:          50,
			NotNull:       true,
			AutoIncrement: false,
		},
		schema.NewColumnParams{
			Name:          "password_hash",
			Type:          schema.Varchar,
			Size:          256,
			NotNull:       false,
			AutoIncrement: false,
		},
		schema.NewColumnParams{
			Name:          "active",
			Type:          schema.Boolean,
			NotNull:       false,
			AutoIncrement: false,
		},
		schema.NewColumnParams{
			Name:          "created_at",
			Type:          schema.DateTime,
			NotNull:       false,
			AutoIncrement: false,
		},
		schema.NewColumnParams{
			Name:          "updated_at",
			Type:          schema.DateTimeOffSet,
			NotNull:       false,
			AutoIncrement: false,
		},
		schema.NewColumnParams{
			Name:          "user_master",
			Type:          schema.Integer,
			NotNull:       false,
			AutoIncrement: false,
		},
		schema.NewColumnParams{
			Name:          "master_user_id",
			Type:          schema.Integer,
			NotNull:       false,
			AutoIncrement: false,
		},
		schema.NewColumnParams{
			Name:          "licence_price",
			Type:          schema.Monetary,
			NotNull:       false,
			AutoIncrement: false,
		},
		schema.NewColumnParams{
			Name:          "relevance",
			Type:          schema.Decimal,
			NotNull:       false,
			AutoIncrement: false,
		},
		schema.NewColumnParams{
			Name:          "photo",
			Type:          schema.Binary,
			NotNull:       false,
			AutoIncrement: false,
		},
		schema.NewColumnParams{
			Name:          "any",
			Type:          schema.ColumnType(-1),
			NotNull:       false,
			AutoIncrement: false,
		},
	)
	if !AssertNilError(t, err, "UserTable.AddColumns") {
		panic(err)
	}

	err = userTable.PrimaryKey("id")
	if !AssertNilError(t, err, "UserTable.PrimaryKey") {
		panic(err)
	}

	err = userTable.AddUniqueConstraint("email")
	if !AssertNilError(t, err, "UserTable.AddUniqueConstraint email") {
		panic(err)
	}

	return userTable
}

func FakeUserTableName(t *testing.T) string {
	return "user"
}

func FakeUpdatedAtColumn(t *testing.T) schema.Column {
	updatedAtColumn := schema.NewColumn(schema.NewColumnParams{
		Name:          "updated_at",
		Type:          schema.DateTime,
		NotNull:       false,
		AutoIncrement: false,
	})
	return updatedAtColumn
}

func FakeCreatedAtColumn(t *testing.T) schema.Column {
	return schema.NewColumn(schema.NewColumnParams{
		Name:          "created_at",
		Type:          schema.DateTimeOffSet,
		NotNull:       true,
		AutoIncrement: false,
	})
}

func FakeActivatedAtColumn(t *testing.T) schema.Column {
	return schema.NewColumn(schema.NewColumnParams{
		Name:          "activated_at",
		Type:          schema.DateTime,
		NotNull:       true,
		AutoIncrement: false,
	})
}

func FakeScoreColumn(t *testing.T) schema.Column {
	return schema.NewColumn(schema.NewColumnParams{
		Name:          "score",
		Type:          schema.Integer,
		NotNull:       true,
		AutoIncrement: false,
	})
}

func FakePriceColumn(t *testing.T) schema.Column {
	return schema.NewColumn(schema.NewColumnParams{
		Name:          "price",
		Type:          schema.Monetary,
		NotNull:       true,
		AutoIncrement: false,
	})
}

func FakeNameColumn(t *testing.T) schema.Column {
	return schema.NewColumn(schema.NewColumnParams{
		Name:          "name",
		Type:          schema.Varchar,
		Size:          50,
		NotNull:       true,
		AutoIncrement: false,
	})
}

func FakeUpdatedAtColumnName(t *testing.T) string {
	return "updated_at"
}

func FakeInteroperateSQLCommandWithNamedArgsFunc(command string, args ...sql.NamedArg) (string, []any) {
	a := make([]any, len(args))
	for i, v := range args {
		a[i] = v
	}
	return "command has been modified", a
}

func FakeColumnDefinitionFunc(column schema.Column) string {
	return "fake"
}

func FakeQuoteIdentifierFunc(identifier string) string {
	return "quoted_" + identifier
}
