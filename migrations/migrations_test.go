package migrations_test

import (
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/sqlite"
	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/migrations"
	"github.com/ordershift/ormshift/schema"
)

func TestMigrate(t *testing.T) {
	lDatabase, lError := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if lError != nil {
		t.Errorf("ormshift.OpenDatabase failed: %v", lError)
		return
	}
	defer func() { _ = lDatabase.Close() }()

	lMigrator, lError := migrations.Migrate(
		lDatabase,
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
	testutils.AssertEqualWithLabel(t, true, lDatabase.DBSchema().ExistsTableColumn(*lUserTableName, *lUpdatedAtColumnName), "Migrator.DBSchema.ExistsTableColumn[user.updated_at]")
	testutils.AssertEqualWithLabel(t, 2, len(lMigrator.AppliedMigrations()), "len(Migrator.AppliedMigrationNames)")
}

func TestMigrateTwice(t *testing.T) {
	lDatabase, lError := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if lError != nil {
		t.Errorf("ormshift.OpenDatabase failed: %v", lError)
		return
	}
	defer func() { _ = lDatabase.Close() }()

	lMigrator, lError := migrations.Migrate(
		lDatabase,
		migrations.NewMigratorConfig(),
		testutils.M001_Create_Table_User{},
		testutils.M002_Alter_Table_User_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNotNilResultAndNilError(t, lMigrator, lError, "migrations.Migrate") {
		return
	}

	lMigrator, lError = migrations.Migrate(
		lDatabase,
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
	testutils.AssertEqualWithLabel(t, true, lDatabase.DBSchema().ExistsTableColumn(*lUserTableName, *lUpdatedAtColumnName), "Migrator.DBSchema.ExistsTableColumn[user.updated_at]")
	testutils.AssertEqualWithLabel(t, 2, len(lMigrator.AppliedMigrations()), "len(Migrator.AppliedMigrations)")
}

func TestMigrateFailsWhenDatabaseIsClosed(t *testing.T) {
	lDatabase, lError := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if lError != nil {
		t.Errorf("ormshift.OpenDatabase failed: %v", lError)
		return
	}
	_ = lDatabase.Close()

	lMigrator, lError := migrations.Migrate(
		lDatabase,
		migrations.NewMigratorConfig(),
		testutils.M001_Create_Table_User{},
		testutils.M002_Alter_Table_User_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNilResultAndNotNilError(t, lMigrator, lError, "migrations.Migrate") {
		return
	}
	testutils.AssertErrorMessage(t, "failed to get applied migration names: sql: database is closed", lError, "migrations.Migrate")
}

func TestMigrateFailsWhenMigrationUpFails(t *testing.T) {
	lDatabase, lError := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if lError != nil {
		t.Errorf("ormshift.OpenDatabase failed: %v", lError)
		return
	}
	defer func() { _ = lDatabase.Close() }()

	lMigrator, lError := migrations.Migrate(
		lDatabase,
		migrations.NewMigratorConfig(),
		testutils.M003_Bad_Migration_Fails_To_Apply{},
	)
	if !testutils.AssertNilResultAndNotNilError(t, lMigrator, lError, "migrations.Migrate") {
		return
	}
	testutils.AssertErrorMessage(t, "failed to apply migration \"M003_Bad_Migration_Fails_To_Apply\": intentionally failed to Up", lError, "migrations.Migrate")
}
