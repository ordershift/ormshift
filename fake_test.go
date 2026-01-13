package ormshift_test

import (
	"testing"

	"github.com/ordershift/ormshift"
)

func fakeProductAttributeTable(t *testing.T) *ormshift.Table {
	lProductAttributeTable, lError := ormshift.NewTable("product_attribute")
	if !assertNotNilResultAndNilError(t, lProductAttributeTable, lError, "ormshift.NewTable") {
		return nil
	}
	lError = lProductAttributeTable.AddColumn(ormshift.NewColumnParams{
		Name:          "product_id",
		Type:          ormshift.Integer,
		PrimaryKey:    true,
		NotNull:       true,
		Autoincrement: false,
	})
	if !assertNilError(t, lError, "ProductAttributeTable.AddColumn") {
		return nil
	}
	lError = lProductAttributeTable.AddColumn(ormshift.NewColumnParams{
		Name:          "attribute_id",
		Type:          ormshift.Integer,
		PrimaryKey:    true,
		NotNull:       true,
		Autoincrement: false,
	})
	if !assertNilError(t, lError, "ProductAttributeTable.AddColumn") {
		return nil
	}
	lError = lProductAttributeTable.AddColumn(ormshift.NewColumnParams{
		Name:          "value",
		Type:          ormshift.Varchar,
		Size:          75,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !assertNilError(t, lError, "ProductAttributeTable.AddColumn") {
		return nil
	}
	lError = lProductAttributeTable.AddColumn(ormshift.NewColumnParams{
		Name:          "position",
		Type:          ormshift.Integer,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !assertNilError(t, lError, "ProductAttributeTable.AddColumn") {
		return nil
	}
	return lProductAttributeTable
}

func fakeUserTable(t *testing.T) *ormshift.Table {
	lUserTable, lError := ormshift.NewTable("user")
	if !assertNotNilResultAndNilError(t, lUserTable, lError, "ormshift.NewTable") {
		return nil
	}
	lError = lUserTable.AddColumn(ormshift.NewColumnParams{
		Name:          "id",
		Type:          ormshift.Integer,
		PrimaryKey:    true,
		NotNull:       true,
		Autoincrement: true,
	})
	if !assertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(ormshift.NewColumnParams{
		Name:          "email",
		Type:          ormshift.Varchar,
		Size:          80,
		PrimaryKey:    true,
		NotNull:       true,
		Autoincrement: false,
	})
	if !assertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(ormshift.NewColumnParams{
		Name:          "name",
		Type:          ormshift.Varchar,
		Size:          50,
		PrimaryKey:    false,
		NotNull:       true,
		Autoincrement: false,
	})
	if !assertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(ormshift.NewColumnParams{
		Name:          "password_hash",
		Type:          ormshift.Varchar,
		Size:          256,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !assertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(ormshift.NewColumnParams{
		Name:          "active",
		Type:          ormshift.Boolean,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !assertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(ormshift.NewColumnParams{
		Name:          "created_at",
		Type:          ormshift.DateTime,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !assertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(ormshift.NewColumnParams{
		Name:          "user_master",
		Type:          ormshift.Integer,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !assertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(ormshift.NewColumnParams{
		Name:          "master_user_id",
		Type:          ormshift.Integer,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !assertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(ormshift.NewColumnParams{
		Name:          "licence_price",
		Type:          ormshift.Monetary,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !assertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(ormshift.NewColumnParams{
		Name:          "relevance",
		Type:          ormshift.Decimal,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !assertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(ormshift.NewColumnParams{
		Name:          "photo",
		Type:          ormshift.Binary,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !assertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	lError = lUserTable.AddColumn(ormshift.NewColumnParams{
		Name:          "any",
		Type:          ormshift.ColumnType(-1),
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !assertNilError(t, lError, "UserTable.AddColumn") {
		return nil
	}
	return lUserTable
}

func fakeUserTableName(t *testing.T) *ormshift.TableName {
	lUserTableName, lError := ormshift.NewTableName("user")
	if !assertNotNilResultAndNilError(t, lUserTableName, lError, "ormshift.NewTableName") {
		return nil
	}
	return lUserTableName
}

func fakeUpdatedAtColumn(t *testing.T) *ormshift.Column {
	lUpdatedAtColumn, lError := ormshift.NewColumn(ormshift.NewColumnParams{
		Name:          "updated_at",
		Type:          ormshift.DateTime,
		PrimaryKey:    false,
		NotNull:       false,
		Autoincrement: false,
	})
	if !assertNotNilResultAndNilError(t, lUpdatedAtColumn, lError, "ormshift.NewColumn") {
		return nil
	}
	return lUpdatedAtColumn
}

func fakeUpdatedAtColumnName(t *testing.T) *ormshift.ColumnName {
	lUpdatedAtColumnName, lError := ormshift.NewColumnName("updated_at")
	if !assertNotNilResultAndNilError(t, lUpdatedAtColumnName, lError, "ormshift.NewColumnName") {
		return nil
	}
	return lUpdatedAtColumnName
}
