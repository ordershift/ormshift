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
	ConnectionString(params ConnectionParams) string
	SQLBuilder() SQLBuilder
	DBSchema(db *sql.DB) (*schema.DBSchema, error)
}

type Database struct {
	driver           DatabaseDriver
	db               *sql.DB
	connectionString string
	sqlBuilder       SQLBuilder
	dbSchema         *schema.DBSchema
}

func OpenDatabase(driver DatabaseDriver, params ConnectionParams) (*Database, error) {
	if driver == nil {
		return nil, errors.New("DatabaseDriver cannot be nil")
	}
	connectionString := driver.ConnectionString(params)
	db, err := sql.Open(driver.Name(), connectionString)
	if err != nil {
		return nil, fmt.Errorf("sql.Open failed: %w", err)
	}
	dbSchema, err := driver.DBSchema(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get DB schema: %w", err)
	}

	return &Database{
		driver:           driver,
		db:               db,
		connectionString: connectionString,
		sqlBuilder:       driver.SQLBuilder(),
		dbSchema:         dbSchema,
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
