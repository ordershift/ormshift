package ormshift

import (
	"fmt"
)

func (d DriverDB) sqliteConnectionString(pParams ConnectionParams) string {
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
	return fmt.Sprintf("file:%s.db?%s_locking=NORMAL", pParams.DBname, lConnetionWithAuth)
}
