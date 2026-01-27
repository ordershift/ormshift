package ormshift_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/sqlite"
	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/migrations"
)

type userRow struct {
	ID        sql.NullInt64
	Name      sql.NullString
	Email     sql.NullString
	UpdatedAt sql.NullTime
	Active    sql.NullByte
	Ammount   sql.NullFloat64
	Percent   sql.NullFloat64
	Photo     *[]byte
}

func TestExecutor(t *testing.T) {
	lDB, lError := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNotNilResultAndNilError(t, lDB, lError, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = lDB.Close() }()

	lMigrationManager, lError := migrations.Migrate(
		lDB,
		migrations.NewMigratorConfig(),
		testutils.M001_Create_Table_User{},
		testutils.M002_Alter_Table_User_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNotNilResultAndNilError(t, lMigrationManager, lError, "migrations.Migrate") {
		return
	}

	lSQLExecutor := lDB.SQLExecutor()
	lSQLBuilder := lDB.SQLBuilder()

	//INSERT
	lValues := ormshift.ColumnsValues{
		"name":       "Jonh Doe",
		"email":      "jonh.doe@mail.com",
		"updated_at": time.Date(2026, time.January, 9, 12, 15, 52, 100952002, time.UTC),
		"active":     true,
		"ammount":    5000.00,
		"percent":    25.54325,
	}
	lSQLInsert, lArgs := lSQLBuilder.InsertWithValues("user", lValues)
	lResult, lError := lSQLExecutor.Exec(lSQLInsert, lArgs...)
	if !testutils.AssertNilError(t, lError, "sqlExecutor.Exec") {
		return
	}
	r, lError := lResult.RowsAffected()
	if !testutils.AssertNilError(t, lError, "RowsAffected") {
		return
	}
	testutils.AssertEqualWithLabel(t, 1, r, "RowsAffected")
	i, lError := lResult.LastInsertId()
	if !testutils.AssertNilError(t, lError, "LastInsertId") {
		return
	}

	//SELECT
	lSQLSelect, lArgs := lSQLBuilder.SelectWithValues(
		"user",
		[]string{"id", "name", "email", "updated_at", "active", "ammount", "percent", "photo"},
		ormshift.ColumnsValues{"id": i},
	)
	lUserRows, lError := lSQLExecutor.Query(lSQLSelect, lArgs...)
	if !testutils.AssertNotNilResultAndNilError(t, lUserRows, lError, "sqlExecutor.Query") {
		return
	}
	defer func() { _ = lUserRows.Close() }()
	if !testutils.AssertEqualWithLabel(t, true, lUserRows.Next(), "Next") {
		return
	}

	//SCAN
	var lUserRow userRow
	lError = lUserRows.Scan(
		&lUserRow.ID,
		&lUserRow.Name,
		&lUserRow.Email,
		&lUserRow.UpdatedAt,
		&lUserRow.Active,
		&lUserRow.Ammount,
		&lUserRow.Percent,
		&lUserRow.Photo,
	)
	if !testutils.AssertNilError(t, lError, "Scan") {
		return
	}
	testutils.AssertEqualWithLabel(t, i, lUserRow.ID.Int64, "user.id")
	testutils.AssertEqualWithLabel(t, "Jonh Doe", lUserRow.Name.String, "user.name")
	testutils.AssertEqualWithLabel(t, "jonh.doe@mail.com", lUserRow.Email.String, "user.email")
	testutils.AssertEqualWithLabel(t, time.Date(2026, time.January, 9, 12, 15, 52, 100952002, time.UTC), lUserRow.UpdatedAt.Time, "user.updated_at")
	testutils.AssertEqualWithLabel(t, 1, lUserRow.Active.Byte, "user.active")
	testutils.AssertEqualWithLabel(t, 5000.00, lUserRow.Ammount.Float64, "user.ammount")
	testutils.AssertEqualWithLabel(t, 25.54325, lUserRow.Percent.Float64, "user.percent")
	testutils.AssertEqualWithLabel(t, nil, lUserRow.Photo, "user.photo")
}
