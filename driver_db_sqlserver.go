package ormshift

import (
	"fmt"
)

func (d DriverDB) sqlServerConnectionString(pParams ConnectionParams) string {
	lHostInstanceAndPort := pParams.Host
	if pParams.Instance != "" {
		lHostInstanceAndPort = fmt.Sprintf("%s\\%s", pParams.Host, pParams.Instance)
	}
	if pParams.Port > 0 {
		lHostInstanceAndPort += fmt.Sprintf(";port=%d", pParams.Port)
	}
	return fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", lHostInstanceAndPort, pParams.User, pParams.Password, pParams.DBname)
}
