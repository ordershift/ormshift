package ormshift

type ConnectionParams struct {
	Host     string
	Instance string
	Port     uint16
	User     string
	Password string
	Database string
	InMemory bool
}
