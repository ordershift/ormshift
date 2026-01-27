package testutils

import (
	"database/sql"
	"testing"

	"github.com/ordershift/ormshift/schema"
)

func FakeProductAttributeTable(t *testing.T) schema.Table {
	lProductAttributeTable := schema.NewTable("product_attribute")
	lError := lProductAttributeTable.AddColumns(
		schema.NewColumnParams{
			Name:       "product_id",
			Type:       schema.Integer,
			PrimaryKey: true,
			NotNull:    true,
		},
		schema.NewColumnParams{
			Name:       "attribute_id",
			Type:       schema.Integer,
			PrimaryKey: true,
			NotNull:    true,
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
	if !AssertNilError(t, lError, "ProductAttributeTable.AddColumns") {
		panic(lError)
	}
	return lProductAttributeTable
}

func FakeUserTable(t *testing.T) schema.Table {
	lUserTable := schema.NewTable("user")
	lError := lUserTable.AddColumns(
		schema.NewColumnParams{
			Name:          "id",
			Type:          schema.Integer,
			PrimaryKey:    true,
			NotNull:       true,
			AutoIncrement: true,
		},
		schema.NewColumnParams{
			Name:          "email",
			Type:          schema.Varchar,
			Size:          80,
			PrimaryKey:    true,
			NotNull:       true,
			AutoIncrement: false,
		},
		schema.NewColumnParams{
			Name:          "name",
			Type:          schema.Varchar,
			Size:          50,
			PrimaryKey:    false,
			NotNull:       true,
			AutoIncrement: false,
		},
		schema.NewColumnParams{
			Name:          "password_hash",
			Type:          schema.Varchar,
			Size:          256,
			PrimaryKey:    false,
			NotNull:       false,
			AutoIncrement: false,
		},
		schema.NewColumnParams{
			Name:          "active",
			Type:          schema.Boolean,
			PrimaryKey:    false,
			NotNull:       false,
			AutoIncrement: false,
		},
		schema.NewColumnParams{
			Name:          "created_at",
			Type:          schema.DateTime,
			PrimaryKey:    false,
			NotNull:       false,
			AutoIncrement: false,
		},
		schema.NewColumnParams{
			Name:          "user_master",
			Type:          schema.Integer,
			PrimaryKey:    false,
			NotNull:       false,
			AutoIncrement: false,
		},
		schema.NewColumnParams{
			Name:          "master_user_id",
			Type:          schema.Integer,
			PrimaryKey:    false,
			NotNull:       false,
			AutoIncrement: false,
		},
		schema.NewColumnParams{
			Name:          "licence_price",
			Type:          schema.Monetary,
			PrimaryKey:    false,
			NotNull:       false,
			AutoIncrement: false,
		},
		schema.NewColumnParams{
			Name:          "relevance",
			Type:          schema.Decimal,
			PrimaryKey:    false,
			NotNull:       false,
			AutoIncrement: false,
		},
		schema.NewColumnParams{
			Name:          "photo",
			Type:          schema.Binary,
			PrimaryKey:    false,
			NotNull:       false,
			AutoIncrement: false,
		},
		schema.NewColumnParams{
			Name:          "any",
			Type:          schema.ColumnType(-1),
			PrimaryKey:    false,
			NotNull:       false,
			AutoIncrement: false,
		},
	)
	if !AssertNilError(t, lError, "UserTable.AddColumns") {
		panic(lError)
	}
	return lUserTable
}

func FakeUserTableName(t *testing.T) string {
	return "user"
}

func FakeUpdatedAtColumn(t *testing.T) schema.Column {
	lUpdatedAtColumn := schema.NewColumn(schema.NewColumnParams{
		Name:          "updated_at",
		Type:          schema.DateTime,
		PrimaryKey:    false,
		NotNull:       false,
		AutoIncrement: false,
	})
	return lUpdatedAtColumn
}

func FakeUpdatedAtColumnName(t *testing.T) string {
	return "updated_at"
}

func FakeInteroperateSQLCommandWithNamedArgsFunc(command string, namedArgs ...sql.NamedArg) (string, []any) {
	lArgs := make([]any, len(namedArgs))
	for i, v := range namedArgs {
		lArgs[i] = v
	}
	return "command has been modified", lArgs
}

func FakeColumnDefinitionFunc(column schema.Column) string {
	return "fake"
}

func FakeQuoteIdentifierFunc(identifier string) string {
	return "quoted_" + identifier
}
