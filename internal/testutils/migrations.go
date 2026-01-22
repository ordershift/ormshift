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
	if pMigrator.Database().DBSchema().ExistsTable(lUserTable.Name()) {
		return nil
	}
	columns := []schema.NewColumnParams{
		{
			Name:          "id",
			Type:          schema.Integer,
			Autoincrement: true,
			PrimaryKey:    true,
			NotNull:       true,
		},
		{
			Name:       "name",
			Type:       schema.Varchar,
			Size:       50,
			PrimaryKey: false,
			NotNull:    false,
		},
		{
			Name:       "email",
			Type:       schema.Varchar,
			Size:       120,
			PrimaryKey: false,
			NotNull:    false,
		},
		{
			Name:       "active",
			Type:       schema.Boolean,
			PrimaryKey: false,
			NotNull:    false,
		},
		{
			Name:       "ammount",
			Type:       schema.Monetary,
			PrimaryKey: false,
			NotNull:    false,
		},
		{
			Name:       "percent",
			Type:       schema.Decimal,
			PrimaryKey: false,
			NotNull:    false,
		},
		{
			Name:       "photo",
			Type:       schema.Binary,
			PrimaryKey: false,
			NotNull:    false,
		},
	}

	for _, col := range columns {
		if err := lUserTable.AddColumn(col); err != nil {
			return err
		}
	}

	_, lError = pMigrator.Database().DB().Exec(pMigrator.Database().SQLBuilder().CreateTable(*lUserTable))
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
	if !pMigrator.Database().DBSchema().ExistsTable(*lUserTableName) {
		return nil
	}
	_, lError = pMigrator.Database().DB().Exec(pMigrator.Database().SQLBuilder().DropTable(*lUserTableName))
	if lError != nil {
		return lError
	}
	return nil
}

type M002_Alter_Table_User_Add_Column_UpdatedAt struct{}

func (m M002_Alter_Table_User_Add_Column_UpdatedAt) Up(pMigrator *migrations.Migrator) error {
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
	if pMigrator.Database().DBSchema().ExistsTableColumn(*lUserTableName, lUpdatedAtColumn.Name()) {
		return nil
	}
	_, lError = pMigrator.Database().DB().Exec(pMigrator.Database().SQLBuilder().AlterTableAddColumn(*lUserTableName, *lUpdatedAtColumn))
	if lError != nil {
		return lError
	}
	return nil
}

func (m M002_Alter_Table_User_Add_Column_UpdatedAt) Down(pMigrator *migrations.Migrator) error {
	lUserTableName, lError := schema.NewTableName("user")
	if lError != nil {
		return lError
	}
	lUpdatedAtColumnName, lError := schema.NewColumnName("updated_at")
	if lError != nil {
		return lError
	}
	if !pMigrator.Database().DBSchema().ExistsTableColumn(*lUserTableName, *lUpdatedAtColumnName) {
		return nil
	}
	_, lError = pMigrator.Database().DB().Exec(pMigrator.Database().SQLBuilder().AlterTableDropColumn(*lUserTableName, *lUpdatedAtColumnName))
	if lError != nil {
		return lError
	}
	return nil
}
