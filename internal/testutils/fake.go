package testutils

import (
	"testing"

	"github.com/ordershift/ormshift/pkg/core"
)

func FakeProductAttributeTable(t *testing.T) *core.Table {
	lProductAttributeTable, lError := core.NewTable("product_attribute")
	if !AssertNotNilResultAndNilError(t, lProductAttributeTable, lError, "core.NewTable") {
		return nil
	}
	lError = lProductAttributeTable.AddColumn(core.NewColumnParams{
		Name:          "product_id",
		Type:          core.Integer,
		PrimaryKey:    true,
		NotNull:       true,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "ProductAttributeTable.AddColumn") {
		return nil
	}
	lError = lProductAttributeTable.AddColumn(core.NewColumnParams{
		Name:          "attribute_id",
		Type:          core.Integer,
		PrimaryKey:    true,
		NotNull:       true,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "ProductAttributeTable.AddColumn") {
		return nil
	}
	lError = lProductAttributeTable.AddColumn(core.NewColumnParams{
		Name:          "value",
		Type:          core.Varchar,
		Size:          75,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "ProductAttributeTable.AddColumn") {
		return nil
	}
	lError = lProductAttributeTable.AddColumn(core.NewColumnParams{
		Name:          "position",
		Type:          core.Integer,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "ProductAttributeTable.AddColumn") {
		return nil
	}
	return lProductAttributeTable
}

func FakeUserTable(t *testing.T) *core.Table {
	lUserTable, lError := core.NewTable("user")
	if !AssertNotNilResultAndNilError(t, lUserTable, lError, "core.NewTable") {
		return nil
	}
	lError = lUserTable.AddColumn(core.NewColumnParams{
		Name:          "id",
		Type:          core.Integer,
		PrimaryKey:    true,
		NotNull:       true,
		Autoincrement: true,
	})
	if !AssertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(core.NewColumnParams{
		Name:          "email",
		Type:          core.Varchar,
		Size:          80,
		PrimaryKey:    true,
		NotNull:       true,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(core.NewColumnParams{
		Name:          "name",
		Type:          core.Varchar,
		Size:          50,
		PrimaryKey:    false,
		NotNull:       true,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(core.NewColumnParams{
		Name:          "password_hash",
		Type:          core.Varchar,
		Size:          256,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(core.NewColumnParams{
		Name:          "active",
		Type:          core.Boolean,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(core.NewColumnParams{
		Name:          "created_at",
		Type:          core.DateTime,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(core.NewColumnParams{
		Name:          "user_master",
		Type:          core.Integer,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(core.NewColumnParams{
		Name:          "master_user_id",
		Type:          core.Integer,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(core.NewColumnParams{
		Name:          "licence_price",
		Type:          core.Monetary,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(core.NewColumnParams{
		Name:          "relevance",
		Type:          core.Decimal,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(core.NewColumnParams{
		Name:          "photo",
		Type:          core.Binary,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(core.NewColumnParams{
		Name:          "any",
		Type:          core.ColumnType(-1),
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !AssertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	return lUserTable
}

func FakeUserTableName(t *testing.T) *core.TableName {
	lUserTableName, lError := core.NewTableName("user")
	if !AssertNotNilResultAndNilError(t, lUserTableName, lError, "core.NewTableName") {
		return nil
	}
	return lUserTableName
}

func FakeUpdatedAtColumn(t *testing.T) *core.Column {
	lUpdatedAtColumn, lError := core.NewColumn(core.NewColumnParams{
		Name:          "updated_at",
		Type:          core.DateTime,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !AssertNotNilResultAndNilError(t, lUpdatedAtColumn, lError, "core.NewColumn") {
		return nil
	}
	return lUpdatedAtColumn
}

func FakeUpdatedAtColumnName(t *testing.T) *core.ColumnName {
	lUpdatedAtColumnName, lError := core.NewColumnName("updated_at")
	if !AssertNotNilResultAndNilError(t, lUpdatedAtColumnName, lError, "core.NewColumnName") {
		return nil
	}
	return lUpdatedAtColumnName
}
