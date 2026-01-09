package ormshift_test

import (
	"database/sql"
	"testing"

	"github.com/ordershift/ormshift"
)

func Test_Migrate_ShouldExecuteWithSuccess(t *testing.T) {
	lDB, lError := sql.Open(ormshift.DriverSQLite.Name(), ormshift.DriverSQLite.ConnectionString(ormshift.ConnectionParams{InMemory: true}))
	if !assertNilError(t, lError, "sql.Open") {
		return
	}
	defer lDB.Close()
	lMigrationManager, lError := ormshift.Migrate(
		lDB,
		ormshift.DriverSQLite,
		m001_Create_Table_User{},
		m002_Alter_Table_Usaer_Add_Column_UpdatedAt{},
	)
	if !assertNotNilResultAndNilError(t, lMigrationManager, lError, "ormshift.Migrate") {
		return
	}
	lUserTableName, lError := ormshift.NewTableName("user")
	if !assertNilError(t, lError, "ormshift.NewTableName") {
		return
	}
	lUpdatedAtColumnName, lError := ormshift.NewColumnName("updated_at")
	if !assertNilError(t, lError, "ormshift.NewColumnName") {
		return
	}
	assertEqualWithLabel(t, true, lMigrationManager.DBSchema().ExistsTableColumn(*lUserTableName, *lUpdatedAtColumnName), "MigrationManager.DBSchema.ExistsTableColumn[user.updated_at]")
	assertEqualWithLabel(t, 2, len(lMigrationManager.UpedMigrationsNames()), "len(MigrationManager.UpedMigrationsNames)")
}

func Test_Migrate_ShouldExecuteWithSuccess_WhenTwiceExecute(t *testing.T) {
	lDB, lError := sql.Open(ormshift.DriverSQLite.Name(), ormshift.DriverSQLite.ConnectionString(ormshift.ConnectionParams{InMemory: true}))
	if !assertNilError(t, lError, "sql.Open") {
		return
	}
	defer lDB.Close()

	lMigrationManager, lError := ormshift.Migrate(
		lDB,
		ormshift.DriverSQLite,
		m001_Create_Table_User{},
		m002_Alter_Table_Usaer_Add_Column_UpdatedAt{},
	)
	if !assertNotNilResultAndNilError(t, lMigrationManager, lError, "ormshift.Migrate") {
		return
	}

	lMigrationManager, lError = ormshift.Migrate(
		lDB,
		ormshift.DriverSQLite,
		m001_Create_Table_User{},
		m002_Alter_Table_Usaer_Add_Column_UpdatedAt{},
	)
	if !assertNotNilResultAndNilError(t, lMigrationManager, lError, "ormshift.Migrate") {
		return
	}

	lUserTableName, lError := ormshift.NewTableName("user")
	if !assertNilError(t, lError, "ormshift.NewTableName") {
		return
	}
	lUpdatedAtColumnName, lError := ormshift.NewColumnName("updated_at")
	if !assertNilError(t, lError, "ormshift.NewColumnName") {
		return
	}
	assertEqualWithLabel(t, true, lMigrationManager.DBSchema().ExistsTableColumn(*lUserTableName, *lUpdatedAtColumnName), "MigrationManager.DBSchema.ExistsTableColumn[user.updated_at]")
	assertEqualWithLabel(t, 2, len(lMigrationManager.UpedMigrationsNames()), "len(MigrationManager.UpedMigrationsNames)")
}

func Test_Migrate_ShouldFail_WhenNilDB(t *testing.T) {
	lMigrationManager, lError := ormshift.Migrate(
		nil,
		ormshift.DriverSQLite,
		m001_Create_Table_User{},
		m002_Alter_Table_Usaer_Add_Column_UpdatedAt{},
	)
	if !assertNilResultAndNotNilError(t, lMigrationManager, lError, "ormshift.Migrate") {
		return
	}
	assertErrorMessage(t, "sql.DB cannot be nil", lError, "ormshift.Migrate")
}

func Test_Migrate_ShouldFail_WhenInvalidDriverDB(t *testing.T) {
	lDB, lError := sql.Open(ormshift.DriverSQLite.Name(), ormshift.DriverSQLite.ConnectionString(ormshift.ConnectionParams{InMemory: true}))
	if !assertNilError(t, lError, "sql.Open") {
		return
	}
	defer lDB.Close()

	lMigrationManager, lError := ormshift.Migrate(
		lDB,
		ormshift.DriverDB(-5),
		m001_Create_Table_User{},
		m002_Alter_Table_Usaer_Add_Column_UpdatedAt{},
	)
	if !assertNilResultAndNotNilError(t, lMigrationManager, lError, "ormshift.Migrate") {
		return
	}
	assertErrorMessage(t, "driver db should be valid", lError, "ormshift.Migrate")
}

func Test_Migrate_ShouldFail_WhenClosedDB(t *testing.T) {
	lDB, lError := sql.Open(ormshift.DriverSQLite.Name(), ormshift.DriverSQLite.ConnectionString(ormshift.ConnectionParams{InMemory: true}))
	if !assertNilError(t, lError, "sql.Open") {
		return
	}
	lDB.Close()

	lMigrationManager, lError := ormshift.Migrate(
		lDB,
		ormshift.DriverSQLite,
		m001_Create_Table_User{},
		m002_Alter_Table_Usaer_Add_Column_UpdatedAt{},
	)
	if !assertNilResultAndNotNilError(t, lMigrationManager, lError, "ormshift.Migrate") {
		return
	}
	assertErrorMessage(t, "sql: database is closed", lError, "ormshift.Migrate")
}

