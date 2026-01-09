package ormshift_test

import (
	"testing"

	"github.com/ordershift/ormshift"
)

func Test_DriverSQLServer_Name_ShouldBeValid(t *testing.T) {
	assertEqualWithLabel(t, "sqlserver", ormshift.DriverSQLServer.Name(), "DriverSQLServer.Name")
	assertEqualWithLabel(t, true, ormshift.DriverSQLServer.IsValid(), "DriverSQLServer.IsValid")
}

func Test_DriverSQLite_Name_ShouldBeValid(t *testing.T) {
	assertEqualWithLabel(t, "sqlite3", ormshift.DriverSQLite.Name(), "DriverSQLite.Name")
	assertEqualWithLabel(t, true, ormshift.DriverSQLite.IsValid(), "DriverSQLite.IsValid")
}

func Test_DriverPostgresql_Name_ShouldBeValid(t *testing.T) {
	assertEqualWithLabel(t, "postgres", ormshift.DriverPostgresql.Name(), "DriverPostgresql.Name")
	assertEqualWithLabel(t, true, ormshift.DriverPostgresql.IsValid(), "DriverPostgresql.IsValid")
}

func Test_DriverInvalid_Name_ShouldBeInvalid(t *testing.T) {
	assertEqualWithLabel(t, "<<invalid>>", ormshift.DriverDB(-1).Name(), "DriverInvalid.Name")
	assertEqualWithLabel(t, false, ormshift.DriverDB(-1).IsValid(), "DriverInvalid.IsValid")
}

func Test_DriverSQLServer_ConnectionString_ShouldBeValid(t *testing.T) {
	lReturnedConnectionString := ormshift.DriverSQLServer.ConnectionString(ormshift.ConnectionParams{
		Host:     "my-server",
		Port:     1433,
		Instance: "sqlexpress",
		User:     "sa",
		Password: "123456",
		DBname:   "my-db",
	})
	lExpectedConnectionString := "server=my-server\\sqlexpress;port=1433;user id=sa;password=123456;database=my-db" //NOSONAR go:S2068
	assertEqualWithLabel(t, lExpectedConnectionString, lReturnedConnectionString, "DriverSQLServer.ConnectionString")
}

func Test_DriverSQLite_ConnectionString_ShouldBeValid(t *testing.T) {
	lReturnedConnectionString := ormshift.DriverSQLite.ConnectionString(ormshift.ConnectionParams{
		User:     "user",
		Password: "123456",
		DBname:   "my-db",
	})
	lExpectedConnectionString := "file:my-db.db?_auth&_auth_user=user&_auth_pass=123456&_locking=NORMAL"
	assertEqualWithLabel(t, lExpectedConnectionString, lReturnedConnectionString, "DriverSQLite.ConnectionString")
}

func Test_DriverSQLite_ConnectionString_ShouldBeValid_WhenInMemory(t *testing.T) {
	lReturnedConnectionString := ormshift.DriverSQLite.ConnectionString(ormshift.ConnectionParams{InMemory: true})
	lExpectedConnectionString := ":memory:"
	assertEqualWithLabel(t, lExpectedConnectionString, lReturnedConnectionString, "DriverSQLite.ConnectionString")
}

func Test_DriverPostgresql_ConnectionString_ShouldBeValid(t *testing.T) {
	lReturnedConnectionString := ormshift.DriverPostgresql.ConnectionString(ormshift.ConnectionParams{
		User:     "pg",
		Password: "123456",
		DBname:   "my-db",
	})
	lExpectedConnectionString := "host=localhost port=5432 user=pg password=123456 dbname=my-db sslmode=disable" //NOSONAR go:S2068
	assertEqualWithLabel(t, lExpectedConnectionString, lReturnedConnectionString, "DriverPostgresql.ConnectionString")
}

func Test_DriverInvalid_ConnectionString_ShouldBeInvalid(t *testing.T) {
	lReturnedConnectionString := ormshift.DriverDB(-1).ConnectionString(ormshift.ConnectionParams{
		User:     "pg",
		Password: "123456",
		DBname:   "my-db",
	})
	lExpectedConnectionString := "<<invalid>>"
	assertEqualWithLabel(t, lExpectedConnectionString, lReturnedConnectionString, "DriverInvalid.ConnectionString")
}
