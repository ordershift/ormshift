package migrations_test

import (
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/sqlite"
	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/migrations"
)

func TestMigrate(t *testing.T) {
	db, err := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNilError(t, err, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = db.Close() }()

	migrator, err := migrations.Migrate(
		db,
		migrations.NewMigratorConfig(),
		testutils.M001_Create_Table_User{},
		testutils.M002_Alter_Table_User_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNotNilResultAndNilError(t, migrator, err, "migrations.Migrate") {
		return
	}
	userTableName := "user"
	updatedAtColumnName := "updated_at"
	testutils.AssertEqualWithLabel(t, true, db.DBSchema().HasColumn(userTableName, updatedAtColumnName), "Migrator.DBSchema.HasColumn[user.updated_at]")
	testutils.AssertEqualWithLabel(t, 2, len(migrator.AppliedMigrations()), "len(Migrator.AppliedMigrations)")
}

func TestMigrateTwice(t *testing.T) {
	db, err := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNilError(t, err, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = db.Close() }()

	migrator, err := migrations.Migrate(
		db,
		migrations.NewMigratorConfig(),
		testutils.M001_Create_Table_User{},
		testutils.M002_Alter_Table_User_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNotNilResultAndNilError(t, migrator, err, "migrations.Migrate") {
		return
	}

	migrator, err = migrations.Migrate(
		db,
		migrations.NewMigratorConfig(),
		testutils.M001_Create_Table_User{},
		testutils.M002_Alter_Table_User_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNotNilResultAndNilError(t, migrator, err, "migrations.Migrate") {
		return
	}

	userTableName := "user"
	updatedAtColumnName := "updated_at"
	testutils.AssertEqualWithLabel(t, true, db.DBSchema().HasColumn(userTableName, updatedAtColumnName), "Migrator.DBSchema.HasColumn[user.updated_at]")
	testutils.AssertEqualWithLabel(t, 2, len(migrator.AppliedMigrations()), "len(Migrator.AppliedMigrations)")
}

func TestMigrateFailsWhenDatabaseIsClosed(t *testing.T) {
	db, err := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNilError(t, err, "ormshift.OpenDatabase") {
		return
	}
	_ = db.Close()

	migrator, err := migrations.Migrate(
		db,
		migrations.NewMigratorConfig(),
		testutils.M001_Create_Table_User{},
		testutils.M002_Alter_Table_User_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNilResultAndNotNilError(t, migrator, err, "migrations.Migrate") {
		return
	}
	testutils.AssertErrorMessage(t, "failed to get applied migration names: sql: database is closed", err, "migrations.Migrate")
}

func TestMigrateFailsWhenMigrationUpFails(t *testing.T) {
	db, err := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNilError(t, err, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = db.Close() }()

	migrator, err := migrations.Migrate(
		db,
		migrations.NewMigratorConfig(),
		testutils.M003_Bad_Migration_Fails_To_Apply{},
	)
	if !testutils.AssertNilResultAndNotNilError(t, migrator, err, "migrations.Migrate") {
		return
	}
	testutils.AssertErrorMessage(t, "failed to apply migration \"M003_Bad_Migration_Fails_To_Apply\": intentionally failed to Up", err, "migrations.Migrate")
}
