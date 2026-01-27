package migrations_test

import (
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/sqlite"
	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/migrations"
)

func TestMigrate(t *testing.T) {
	lDB, lError := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNilError(t, lError, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = lDB.Close() }()

	lMigrator, lError := migrations.Migrate(
		lDB,
		migrations.NewMigratorConfig(),
		testutils.M001_Create_Table_User{},
		testutils.M002_Alter_Table_User_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNotNilResultAndNilError(t, lMigrator, lError, "migrations.Migrate") {
		return
	}
	lUserTableName := "user"
	lUpdatedAtColumnName := "updated_at"
	testutils.AssertEqualWithLabel(t, true, lDB.DBSchema().HasColumn(lUserTableName, lUpdatedAtColumnName), "Migrator.DBSchema.HasColumn[user.updated_at]")
	testutils.AssertEqualWithLabel(t, 2, len(lMigrator.AppliedMigrations()), "len(Migrator.AppliedMigrations)")
}

func TestMigrateTwice(t *testing.T) {
	lDB, lError := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNilError(t, lError, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = lDB.Close() }()

	lMigrator, lError := migrations.Migrate(
		lDB,
		migrations.NewMigratorConfig(),
		testutils.M001_Create_Table_User{},
		testutils.M002_Alter_Table_User_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNotNilResultAndNilError(t, lMigrator, lError, "migrations.Migrate") {
		return
	}

	lMigrator, lError = migrations.Migrate(
		lDB,
		migrations.NewMigratorConfig(),
		testutils.M001_Create_Table_User{},
		testutils.M002_Alter_Table_User_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNotNilResultAndNilError(t, lMigrator, lError, "migrations.Migrate") {
		return
	}

	lUserTableName := "user"
	lUpdatedAtColumnName := "updated_at"
	testutils.AssertEqualWithLabel(t, true, lDB.DBSchema().HasColumn(lUserTableName, lUpdatedAtColumnName), "Migrator.DBSchema.HasColumn[user.updated_at]")
	testutils.AssertEqualWithLabel(t, 2, len(lMigrator.AppliedMigrations()), "len(Migrator.AppliedMigrations)")
}

func TestMigrateFailsWhenDatabaseIsClosed(t *testing.T) {
	lDB, lError := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNilError(t, lError, "ormshift.OpenDatabase") {
		return
	}
	_ = lDB.Close()

	lMigrator, lError := migrations.Migrate(
		lDB,
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
	lDB, lError := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNilError(t, lError, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = lDB.Close() }()

	lMigrator, lError := migrations.Migrate(
		lDB,
		migrations.NewMigratorConfig(),
		testutils.M003_Bad_Migration_Fails_To_Apply{},
	)
	if !testutils.AssertNilResultAndNotNilError(t, lMigrator, lError, "migrations.Migrate") {
		return
	}
	testutils.AssertErrorMessage(t, "failed to apply migration \"M003_Bad_Migration_Fails_To_Apply\": intentionally failed to Up", lError, "migrations.Migrate")
}