func Test_MigrationManager_DownLast_ShouldExecuteWithSuccess(t *testing.T) {
	lDB, lError := sql.Open(ormshift.DriverSQLite.Name(), ormshift.DriverSQLite.ConnectionString(ormshift.ConnectionParams{InMemory: true}))
	if !assertNilError(t, lError, "sql.Open") {
		return
	}
	defer lDB.Close()

	lMigrationManager, lError := ormshift.Migrate(
		lDB,
		ormshift.DriverSQLite,
		m001_Create_Table_User{},
		m002_Alter_Table_Usaer_Add_Column_UpdatedAt{},
	)
	if !assertNotNilResultAndNilError(t, lMigrationManager, lError, "ormshift.NewMigrationManager") {
		return
	}

	lUserTableName, lError := ormshift.NewTableName("user")
	if !assertNilError(t, lError, "ormshift.NewTableName") {
		return
	}
	assertEqualWithLabel(t, true, lMigrationManager.DBSchema().ExistsTable(*lUserTableName), "MigrationManager.DBSchema.ExistsTable[user]")

	lError = lMigrationManager.DownLast()
	if !assertNilError(t, lError, "MigrationManager.DownLast") {
		return
	}
	lUpdatedAtColumnName, lError := ormshift.NewColumnName("updated_at")
	if !assertNilError(t, lError, "ormshift.NewColumnName") {
		return
	}
	assertEqualWithLabel(t, false, lMigrationManager.DBSchema().ExistsTableColumn(*lUserTableName, *lUpdatedAtColumnName), "MigrationManager.DBSchema.ExistsTableColumn[user.updated_at]")
}

type m001_Create_Table_User struct{}

func (m m001_Create_Table_User) Up(pMigrationManager *ormshift.MigrationManager) error {
	lUserTable, lError := ormshift.NewTable("user")
	if lError != nil {
		return lError
	}
	if pMigrationManager.DBSchema().ExistsTable(lUserTable.Name()) {
		return nil
	}
	lUserTable.AddColumn(ormshift.NewColumnParams{
		Name:          "id",
		Type:          ormshift.Integer,
		Autoincrement: true,
		PrimaryKey:    true,
		NotNull:       true,
	})
	lUserTable.AddColumn(ormshift.NewColumnParams{
		Name:       "name",
		Type:       ormshift.Varchar,
		Size:       50,
		PrimaryKey: false,
		NotNull:    false,
	})
	lUserTable.AddColumn(ormshift.NewColumnParams{
		Name:       "email",
		Type:       ormshift.Varchar,
		Size:       120,
		PrimaryKey: false,
		NotNull:    false,
	})
	lUserTable.AddColumn(ormshift.NewColumnParams{
		Name:       "active",
		Type:       ormshift.Boolean,
		PrimaryKey: false,
		NotNull:    false,
	})
	lUserTable.AddColumn(ormshift.NewColumnParams{
		Name:       "ammount",
		Type:       ormshift.Monetary,
		PrimaryKey: false,
		NotNull:    false,
	})
	lUserTable.AddColumn(ormshift.NewColumnParams{
		Name:       "percent",
		Type:       ormshift.Decimal,
		PrimaryKey: false,
		NotNull:    false,
	})
	lUserTable.AddColumn(ormshift.NewColumnParams{
		Name:       "photo",
		Type:       ormshift.Binary,
		PrimaryKey: false,
		NotNull:    false,
	})
	_, lError = pMigrationManager.DB().Exec(pMigrationManager.SQLBuilder().CreateTable(*lUserTable))
	if lError != nil {
		return lError
	}
	return nil
}

func (m m001_Create_Table_User) Down(pMigrationManager *ormshift.MigrationManager) error {
	lUserTableName, lError := ormshift.NewTableName("user")
	if lError != nil {
		return lError
	}
	if !pMigrationManager.DBSchema().ExistsTable(*lUserTableName) {
		return nil
	}
	_, lError = pMigrationManager.DB().Exec(pMigrationManager.SQLBuilder().DropTable(*lUserTableName))
	if lError != nil {
		return lError
	}
	return nil
}

type m002_Alter_Table_Usaer_Add_Column_UpdatedAt struct{}

func (m m002_Alter_Table_Usaer_Add_Column_UpdatedAt) Up(pMigrationManager *ormshift.MigrationManager) error {
	lUserTableName, lError := ormshift.NewTableName("user")
	if lError != nil {
		return lError
	}
	lUpdatedAtColumn, lError := ormshift.NewColumn(ormshift.NewColumnParams{
		Name: "updated_at",
		Type: ormshift.DateTime,
	})
	if lError != nil {
		return lError
	}
	if pMigrationManager.DBSchema().ExistsTableColumn(*lUserTableName, lUpdatedAtColumn.Name()) {
		return nil
	}
	_, lError = pMigrationManager.DB().Exec(pMigrationManager.SQLBuilder().AlterTableAddColumn(*lUserTableName, *lUpdatedAtColumn))
	if lError != nil {
		return lError
	}
	return nil
}

func (m m002_Alter_Table_Usaer_Add_Column_UpdatedAt) Down(pMigrationManager *ormshift.MigrationManager) error {
	lUserTableName, lError := ormshift.NewTableName("user")
	if lError != nil {
		return lError
	}
	lUpdatedAtColumnName, lError := ormshift.NewColumnName("updated_at")
	if lError != nil {
		return lError
	}
	if !pMigrationManager.DBSchema().ExistsTableColumn(*lUserTableName, *lUpdatedAtColumnName) {
		return nil
	}
	_, lError = pMigrationManager.DB().Exec(pMigrationManager.SQLBuilder().AlterTableDropColumn(*lUserTableName, *lUpdatedAtColumnName))
	if lError != nil {
		return lError
	}
	return nil
}
