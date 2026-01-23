package testutils

import (
	"database/sql"
	"fmt"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/schema"
)

// FakeDriverInvalidConnectionString always returns an empty connection string.
type FakeDriverInvalidConnectionString struct {
	underlyingDriver ormshift.DatabaseDriver
}

func NewFakeDriverInvalidConnectionString(underlyingDriver ormshift.DatabaseDriver) *FakeDriverInvalidConnectionString {
	return &FakeDriverInvalidConnectionString{
		underlyingDriver: underlyingDriver,
	}
}
func (d *FakeDriverInvalidConnectionString) ConnectionString(params ormshift.ConnectionParams) string {
	return "invalid-connection-string"
}
func (d *FakeDriverInvalidConnectionString) Name() string {
	return d.underlyingDriver.Name()
}
func (d *FakeDriverInvalidConnectionString) SQLBuilder() ormshift.SQLBuilder {
	return d.underlyingDriver.SQLBuilder()
}
func (d *FakeDriverInvalidConnectionString) DBSchema(pDB *sql.DB) (*schema.DBSchema, error) {
	return d.underlyingDriver.DBSchema(pDB)
}

// FakeDriverBadSchema always returns an error when attempting to get the DB schema.
type FakeDriverBadSchema struct {
	underlyingDriver ormshift.DatabaseDriver
}

func NewFakeDriverBadSchema(underlyingDriver ormshift.DatabaseDriver) *FakeDriverBadSchema {
	return &FakeDriverBadSchema{
		underlyingDriver: underlyingDriver,
	}
}
func (d *FakeDriverBadSchema) ConnectionString(params ormshift.ConnectionParams) string {
	return d.underlyingDriver.ConnectionString(params)
}
func (d *FakeDriverBadSchema) Name() string {
	return d.underlyingDriver.Name()
}
func (d *FakeDriverBadSchema) SQLBuilder() ormshift.SQLBuilder {
	return d.underlyingDriver.SQLBuilder()
}
func (d *FakeDriverBadSchema) DBSchema(pDB *sql.DB) (*schema.DBSchema, error) {
	return nil, fmt.Errorf("intentionally bad schema")
}

// FakeDriverBadName always returns an invalid name.
type FakeDriverBadName struct {
	underlyingDriver ormshift.DatabaseDriver
}

func NewFakeDriverBadName(underlyingDriver ormshift.DatabaseDriver) *FakeDriverBadName {
	return &FakeDriverBadName{
		underlyingDriver: underlyingDriver,
	}
}
func (d *FakeDriverBadName) ConnectionString(params ormshift.ConnectionParams) string {
	return d.underlyingDriver.ConnectionString(params)
}
func (d *FakeDriverBadName) Name() string {
	return "bad-driver-name"
}
func (d *FakeDriverBadName) SQLBuilder() ormshift.SQLBuilder {
	return d.underlyingDriver.SQLBuilder()
}
func (d *FakeDriverBadName) DBSchema(pDB *sql.DB) (*schema.DBSchema, error) {
	return d.underlyingDriver.DBSchema(pDB)
}
