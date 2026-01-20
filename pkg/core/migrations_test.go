package core_test

import (
	"database/sql"
	"testing"

	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/pkg/core"
	"github.com/ordershift/ormshift/pkg/dialects/sqlite"
)

func Test_Migrate_ShouldExecuteWithSuccess(t *testing.T) {
	lDB, lError := sql.Open(sqlite.DriverName(), sqlite.ConnectionString(core.ConnectionParams{InMemory: true}))
	if !testutils.AssertNilError(t, lError, "sql.Open") {
		return
	}
	defer lDB.Close()

	lSQLBuilder := sqlite.SQLBuilder()

	lDBSchema, lError := sqlite.DBSchema(lDB)
	if !testutils.AssertNotNilResultAndNilError(t, lDBSchema, lError, "sqlite.DBSchema") {
		return
	}

	lMigrator, lError := core.Migrate(
		lDB,
		lSQLBuilder,
		lDBSchema,
		core.NewMigratorConfig(),
		m001_Create_Table_User{},
		m002_Alter_Table_Usaer_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNotNilResultAndNilError(t, lMigrator, lError, "core.Migrate") {
		return
	}
	lUserTableName, lError := core.NewTableName("user")
	if !testutils.AssertNilError(t, lError, "core.NewTableName") {
		return
	}
	lUpdatedAtColumnName, lError := core.NewColumnName("updated_at")
	if !testutils.AssertNilError(t, lError, "core.NewColumnName") {
		return
	}
	testutils.AssertEqualWithLabel(t, true, lMigrator.DBSchema().ExistsTableColumn(*lUserTableName, *lUpdatedAtColumnName), "Migrator.DBSchema.ExistsTableColumn[user.updated_at]")
	testutils.AssertEqualWithLabel(t, 2, len(lMigrator.AppliedMigrationNames()), "len(Migrator.AppliedMigrationNames)")
}

func Test_Migrate_ShouldExecuteWithSuccess_WhenTwiceExecute(t *testing.T) {
	lDB, lError := sql.Open(sqlite.DriverName(), sqlite.ConnectionString(core.ConnectionParams{InMemory: true}))
	if !testutils.AssertNilError(t, lError, "sql.Open") {
		return
	}
	defer lDB.Close()

	lSQLBuilder := sqlite.SQLBuilder()

	lDBSchema, lError := sqlite.DBSchema(lDB)
	if !testutils.AssertNotNilResultAndNilError(t, lDBSchema, lError, "sqlite.DBSchema") {
		return
	}

	lMigrator, lError := core.Migrate(
		lDB,
		lSQLBuilder,
		lDBSchema,
		core.NewMigratorConfig(),
		m001_Create_Table_User{},
		m002_Alter_Table_Usaer_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNotNilResultAndNilError(t, lMigrator, lError, "core.Migrate") {
		return
	}

	lMigrator, lError = core.Migrate(
		lDB,
		lSQLBuilder,
		lDBSchema,
		core.NewMigratorConfig(),
		m001_Create_Table_User{},
		m002_Alter_Table_Usaer_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNotNilResultAndNilError(t, lMigrator, lError, "core.Migrate") {
		return
	}

	lUserTableName, lError := core.NewTableName("user")
	if !testutils.AssertNilError(t, lError, "core.NewTableName") {
		return
	}
	lUpdatedAtColumnName, lError := core.NewColumnName("updated_at")
	if !testutils.AssertNilError(t, lError, "core.NewColumnName") {
		return
	}
	testutils.AssertEqualWithLabel(t, true, lMigrator.DBSchema().ExistsTableColumn(*lUserTableName, *lUpdatedAtColumnName), "Migrator.DBSchema.ExistsTableColumn[user.updated_at]")
	testutils.AssertEqualWithLabel(t, 2, len(lMigrator.AppliedMigrationNames()), "len(Migrator.AppliedMigrationNames)")
}

func Test_Migrate_ShouldFail_WhenNilDB(t *testing.T) {
	lMigrator, lError := core.Migrate(
		nil,
		sqlite.SQLBuilder(),
		nil,
		core.NewMigratorConfig(),
		m001_Create_Table_User{},
		m002_Alter_Table_Usaer_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNilResultAndNotNilError(t, lMigrator, lError, "core.Migrate") {
		return
	}
	testutils.AssertErrorMessage(t, "sql.DB cannot be nil", lError, "core.Migrate")
}

func Test_Migrate_ShouldFail_WhenClosedDB(t *testing.T) {
	lDB, lError := sql.Open(sqlite.DriverName(), sqlite.ConnectionString(core.ConnectionParams{InMemory: true}))
	if !testutils.AssertNilError(t, lError, "sql.Open") {
		return
	}

	lSQLBuilder := sqlite.SQLBuilder()

	lDBSchema, lError := sqlite.DBSchema(lDB)
	if !testutils.AssertNotNilResultAndNilError(t, lDBSchema, lError, "sqlite.DBSchema") {
		return
	}

	lDB.Close()

	lMigrator, lError := core.Migrate(
		lDB,
		lSQLBuilder,
		lDBSchema,
		core.NewMigratorConfig(),
		m001_Create_Table_User{},
		m002_Alter_Table_Usaer_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNilResultAndNotNilError(t, lMigrator, lError, "core.Migrate") {
		return
	}
	testutils.AssertErrorMessage(t, "sql: database is closed", lError, "core.Migrate")
}

func Test_Migrator_DownLast_ShouldExecuteWithSuccess(t *testing.T) {
	lDB, lError := sql.Open(sqlite.DriverName(), sqlite.ConnectionString(core.ConnectionParams{InMemory: true}))
	if !testutils.AssertNilError(t, lError, "sql.Open") {
		return
	}
	defer lDB.Close()

	lSQLBuilder := sqlite.SQLBuilder()

	lDBSchema, lError := sqlite.DBSchema(lDB)
	if !testutils.AssertNotNilResultAndNilError(t, lDBSchema, lError, "sqlite.DBSchema") {
		return
	}

	lMigrator, lError := core.Migrate(
		lDB,
		lSQLBuilder,
		lDBSchema,
		core.NewMigratorConfig(),
		m001_Create_Table_User{},
		m002_Alter_Table_Usaer_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNotNilResultAndNilError(t, lMigrator, lError, "core.NewMigrator") {
		return
	}

	lUserTableName, lError := core.NewTableName("user")
	if !testutils.AssertNilError(t, lError, "core.NewTableName") {
		return
	}
	testutils.AssertEqualWithLabel(t, true, lMigrator.DBSchema().ExistsTable(*lUserTableName), "Migrator.DBSchema.ExistsTable[user]")

	lError = lMigrator.RevertLatestMigration()
	if !testutils.AssertNilError(t, lError, "Migrator.DownLast") {
		return
	}
	lUpdatedAtColumnName, lError := core.NewColumnName("updated_at")
	if !testutils.AssertNilError(t, lError, "core.NewColumnName") {
		return
	}
	testutils.AssertEqualWithLabel(t, false, lMigrator.DBSchema().ExistsTableColumn(*lUserTableName, *lUpdatedAtColumnName), "Migrator.DBSchema.ExistsTableColumn[user.updated_at]")
}

type m001_Create_Table_User struct{}

func (m m001_Create_Table_User) Up(pMigrator *core.Migrator) error {
	lUserTable, lError := core.NewTable("user")
	if lError != nil {
		return lError
	}
	if pMigrator.DBSchema().ExistsTable(lUserTable.Name()) {
		return nil
	}
	lUserTable.AddColumn(core.NewColumnParams{
		Name:          "id",
		Type:          core.Integer,
		Autoincrement: true,
		PrimaryKey:    true,
		NotNull:       true,
	})
	lUserTable.AddColumn(core.NewColumnParams{
		Name:       "name",
		Type:       core.Varchar,
		Size:       50,
		PrimaryKey: false,
		NotNull:    false,
	})
	lUserTable.AddColumn(core.NewColumnParams{
		Name:       "email",
		Type:       core.Varchar,
		Size:       120,
		PrimaryKey: false,
		NotNull:    false,
	})
	lUserTable.AddColumn(core.NewColumnParams{
		Name:       "active",
		Type:       core.Boolean,
		PrimaryKey: false,
		NotNull:    false,
	})
	lUserTable.AddColumn(core.NewColumnParams{
		Name:       "ammount",
		Type:       core.Monetary,
		PrimaryKey: false,
		NotNull:    false,
	})
	lUserTable.AddColumn(core.NewColumnParams{
		Name:       "percent",
		Type:       core.Decimal,
		PrimaryKey: false,
		NotNull:    false,
	})
	lUserTable.AddColumn(core.NewColumnParams{
		Name:       "photo",
		Type:       core.Binary,
		PrimaryKey: false,
		NotNull:    false,
	})
	_, lError = pMigrator.DB().Exec(pMigrator.SQLBuilder().CreateTable(*lUserTable))
	if lError != nil {
		return lError
	}
	return nil
}

func (m m001_Create_Table_User) Down(pMigrator *core.Migrator) error {
	lUserTableName, lError := core.NewTableName("user")
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

func (m m002_Alter_Table_Usaer_Add_Column_UpdatedAt) Up(pMigrator *core.Migrator) error {
	lUserTableName, lError := core.NewTableName("user")
	if lError != nil {
		return lError
	}
	lUpdatedAtColumn, lError := core.NewColumn(core.NewColumnParams{
		Name: "updated_at",
		Type: core.DateTime,
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

func (m m002_Alter_Table_Usaer_Add_Column_UpdatedAt) Down(pMigrator *core.Migrator) error {
	lUserTableName, lError := core.NewTableName("user")
	if lError != nil {
		return lError
	}
	lUpdatedAtColumnName, lError := core.NewColumnName("updated_at")
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
