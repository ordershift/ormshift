package ormshift

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ordershift/ormshift/schema"
)

type ConnectionParams struct {
	Host     string
	Instance string
	Port     uint16
	User     string
	Password string
	Database string
	InMemory bool
}

type DatabaseDriver interface {
	Name() string
	ConnectionString(pParams ConnectionParams) string
	SQLBuilder() SQLBuilder
	DBSchema(pDB *sql.DB) (*schema.DBSchema, error)
}

type Database struct {
	driver           DatabaseDriver
	db               *sql.DB
	connectionString string
	sqlBuilder       SQLBuilder
	dbSchema         *schema.DBSchema
}

func OpenDatabase(pDriver DatabaseDriver, pParams ConnectionParams) (*Database, error) {
	if pDriver == nil {
		return nil, errors.New("DatabaseDriver cannot be nil")
	}
	lConnectionString := pDriver.ConnectionString(pParams)
	lDB, lError := sql.Open(pDriver.Name(), lConnectionString)
	if lError != nil {
		return nil, fmt.Errorf("sql.Open failed: %w", lError)
	}
	lDBSchema, lError := pDriver.DBSchema(lDB)
	if lError != nil {
		return nil, fmt.Errorf("failed to get DB schema: %w", lError)
	}

	return &Database{
		driver:           pDriver,
		db:               lDB,
		connectionString: lConnectionString,
		sqlBuilder:       pDriver.SQLBuilder(),
		dbSchema:         lDBSchema,
	}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) DB() *sql.DB {
	return d.db
}

func (d *Database) SQLExecutor() SQLExecutor {
	return d.db
}

func (d *Database) DriverName() string {
	return d.driver.Name()
}

func (d *Database) ConnectionString() string {
	// ConnectionParams is mutable for simplicity, so we store the connection string used to open the database
	return d.connectionString
}

func (d *Database) SQLBuilder() SQLBuilder {
	return d.sqlBuilder
}

func (d *Database) DBSchema() *schema.DBSchema {
	return d.dbSchema
}
