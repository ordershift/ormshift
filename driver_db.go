package ormshift

import (
	// Import for PostgreSQL
	_ "github.com/lib/pq"
	// Import for SQLite
	_ "github.com/mattn/go-sqlite3"
	// Import for SQL Server
	_ "github.com/microsoft/go-mssqldb"
)

type ConnectionParams struct {
	Host     string
	Instance string
	Port     uint
	User     string
	Password string
	DBname   string
	InMemory bool
}

type DriverDB int

const (
	DriverSQLServer DriverDB = iota
	DriverSQLite
	DriverPostgresql
)

const invalid string = "<<invalid>>"

func (d DriverDB) Name() string {
	switch d {
	case DriverSQLServer:
		return "sqlserver"
	case DriverSQLite:
		return "sqlite3"
	case DriverPostgresql:
		return "postgres"
	default:
		return invalid
	}
}

func (d DriverDB) ConnectionString(pParams ConnectionParams) string {
	switch d {
	case DriverSQLServer:
		return d.sqlServerConnectionString(pParams)
	case DriverSQLite:
		return d.sqliteConnectionString(pParams)
	case DriverPostgresql:
		return d.postgresqlConnectionString(pParams)
	default:
		return invalid
	}
}

func (d DriverDB) IsValid() bool {
	return d.Name() != invalid
}

func (d DriverDB) SQLBuilder() SQLBuilder {
	switch d {
	case DriverSQLServer:
		return sqlserverSQLBuilder{}
	case DriverSQLite:
		return sqliteSQLBuilder{}
	case DriverPostgresql:
		return postgresqlSQLBuilder{}
	default:
		return genericSQLBuilder{}
	}
}
