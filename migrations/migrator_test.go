package migrations_test

import (
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/postgresql"
	"github.com/ordershift/ormshift/dialects/sqlite"
	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/migrations"
)

func TestNewMigratorWhenDatabaseIsNil(t *testing.T) {
	migrator, err := migrations.NewMigrator(nil, migrations.NewMigratorConfig())
	testutils.AssertNilResultAndNotNilError(t, migrator, err, "migrations.NewMigrator[database=nil]")
	testutils.AssertErrorMessage(t, "failed to migrate: database cannot be nil", err, "migrations.NewMigrator[database=nil]")
}

func TestNewMigratorWhenConfigIsNil(t *testing.T) {
	db, err := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNotNilResultAndNilError(t, db, err, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = db.Close() }()

	migrator, err := migrations.NewMigrator(db, nil)
	testutils.AssertNilResultAndNotNilError(t, migrator, err, "migrations.NewMigrator[config=nil]")
	testutils.AssertErrorMessage(t, "failed to migrate: migrator config cannot be nil", err, "migrations.NewMigrator[config=nil]")
}

func TestNewMigratorWhenDatabaseIsInvalid(t *testing.T) {
	driver := testutils.NewFakeDriverInvalidConnectionString(postgresql.Driver())
	db, err := ormshift.OpenDatabase(driver, ormshift.ConnectionParams{})
	if !testutils.AssertNotNilResultAndNilError(t, db, err, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = db.Close() }()

	migrator, err := migrations.NewMigrator(db, migrations.NewMigratorConfig())
	testutils.AssertNilResultAndNotNilError(t, migrator, err, "migrations.NewMigrator[database=invalid]")
	testutils.AssertErrorMessage(t, "failed to migrate: failed to get applied migration names: missing \"=\" after \"invalid-connection-string\" in connection info string\"", err, "migrations.NewMigrator[database=invalid]")
}

func TestApplyAllMigrationsFailsWhenRecordingFails(t *testing.T) {
	db, err := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNotNilResultAndNilError(t, db, err, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = db.Close() }()

	migrator, err := migrations.NewMigrator(db, migrations.NewMigratorConfig())
	if !testutils.AssertNotNilResultAndNilError(t, migrator, err, "migrations.NewMigrator") {
		return
	}
	migrator.Add(testutils.M005_Blank_Migration{})

	_ = db.Close()

	err = migrator.ApplyAllMigrations()
	if !testutils.AssertNotNilError(t, err, "Migrator.ApplyAllMigrations") {
		return
	}
	testutils.AssertErrorMessage(t, "failed to record applied migration \"M005_Blank_Migration\": sql: database is closed", err, "Migrator.ApplyAllMigrations")
}

func TestRevertLastAppliedMigration(t *testing.T) {
	db, err := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNotNilResultAndNilError(t, db, err, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = db.Close() }()

	migrator, err := migrations.Migrate(
		db,
		migrations.NewMigratorConfig(),
		testutils.M001_Create_Table_User{},
		testutils.M002_Alter_Table_User_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNotNilResultAndNilError(t, migrator, err, "migrations.NewMigrator") {
		return
	}

	userTableName := "user"
	testutils.AssertEqualWithLabel(t, true, db.DBSchema().HasTable(userTableName), "Migrator.DBSchema.HasTable[user]")

	err = migrator.RevertLastAppliedMigration()
	if !testutils.AssertNilError(t, err, "Migrator.RevertLastAppliedMigration") {
		return
	}

	updatedAtColumnName := "updated_at"
	testutils.AssertEqualWithLabel(t, false, db.DBSchema().HasColumn(userTableName, updatedAtColumnName), "Migrator.DBSchema.HasColumn[user.updated_at]")
}

func TestRevertLastAppliedMigrationWhenNoMigrationsApplied(t *testing.T) {
	db, err := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNotNilResultAndNilError(t, db, err, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = db.Close() }()

	migrator, err := migrations.NewMigrator(db, migrations.NewMigratorConfig())
	if !testutils.AssertNotNilResultAndNilError(t, migrator, err, "migrations.NewMigrator") {
		return
	}
	err = migrator.RevertLastAppliedMigration()
	if !testutils.AssertNilError(t, err, "Migrator.RevertLastAppliedMigration") {
		return
	}
}

func TestRevertLastAppliedMigrationFailsWhenDownFails(t *testing.T) {
	db, err := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNotNilResultAndNilError(t, db, err, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = db.Close() }()

	migrator, err := migrations.Migrate(
		db,
		migrations.NewMigratorConfig(),
		testutils.M004_Bad_Migration_Fails_To_Revert{},
	)
	if !testutils.AssertNotNilResultAndNilError(t, migrator, err, "migrations.Migrate") {
		return
	}

	err = migrator.RevertLastAppliedMigration()
	if !testutils.AssertNotNilError(t, err, "Migrator.RevertLastAppliedMigration") {
		return
	}
	testutils.AssertErrorMessage(t, "failed to revert migration \"M004_Bad_Migration_Fails_To_Revert\": intentionally failed to Down", err, "Migrator.RevertLastAppliedMigration")
}

func TestRevertLastAppliedMigrationFailsWhenDeletingFails(t *testing.T) {
	db, err := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNotNilResultAndNilError(t, db, err, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = db.Close() }()

	migrator, err := migrations.NewMigrator(db, migrations.NewMigratorConfig())
	if !testutils.AssertNotNilResultAndNilError(t, migrator, err, "migrations.NewMigrator") {
		return
	}
	migrator.Add(testutils.M005_Blank_Migration{})
	err = migrator.ApplyAllMigrations()
	if !testutils.AssertNilError(t, err, "Migrator.ApplyAllMigrations") {
		return
	}

	_ = db.Close()

	err = migrator.RevertLastAppliedMigration()
	if !testutils.AssertNotNilError(t, err, "Migrator.RevertLastAppliedMigration") {
		return
	}
	testutils.AssertErrorMessage(t, "failed to delete applied migration \"M005_Blank_Migration\": sql: database is closed", err, "Migrator.RevertLastAppliedMigration")
}
