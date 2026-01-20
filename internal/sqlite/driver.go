package sqlite

import (
	"fmt"

	"github.com/ordershift/ormshift"
	_ "modernc.org/sqlite"
)

func DriverName() string {
	return "sqlite"
}

func ConnectionString(pParams ormshift.ConnectionParams) string {
	if pParams.InMemory {
		return ":memory:"
	}
	lConnetionWithAuth := ""
	if pParams.User != "" {
		lConnetionWithAuth += fmt.Sprintf("_auth&_auth_user=%s&", pParams.User)
		if pParams.Password != "" {
			lConnetionWithAuth += fmt.Sprintf("_auth_pass=%s&", pParams.Password)
		}
	}
	return fmt.Sprintf("file:%s.db?%s_locking=NORMAL", pParams.Database, lConnetionWithAuth)
}

func SQLBuilder() ormshift.SQLBuilder {
	return sqliteSQLBuilder{}
}
