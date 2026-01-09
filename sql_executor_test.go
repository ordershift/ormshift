package ormshift_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/ordershift/ormshift"
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

func Test_DBExecQuery_MigrateInsertSelectScan_ShouldSuccess(t *testing.T) {
	var lDriverDB ormshift.DriverDB = ormshift.DriverSQLite
	var lConnectionParams ormshift.ConnectionParams = ormshift.ConnectionParams{InMemory: true}
	var lConnectionString string = lDriverDB.ConnectionString(lConnectionParams)

	//MIGRATE
	lDB, lError := sql.Open(lDriverDB.Name(), lConnectionString)
	if !assertNilError(t, lError, "sql.Open") {
		return
	}
	defer lDB.Close()
	lMigrationManager, lError := ormshift.Migrate(
		lDB,
		lDriverDB,
		m001_Create_Table_User{},
		m002_Alter_Table_Usaer_Add_Column_UpdatedAt{},
	)
	if !assertNotNilResultAndNilError(t, lMigrationManager, lError, "ormshift.Migrate") {
		return
	}

	var lSQLExecutor ormshift.SQLExecutor = lDB
	var lSQLBuilder ormshift.SQLBuilder = lDriverDB.SQLBuilder()

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
	if !assertNilError(t, lError, "sqlExecutor.Exec") {
		return
	}
	r, lError := lResult.RowsAffected()
	if !assertNilError(t, lError, "RowsAffected") {
		return
	}
	assertEqualWithLabel(t, 1, r, "RowsAffected")
	i, lError := lResult.LastInsertId()
	if !assertNilError(t, lError, "LastInsertId") {
		return
	}

	//SELECT
	lSQLSelect, lArgs := lSQLBuilder.SelectWithValues(
		"user",
		[]string{"id", "name", "email", "updated_at", "active", "ammount", "percent", "photo"},
		ormshift.ColumnsValues{"id": i},
	)
	lUserRows, lError := lSQLExecutor.Query(lSQLSelect, lArgs...)
	if !assertNotNilResultAndNilError(t, lUserRows, lError, "sqlExecutor.Query") {
		return
	}
	defer lUserRows.Close()
	if !assertEqualWithLabel(t, true, lUserRows.Next(), "Next") {
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
	if !assertNilError(t, lError, "Scan") {
		return
	}
	assertEqualWithLabel(t, i, lUserRow.ID.Int64, "user.id")
	assertEqualWithLabel(t, "Jonh Doe", lUserRow.Name.String, "user.name")
	assertEqualWithLabel(t, "jonh.doe@mail.com", lUserRow.Email.String, "user.email")
	assertEqualWithLabel(t, time.Date(2026, time.January, 9, 12, 15, 52, 100952002, time.UTC), lUserRow.UpdatedAt.Time, "user.updated_at")
	assertEqualWithLabel(t, 1, lUserRow.Active.Byte, "user.active")
	assertEqualWithLabel(t, 5000.00, lUserRow.Ammount.Float64, "user.ammount")
	assertEqualWithLabel(t, 25.54325, lUserRow.Percent.Float64, "user.percent")
	assertEqualWithLabel(t, nil, lUserRow.Photo, "user.photo")
}
