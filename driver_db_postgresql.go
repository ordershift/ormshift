package ormshift

import (
	"fmt"
)

func (d DriverDB) postgresqlConnectionString(pParams ConnectionParams) string {
	lHost := pParams.Host
	if lHost == "" {
		lHost = "localhost"
	}
	lPorta := pParams.Port
	if lPorta == 0 {
		lPorta = 5432
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", lHost, lPorta, pParams.User, pParams.Password, pParams.DBname)
}
