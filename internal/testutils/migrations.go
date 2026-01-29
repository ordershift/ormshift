package testutils

import (
	"fmt"

	"github.com/ordershift/ormshift/migrations"
	"github.com/ordershift/ormshift/schema"
)

// M001 Create_Table_User creates the "user" table.
type M001_Create_Table_User struct{}

func (m M001_Create_Table_User) Up(pMigrator *migrations.Migrator) error {
	userTable := schema.NewTable("user")
	if pMigrator.Database().DBSchema().HasTable(userTable.Name()) {
		return nil
	}

	err := userTable.AddColumns(
		schema.NewColumnParams{
			Name:          "id",
			Type:          schema.Integer,
			AutoIncrement: true,
			PrimaryKey:    true,
			NotNull:       true,
		},
		schema.NewColumnParams{
			Name:       "name",
			Type:       schema.Varchar,
			Size:       50,
			PrimaryKey: false,
			NotNull:    false,
		},
		schema.NewColumnParams{
			Name:       "email",
			Type:       schema.Varchar,
			Size:       120,
			PrimaryKey: false,
			NotNull:    false,
		},
		schema.NewColumnParams{
			Name:       "active",
			Type:       schema.Boolean,
			PrimaryKey: false,
			NotNull:    false,
		},
		schema.NewColumnParams{
			Name:       "ammount",
			Type:       schema.Monetary,
			PrimaryKey: false,
			NotNull:    false,
		},
		schema.NewColumnParams{
			Name:       "percent",
			Type:       schema.Decimal,
			PrimaryKey: false,
			NotNull:    false,
		},
		schema.NewColumnParams{
			Name:       "photo",
			Type:       schema.Binary,
			PrimaryKey: false,
			NotNull:    false,
		},
	)
	if err != nil {
		return err
	}

	_, err = pMigrator.Database().SQLExecutor().Exec(pMigrator.Database().SQLBuilder().CreateTable(userTable))
	return err
}

func (m M001_Create_Table_User) Down(pMigrator *migrations.Migrator) error {
	userTableName := "user"
	if !pMigrator.Database().DBSchema().HasTable(userTableName) {
		return nil
	}

	_, err := pMigrator.Database().SQLExecutor().Exec(pMigrator.Database().SQLBuilder().DropTable(userTableName))
	return err
}

// M002_Alter_Table_User_Add_Column_UpdatedAt adds the "updated_at" column to the "user" table.
type M002_Alter_Table_User_Add_Column_UpdatedAt struct{}

func (m M002_Alter_Table_User_Add_Column_UpdatedAt) Up(pMigrator *migrations.Migrator) error {
	userTableName := "user"
	updatedAtColumn := schema.NewColumn(schema.NewColumnParams{
		Name: "updated_at",
		Type: schema.DateTime,
	})
	if pMigrator.Database().DBSchema().HasColumn(userTableName, updatedAtColumn.Name()) {
		return nil
	}
	_, err := pMigrator.Database().SQLExecutor().Exec(pMigrator.Database().SQLBuilder().AlterTableAddColumn(userTableName, updatedAtColumn))
	return err
}

func (m M002_Alter_Table_User_Add_Column_UpdatedAt) Down(pMigrator *migrations.Migrator) error {
	userTableName := "user"
	updatedAtColumnName := "updated_at"
	if !pMigrator.Database().DBSchema().HasColumn(userTableName, updatedAtColumnName) {
		return nil
	}
	_, err := pMigrator.Database().SQLExecutor().Exec(pMigrator.Database().SQLBuilder().AlterTableDropColumn(userTableName, updatedAtColumnName))
	return err
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
