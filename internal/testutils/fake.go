package testutils

import (
	"testing"

	"github.com/ordershift/ormshift/schema"
)

func FakeProductAttributeTable(t *testing.T) *schema.Table {
	lProductAttributeTable, lError := schema.NewTable("product_attribute")
	if !AssertNotNilResultAndNilError(t, lProductAttributeTable, lError, "schema.NewTable") {
		return nil
	}
	lError = lProductAttributeTable.AddColumn(schema.NewColumnParams{
		Name:          "product_id",
		Type:          schema.Integer,
		PrimaryKey:    true,
		NotNull:       true,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "ProductAttributeTable.AddColumn") {
		return nil
	}
	lError = lProductAttributeTable.AddColumn(schema.NewColumnParams{
		Name:          "attribute_id",
		Type:          schema.Integer,
		PrimaryKey:    true,
		NotNull:       true,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "ProductAttributeTable.AddColumn") {
		return nil
	}
	lError = lProductAttributeTable.AddColumn(schema.NewColumnParams{
		Name:          "value",
		Type:          schema.Varchar,
		Size:          75,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "ProductAttributeTable.AddColumn") {
		return nil
	}
	lError = lProductAttributeTable.AddColumn(schema.NewColumnParams{
		Name:          "position",
		Type:          schema.Integer,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "ProductAttributeTable.AddColumn") {
		return nil
	}
	return lProductAttributeTable
}

func FakeUserTable(t *testing.T) *schema.Table {
	lUserTable, lError := schema.NewTable("user")
	if !AssertNotNilResultAndNilError(t, lUserTable, lError, "schema.NewTable") {
		return nil
	}
	lError = lUserTable.AddColumn(schema.NewColumnParams{
		Name:          "id",
		Type:          schema.Integer,
		PrimaryKey:    true,
		NotNull:       true,
		Autoincrement: true,
	})
	if !AssertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(schema.NewColumnParams{
		Name:          "email",
		Type:          schema.Varchar,
		Size:          80,
		PrimaryKey:    true,
		NotNull:       true,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(schema.NewColumnParams{
		Name:          "name",
		Type:          schema.Varchar,
		Size:          50,
		PrimaryKey:    false,
		NotNull:       true,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(schema.NewColumnParams{
		Name:          "password_hash",
		Type:          schema.Varchar,
		Size:          256,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(schema.NewColumnParams{
		Name:          "active",
		Type:          schema.Boolean,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(schema.NewColumnParams{
		Name:          "created_at",
		Type:          schema.DateTime,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(schema.NewColumnParams{
		Name:          "user_master",
		Type:          schema.Integer,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(schema.NewColumnParams{
		Name:          "master_user_id",
		Type:          schema.Integer,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(schema.NewColumnParams{
		Name:          "licence_price",
		Type:          schema.Monetary,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(schema.NewColumnParams{
		Name:          "relevance",
		Type:          schema.Decimal,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(schema.NewColumnParams{
		Name:          "photo",
		Type:          schema.Binary,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(schema.NewColumnParams{
		Name:          "any",
		Type:          schema.ColumnType(-1),
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	return lUserTable
}

func FakeUserTableName(t *testing.T) *schema.TableName {
	lUserTableName, lError := schema.NewTableName("user")
	if !AssertNotNilResultAndNilError(t, lUserTableName, lError, "schema.NewTableName") {
		return nil
	}
	return lUserTableName
}

func FakeUpdatedAtColumn(t *testing.T) *schema.Column {
	lUpdatedAtColumn, lError := schema.NewColumn(schema.NewColumnParams{
		Name:          "updated_at",
		Type:          schema.DateTime,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !AssertNotNilResultAndNilError(t, lUpdatedAtColumn, lError, "schema.NewColumn") {
		return nil
	}
	return lUpdatedAtColumn
}

func FakeUpdatedAtColumnName(t *testing.T) *schema.ColumnName {
	lUpdatedAtColumnName, lError := schema.NewColumnName("updated_at")
	if !AssertNotNilResultAndNilError(t, lUpdatedAtColumnName, lError, "schema.NewColumnName") {
		return nil
	}
	return lUpdatedAtColumnName
}
