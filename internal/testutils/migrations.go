package testutils

import (
	"fmt"

	"github.com/ordershift/ormshift/migrations"
	"github.com/ordershift/ormshift/schema"
)

// M001 Create_Table_User creates the "user" table.
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
			AutoIncrement: true,
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

// M002_Alter_Table_User_Add_Column_UpdatedAt adds the "updated_at" column to the "user" table.
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

// M003_Bad_Migration_Fails_To_Apply is a migration that always fails to apply.
type M003_Bad_Migration_Fails_To_Apply struct{}

func (m M003_Bad_Migration_Fails_To_Apply) Up(pMigrator *migrations.Migrator) error {
	return fmt.Errorf("intentionally failed to Up")
}
func (m M003_Bad_Migration_Fails_To_Apply) Down(pMigrator *migrations.Migrator) error {
	return nil
}

// M004_Bad_Migration_Fails_To_Revert is a migration that always fails to revert.
type M004_Bad_Migration_Fails_To_Revert struct{}

func (m M004_Bad_Migration_Fails_To_Revert) Up(pMigrator *migrations.Migrator) error {
	return nil
}
func (m M004_Bad_Migration_Fails_To_Revert) Down(pMigrator *migrations.Migrator) error {
	return fmt.Errorf("intentionally failed to Down")
}

// M005_Blank_Migration is a migration that does nothing, always succeeding regardless of direction and database state.
type M005_Blank_Migration struct{}

func (m M005_Blank_Migration) Up(pMigrator *migrations.Migrator) error {
	return nil
}
func (m M005_Blank_Migration) Down(pMigrator *migrations.Migrator) error {
	return nil
}
