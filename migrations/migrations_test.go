package migrations_test

import (
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/sqlite"
	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/migrations"
	"github.com/ordershift/ormshift/schema"
)

func Test_Migrate_ShouldExecuteWithSuccess(t *testing.T) {
	lDatabase, lError := ormshift.OpenDatabase(sqlite.SQLiteDriver{}, ormshift.ConnectionParams{InMemory: true})
	if lError != nil {
		t.Errorf("ormshift.OpenDatabase failed: %v", lError)
		return
	}
	defer func() { _ = lDatabase.Close() }()

	lMigrator, lError := migrations.Migrate(
		*lDatabase,
		migrations.NewMigratorConfig(),
		testutils.M001_Create_Table_User{},
		testutils.M002_Alter_Table_User_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNotNilResultAndNilError(t, lMigrator, lError, "migrations.Migrate") {
		return
	}
	lUserTableName, lError := schema.NewTableName("user")
	if !testutils.AssertNilError(t, lError, "migrations.NewTableName") {
		return
	}
	lUpdatedAtColumnName, lError := schema.NewColumnName("updated_at")
	if !testutils.AssertNilError(t, lError, "migrations.NewColumnName") {
		return
	}
	testutils.AssertEqualWithLabel(t, true, lMigrator.DBSchema().ExistsTableColumn(*lUserTableName, *lUpdatedAtColumnName), "Migrator.DBSchema.ExistsTableColumn[user.updated_at]")
	testutils.AssertEqualWithLabel(t, 2, len(lMigrator.AppliedMigrationNames()), "len(Migrator.AppliedMigrationNames)")
}

func Test_Migrate_ShouldExecuteWithSuccess_WhenTwiceExecute(t *testing.T) {
	lDatabase, lError := ormshift.OpenDatabase(sqlite.SQLiteDriver{}, ormshift.ConnectionParams{InMemory: true})
	if lError != nil {
		t.Errorf("ormshift.OpenDatabase failed: %v", lError)
		return
	}
	defer func() { _ = lDatabase.Close() }()

	lMigrator, lError := migrations.Migrate(
		*lDatabase,
		migrations.NewMigratorConfig(),
		testutils.M001_Create_Table_User{},
		testutils.M002_Alter_Table_User_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNotNilResultAndNilError(t, lMigrator, lError, "migrations.Migrate") {
		return
	}

	lMigrator, lError = migrations.Migrate(
		*lDatabase,
		migrations.NewMigratorConfig(),
		testutils.M001_Create_Table_User{},
		testutils.M002_Alter_Table_User_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNotNilResultAndNilError(t, lMigrator, lError, "migrations.Migrate") {
		return
	}

	lUserTableName, lError := schema.NewTableName("user")
	if !testutils.AssertNilError(t, lError, "migrations.NewTableName") {
		return
	}
	lUpdatedAtColumnName, lError := schema.NewColumnName("updated_at")
	if !testutils.AssertNilError(t, lError, "migrations.NewColumnName") {
		return
	}
	testutils.AssertEqualWithLabel(t, true, lMigrator.DBSchema().ExistsTableColumn(*lUserTableName, *lUpdatedAtColumnName), "Migrator.DBSchema.ExistsTableColumn[user.updated_at]")
	testutils.AssertEqualWithLabel(t, 2, len(lMigrator.AppliedMigrationNames()), "len(Migrator.AppliedMigrationNames)")
}

func Test_Migrate_ShouldFail_WhenDatabaseIsInvalid(t *testing.T) {
	lDatabase, lError := ormshift.OpenDatabase(sqlite.SQLiteDriver{}, ormshift.ConnectionParams{InMemory: true})
	if lError != nil {
		t.Errorf("ormshift.OpenDatabase failed: %v", lError)
		return
	}
	_ = lDatabase.Close()
	lMigrator, lError := migrations.Migrate(
		*lDatabase,
		migrations.NewMigratorConfig(),
		testutils.M001_Create_Table_User{},
		testutils.M002_Alter_Table_User_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNilResultAndNotNilError(t, lMigrator, lError, "migrations.Migrate") {
		return
	}
	testutils.AssertErrorMessage(t, "invalid database: database ping failed: sql: database is closed", lError, "migrations.Migrate")
}

func Test_Migrator_DownLast_ShouldExecuteWithSuccess(t *testing.T) {
	lDatabase, lError := ormshift.OpenDatabase(sqlite.SQLiteDriver{}, ormshift.ConnectionParams{InMemory: true})
	if lError != nil {
		t.Errorf("ormshift.OpenDatabase failed: %v", lError)
		return
	}
	defer func() { _ = lDatabase.Close() }()

	lMigrator, lError := migrations.Migrate(
		*lDatabase,
		migrations.NewMigratorConfig(),
		testutils.M001_Create_Table_User{},
		testutils.M002_Alter_Table_User_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNotNilResultAndNilError(t, lMigrator, lError, "migrations.NewMigrator") {
		return
	}

	lUserTableName, lError := schema.NewTableName("user")
	if !testutils.AssertNilError(t, lError, "migrations.NewTableName") {
		return
	}
	testutils.AssertEqualWithLabel(t, true, lMigrator.DBSchema().ExistsTable(*lUserTableName), "Migrator.DBSchema.ExistsTable[user]")

	lError = lMigrator.RevertLatestMigration()
	if !testutils.AssertNilError(t, lError, "Migrator.DownLast") {
		return
	}
	lUpdatedAtColumnName, lError := schema.NewColumnName("updated_at")
	if !testutils.AssertNilError(t, lError, "migrations.NewColumnName") {
		return
	}
	testutils.AssertEqualWithLabel(t, false, lMigrator.DBSchema().ExistsTableColumn(*lUserTableName, *lUpdatedAtColumnName), "Migrator.DBSchema.ExistsTableColumn[user.updated_at]")
}
