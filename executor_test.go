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
	db, err := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNotNilResultAndNilError(t, db, err, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = db.Close() }()

	migrationManager, err := migrations.Migrate(
		db,
		migrations.NewMigratorConfig(),
		testutils.M001_Create_Table_User{},
		testutils.M002_Alter_Table_User_Add_Column_UpdatedAt{},
	)
	if !testutils.AssertNotNilResultAndNilError(t, migrationManager, err, "migrations.Migrate") {
		return
	}

	sqlExecutor := db.SQLExecutor()
	sqlBuilder := db.SQLBuilder()

	//INSERT
	values := ormshift.ColumnsValues{
		"name":       "Jonh Doe",
		"email":      "jonh.doe@mail.com",
		"updated_at": time.Date(2026, time.January, 9, 12, 15, 52, 100952002, time.UTC),
		"active":     true,
		"ammount":    5000.00,
		"percent":    25.54325,
	}
	sqlInsert, args := sqlBuilder.InsertWithValues("user", values)
	result, err := sqlExecutor.Exec(sqlInsert, args...)
	if !testutils.AssertNilError(t, err, "sqlExecutor.Exec") {
		return
	}
	r, err := result.RowsAffected()
	if !testutils.AssertNilError(t, err, "RowsAffected") {
		return
	}
	testutils.AssertEqualWithLabel(t, 1, r, "RowsAffected")
	i, err := result.LastInsertId()
	if !testutils.AssertNilError(t, err, "LastInsertId") {
		return
	}

	//SELECT
	sqlSelect, args := sqlBuilder.SelectWithValues(
		"user",
		[]string{"id", "name", "email", "updated_at", "active", "ammount", "percent", "photo"},
		ormshift.ColumnsValues{"id": i},
	)
	userRows, err := sqlExecutor.Query(sqlSelect, args...)
	if !testutils.AssertNotNilResultAndNilError(t, userRows, err, "sqlExecutor.Query") {
		return
	}
	defer func() { _ = userRows.Close() }()
	if !testutils.AssertEqualWithLabel(t, true, userRows.Next(), "Next") {
		return
	}

	//SCAN
	var userRow userRow
	err = userRows.Scan(
		&userRow.ID,
		&userRow.Name,
		&userRow.Email,
		&userRow.UpdatedAt,
		&userRow.Active,
		&userRow.Ammount,
		&userRow.Percent,
		&userRow.Photo,
	)
	if !testutils.AssertNilError(t, err, "Scan") {
		return
	}
	testutils.AssertEqualWithLabel(t, i, userRow.ID.Int64, "user.id")
	testutils.AssertEqualWithLabel(t, "Jonh Doe", userRow.Name.String, "user.name")
	testutils.AssertEqualWithLabel(t, "jonh.doe@mail.com", userRow.Email.String, "user.email")
	testutils.AssertEqualWithLabel(t, time.Date(2026, time.January, 9, 12, 15, 52, 100952002, time.UTC), userRow.UpdatedAt.Time, "user.updated_at")
	testutils.AssertEqualWithLabel(t, 1, userRow.Active.Byte, "user.active")
	testutils.AssertEqualWithLabel(t, 5000.00, userRow.Ammount.Float64, "user.ammount")
	testutils.AssertEqualWithLabel(t, 25.54325, userRow.Percent.Float64, "user.percent")
	testutils.AssertEqualWithLabel(t, nil, userRow.Photo, "user.photo")
}
