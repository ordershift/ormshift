package sqlserver

import (
	"fmt"

	_ "github.com/microsoft/go-mssqldb"

	"github.com/ordershift/ormshift/pkg/core"
)

func DriverName() string {
	return "sqlserver"
}

func ConnectionString(pParams core.ConnectionParams) string {
	lHostInstanceAndPort := pParams.Host
	if pParams.Instance != "" {
		lHostInstanceAndPort = fmt.Sprintf("%s\\%s", pParams.Host, pParams.Instance)
	}
	if pParams.Port > 0 {
		lHostInstanceAndPort += fmt.Sprintf(";port=%d", pParams.Port)
	}
	return fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", lHostInstanceAndPort, pParams.User, pParams.Password, pParams.Database)
}

func SQLBuilder() core.SQLBuilder {
	return sqlserverSQLBuilder{}
}
