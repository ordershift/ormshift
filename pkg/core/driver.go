package core

type ConnectionParams struct {
	Host     string
	Instance string
	Port     uint
	User     string
	Password string
	Database string
	InMemory bool
}
