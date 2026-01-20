package ormshift

import (
	"database/sql"

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
