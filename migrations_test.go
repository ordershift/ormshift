package ormshift_test

import (
	"database/sql"
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/internal/sqlite"
	"github.com/ordershift/ormshift/internal/testutils"
)

func Test_Migrate_ShouldExecuteWithSuccess(t *testing.T) {
	lDB, lError := sql.Open(sqlite.DriverName(), sqlite.ConnectionString(ormshift.ConnectionParams{InMemory: true}))
	if !testutils.AssertNilError(t, lError, "sql.Open") {
		return
	}
	defer lDB.Close()

	lSQLBuilder := sqlite.SQLBuilder()

	lDBSchema, lError := sqlite.DBSchema(lDB)
	if !testutils.AssertNotNilResultAndNilError(t, lDBSchema, lError, "sqlite.DBSchema") {
		return
	}

	lMigrator, lError := ormshift.Migrate(
		lDB,
		lSQLBuilder,
		lDBSchema,
		ormshift.NewMigratorConfig(),
		m001_Create_Table_User{},
		m002_Alter_Table_Usaer_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNotNilResultAndNilError(t, lMigrator, lError, "ormshift.Migrate") {
		return
	}
	lUserTableName, lError := ormshift.NewTableName("user")
	if !testutils.AssertNilError(t, lError, "ormshift.NewTableName") {
		return
	}
	lUpdatedAtColumnName, lError := ormshift.NewColumnName("updated_at")
	if !testutils.AssertNilError(t, lError, "ormshift.NewColumnName") {
		return
	}
	testutils.AssertEqualWithLabel(t, true, lMigrator.DBSchema().ExistsTableColumn(*lUserTableName, *lUpdatedAtColumnName), "Migrator.DBSchema.ExistsTableColumn[user.updated_at]")
	testutils.AssertEqualWithLabel(t, 2, len(lMigrator.AppliedMigrationNames()), "len(Migrator.AppliedMigrationNames)")
}

func Test_Migrate_ShouldExecuteWithSuccess_WhenTwiceExecute(t *testing.T) {
	lDB, lError := sql.Open(sqlite.DriverName(), sqlite.ConnectionString(ormshift.ConnectionParams{InMemory: true}))
	if !testutils.AssertNilError(t, lError, "sql.Open") {
		return
	}
	defer lDB.Close()

	lSQLBuilder := sqlite.SQLBuilder()

	lDBSchema, lError := sqlite.DBSchema(lDB)
	if !testutils.AssertNotNilResultAndNilError(t, lDBSchema, lError, "sqlite.DBSchema") {
		return
	}

	lMigrator, lError := ormshift.Migrate(
		lDB,
		lSQLBuilder,
		lDBSchema,
		ormshift.NewMigratorConfig(),
		m001_Create_Table_User{},
		m002_Alter_Table_Usaer_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNotNilResultAndNilError(t, lMigrator, lError, "ormshift.Migrate") {
		return
	}

	lMigrator, lError = ormshift.Migrate(
		lDB,
		lSQLBuilder,
		lDBSchema,
		ormshift.NewMigratorConfig(),
		m001_Create_Table_User{},
		m002_Alter_Table_Usaer_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNotNilResultAndNilError(t, lMigrator, lError, "ormshift.Migrate") {
		return
	}

	lUserTableName, lError := ormshift.NewTableName("user")
	if !testutils.AssertNilError(t, lError, "ormshift.NewTableName") {
		return
	}
	lUpdatedAtColumnName, lError := ormshift.NewColumnName("updated_at")
	if !testutils.AssertNilError(t, lError, "ormshift.NewColumnName") {
		return
	}
	testutils.AssertEqualWithLabel(t, true, lMigrator.DBSchema().ExistsTableColumn(*lUserTableName, *lUpdatedAtColumnName), "Migrator.DBSchema.ExistsTableColumn[user.updated_at]")
	testutils.AssertEqualWithLabel(t, 2, len(lMigrator.AppliedMigrationNames()), "len(Migrator.AppliedMigrationNames)")
}

func Test_Migrate_ShouldFail_WhenNilDB(t *testing.T) {
	lMigrator, lError := ormshift.Migrate(
		nil,
		sqlite.SQLBuilder(),
		nil,
		ormshift.NewMigratorConfig(),
		m001_Create_Table_User{},
		m002_Alter_Table_Usaer_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNilResultAndNotNilError(t, lMigrator, lError, "ormshift.Migrate") {
		return
	}
	testutils.AssertErrorMessage(t, "sql.DB cannot be nil", lError, "ormshift.Migrate")
}

func Test_Migrate_ShouldFail_WhenClosedDB(t *testing.T) {
	lDB, lError := sql.Open(sqlite.DriverName(), sqlite.ConnectionString(ormshift.ConnectionParams{InMemory: true}))
	if !testutils.AssertNilError(t, lError, "sql.Open") {
		return
	}

	lSQLBuilder := sqlite.SQLBuilder()

	lDBSchema, lError := sqlite.DBSchema(lDB)
	if !testutils.AssertNotNilResultAndNilError(t, lDBSchema, lError, "sqlite.DBSchema") {
		return
	}

	lDB.Close()

	lMigrator, lError := ormshift.Migrate(
		lDB,
		lSQLBuilder,
		lDBSchema,
		ormshift.NewMigratorConfig(),
		m001_Create_Table_User{},
		m002_Alter_Table_Usaer_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNilResultAndNotNilError(t, lMigrator, lError, "ormshift.Migrate") {
		return
	}
	testutils.AssertErrorMessage(t, "sql: database is closed", lError, "ormshift.Migrate")
}

func Test_Migrator_DownLast_ShouldExecuteWithSuccess(t *testing.T) {
	lDB, lError := sql.Open(sqlite.DriverName(), sqlite.ConnectionString(ormshift.ConnectionParams{InMemory: true}))
	if !testutils.AssertNilError(t, lError, "sql.Open") {
		return
	}
	defer lDB.Close()

	lSQLBuilder := sqlite.SQLBuilder()

	lDBSchema, lError := sqlite.DBSchema(lDB)
	if !testutils.AssertNotNilResultAndNilError(t, lDBSchema, lError, "sqlite.DBSchema") {
		return
	}

	lMigrator, lError := ormshift.Migrate(
		lDB,
		lSQLBuilder,
		lDBSchema,
		ormshift.NewMigratorConfig(),
		m001_Create_Table_User{},
		m002_Alter_Table_Usaer_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNotNilResultAndNilError(t, lMigrator, lError, "ormshift.NewMigrator") {
		return
	}

	lUserTableName, lError := ormshift.NewTableName("user")
	if !testutils.AssertNilError(t, lError, "ormshift.NewTableName") {
		return
	}
	testutils.AssertEqualWithLabel(t, true, lMigrator.DBSchema().ExistsTable(*lUserTableName), "Migrator.DBSchema.ExistsTable[user]")

	lError = lMigrator.RevertLatestMigration()
	if !testutils.AssertNilError(t, lError, "Migrator.DownLast") {
		return
	}
	lUpdatedAtColumnName, lError := ormshift.NewColumnName("updated_at")
	if !testutils.AssertNilError(t, lError, "ormshift.NewColumnName") {
		return
	}
	testutils.AssertEqualWithLabel(t, false, lMigrator.DBSchema().ExistsTableColumn(*lUserTableName, *lUpdatedAtColumnName), "Migrator.DBSchema.ExistsTableColumn[user.updated_at]")
}

type m001_Create_Table_User struct{}

func (m m001_Create_Table_User) Up(pMigrator *ormshift.Migrator) error {
	lUserTable, lError := ormshift.NewTable("user")
	if lError != nil {
		return lError
	}
	if pMigrator.DBSchema().ExistsTable(lUserTable.Name()) {
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
	_, lError = pMigrator.DB().Exec(pMigrator.SQLBuilder().CreateTable(*lUserTable))
	if lError != nil {
		return lError
	}
	return nil
}

func (m m001_Create_Table_User) Down(pMigrator *ormshift.Migrator) error {
	lUserTableName, lError := ormshift.NewTableName("user")
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

type m002_Alter_Table_Usaer_Add_Column_UpdatedAt struct{}

func (m m002_Alter_Table_Usaer_Add_Column_UpdatedAt) Up(pMigrator *ormshift.Migrator) error {
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
	if pMigrator.DBSchema().ExistsTableColumn(*lUserTableName, lUpdatedAtColumn.Name()) {
		return nil
	}
	_, lError = pMigrator.DB().Exec(pMigrator.SQLBuilder().AlterTableAddColumn(*lUserTableName, *lUpdatedAtColumn))
	if lError != nil {
		return lError
	}
	return nil
}

func (m m002_Alter_Table_Usaer_Add_Column_UpdatedAt) Down(pMigrator *ormshift.Migrator) error {
	lUserTableName, lError := ormshift.NewTableName("user")
	if lError != nil {
		return lError
	}
	lUpdatedAtColumnName, lError := ormshift.NewColumnName("updated_at")
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
