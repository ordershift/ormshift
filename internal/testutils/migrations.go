package testutils

import (
	"github.com/ordershift/ormshift/migrations"
	"github.com/ordershift/ormshift/schema"
)

type M001_Create_Table_User struct{}

func (m M001_Create_Table_User) Up(pMigrator *migrations.Migrator) error {
	lUserTable, lError := schema.NewTable("user")
	if lError != nil {
		return lError
	}
	if pMigrator.DBSchema().ExistsTable(lUserTable.Name()) {
		return nil
	}
	lUserTable.AddColumn(schema.NewColumnParams{
		Name:          "id",
		Type:          schema.Integer,
		Autoincrement: true,
		PrimaryKey:    true,
		NotNull:       true,
	})
	lUserTable.AddColumn(schema.NewColumnParams{
		Name:       "name",
		Type:       schema.Varchar,
		Size:       50,
		PrimaryKey: false,
		NotNull:    false,
	})
	lUserTable.AddColumn(schema.NewColumnParams{
		Name:       "email",
		Type:       schema.Varchar,
		Size:       120,
		PrimaryKey: false,
		NotNull:    false,
	})
	lUserTable.AddColumn(schema.NewColumnParams{
		Name:       "active",
		Type:       schema.Boolean,
		PrimaryKey: false,
		NotNull:    false,
	})
	lUserTable.AddColumn(schema.NewColumnParams{
		Name:       "ammount",
		Type:       schema.Monetary,
		PrimaryKey: false,
		NotNull:    false,
	})
	lUserTable.AddColumn(schema.NewColumnParams{
		Name:       "percent",
		Type:       schema.Decimal,
		PrimaryKey: false,
		NotNull:    false,
	})
	lUserTable.AddColumn(schema.NewColumnParams{
		Name:       "photo",
		Type:       schema.Binary,
		PrimaryKey: false,
		NotNull:    false,
	})
	_, lError = pMigrator.DB().Exec(pMigrator.SQLBuilder().CreateTable(*lUserTable))
	if lError != nil {
		return lError
	}
	return nil
}

func (m M001_Create_Table_User) Down(pMigrator *migrations.Migrator) error {
	lUserTableName, lError := schema.NewTableName("user")
	if lError != nil {
		return lError
	}
	if !pMigrator.DBSchema().ExistsTable(*lUserTableName) {
		return nil
	}
	_, lError = pMigrator.DB().Exec(pMigrator.SQLBuilder().DropTable(*lUserTableName))
	if lError != nil {
		return lError
	}
	return nil
}

type M002_Alter_Table_Usaer_Add_Column_UpdatedAt struct{}

func (m M002_Alter_Table_Usaer_Add_Column_UpdatedAt) Up(pMigrator *migrations.Migrator) error {
	lUserTableName, lError := schema.NewTableName("user")
	if lError != nil {
		return lError
	}
	lUpdatedAtColumn, lError := schema.NewColumn(schema.NewColumnParams{
		Name: "updated_at",
		Type: schema.DateTime,
	})
	if lError != nil {
		return lError
	}
	if pMigrator.DBSchema().ExistsTableColumn(*lUserTableName, lUpdatedAtColumn.Name()) {
		return nil
	}
	_, lError = pMigrator.DB().Exec(pMigrator.SQLBuilder().AlterTableAddColumn(*lUserTableName, *lUpdatedAtColumn))
	if lError != nil {
		return lError
	}
	return nil
}

func (m M002_Alter_Table_Usaer_Add_Column_UpdatedAt) Down(pMigrator *migrations.Migrator) error {
	lUserTableName, lError := schema.NewTableName("user")
	if lError != nil {
		return lError
	}
	lUpdatedAtColumnName, lError := schema.NewColumnName("updated_at")
	if lError != nil {
		return lError
	}
	if !pMigrator.DBSchema().ExistsTableColumn(*lUserTableName, *lUpdatedAtColumnName) {
		return nil
	}
	_, lError = pMigrator.DB().Exec(pMigrator.SQLBuilder().AlterTableDropColumn(*lUserTableName, *lUpdatedAtColumnName))
	if lError != nil {
		return lError
	}
	return nil
}
