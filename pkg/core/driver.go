package core

type ConnectionParams struct {
	Host     string
	Instance string
	Port     uint
	User     string
	Password string
	DBname   string
	InMemory bool
}
