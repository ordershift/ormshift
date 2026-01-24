package migrations_test

import (
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/postgresql"
	"github.com/ordershift/ormshift/dialects/sqlite"
	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/migrations"
	"github.com/ordershift/ormshift/schema"
)

func TestNewMigratorWhenDatabaseIsNil(t *testing.T) {
	lMigrator, lError := migrations.NewMigrator(nil, migrations.NewMigratorConfig())
	testutils.AssertNilResultAndNotNilError(t, lMigrator, lError, "migrations.NewMigrator[database=nil]")
	testutils.AssertErrorMessage(t, "database cannot be nil", lError, "migrations.NewMigrator[database=nil]")
}

func TestNewMigratorWhenDatabaseIsInvalid(t *testing.T) {
	lDriver := testutils.NewFakeDriverInvalidConnectionString(postgresql.Driver())
	lDB, lError := ormshift.OpenDatabase(lDriver, ormshift.ConnectionParams{})
	if !testutils.AssertNotNilResultAndNilError(t, lDB, lError, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = lDB.Close() }()

	lMigrator, lError := migrations.NewMigrator(lDB, migrations.NewMigratorConfig())
	testutils.AssertNilResultAndNotNilError(t, lMigrator, lError, "migrations.NewMigrator[database=invalid]")
	testutils.AssertErrorMessage(t, "failed to get applied migration names: missing \"=\" after \"invalid-connection-string\" in connection info string\"", lError, "migrations.NewMigrator[database=invalid]")
}

func TestApplyAllMigrationsFailsWhenRecordingFails(t *testing.T) {
	lDB, lError := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if lError != nil {
		t.Errorf("ormshift.OpenDatabase failed: %v", lError)
		return
	}

	lMigrator, lError := migrations.NewMigrator(lDB, migrations.NewMigratorConfig())
	if !testutils.AssertNotNilResultAndNilError(t, lMigrator, lError, "migrations.NewMigrator") {
		return
	}
	lMigrator.Add(testutils.M005_Blank_Migration{})

	_ = lDB.Close()

	lError = lMigrator.ApplyAllMigrations()
	if !testutils.AssertNotNilError(t, lError, "Migrator.ApplyAllMigrations") {
		return
	}
	testutils.AssertErrorMessage(t, "failed to record applied migration \"M005_Blank_Migration\": sql: database is closed", lError, "Migrator.ApplyAllMigrations")
}

func TestRevertLastAppliedMigration(t *testing.T) {
	lDB, lError := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if lError != nil {
		t.Errorf("ormshift.OpenDatabase failed: %v", lError)
		return
	}
	defer func() { _ = lDB.Close() }()

	lMigrator, lError := migrations.Migrate(
		lDB,
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
	testutils.AssertEqualWithLabel(t, true, lDB.DBSchema().ExistsTable(*lUserTableName), "Migrator.DBSchema.ExistsTable[user]")

	lError = lMigrator.RevertLastAppliedMigration()
	if !testutils.AssertNilError(t, lError, "Migrator.RevertLastAppliedMigration") {
		return
	}
	lUpdatedAtColumnName, lError := schema.NewColumnName("updated_at")
	if !testutils.AssertNilError(t, lError, "migrations.NewColumnName") {
		return
	}
	testutils.AssertEqualWithLabel(t, false, lDB.DBSchema().ExistsTableColumn(*lUserTableName, *lUpdatedAtColumnName), "Migrator.DBSchema.ExistsTableColumn[user.updated_at]")
}

func TestRevertLastAppliedMigrationFailsWhenDownFails(t *testing.T) {
	lDB, lError := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if lError != nil {
		t.Errorf("ormshift.OpenDatabase failed: %v", lError)
		return
	}
	defer func() { _ = lDB.Close() }()

	lMigrator, lError := migrations.Migrate(
		lDB,
		migrations.NewMigratorConfig(),
		testutils.M004_Bad_Migration_Fails_To_Revert{},
	)
	if !testutils.AssertNotNilResultAndNilError(t, lMigrator, lError, "migrations.Migrate") {
		return
	}

	lError = lMigrator.RevertLastAppliedMigration()
	if !testutils.AssertNotNilError(t, lError, "Migrator.RevertLastAppliedMigration") {
		return
	}
	testutils.AssertErrorMessage(t, "failed to revert migration \"M004_Bad_Migration_Fails_To_Revert\": intentionally failed to Down", lError, "Migrator.RevertLastAppliedMigration")
}

func TestRevertLastAppliedMigrationFailsWhenDeletingFails(t *testing.T) {
	lDB, lError := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if lError != nil {
		t.Errorf("ormshift.OpenDatabase failed: %v", lError)
		return
	}
	defer func() { _ = lDB.Close() }()

	lMigrator, lError := migrations.NewMigrator(lDB, migrations.NewMigratorConfig())
	if !testutils.AssertNotNilResultAndNilError(t, lMigrator, lError, "migrations.NewMigrator") {
		return
	}
	lMigrator.Add(testutils.M005_Blank_Migration{})
	lError = lMigrator.ApplyAllMigrations()
	if !testutils.AssertNilError(t, lError, "Migrator.ApplyAllMigrations") {
		return
	}

	_ = lDB.Close()

	lError = lMigrator.RevertLastAppliedMigration()
	if !testutils.AssertNotNilError(t, lError, "Migrator.RevertLastAppliedMigration") {
		return
	}
	testutils.AssertErrorMessage(t, "failed to delete applied migration \"M005_Blank_Migration\": sql: database is closed", lError, "Migrator.RevertLastAppliedMigration")
}
