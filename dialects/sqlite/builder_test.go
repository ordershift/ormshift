package sqlite_test

import (
	"database/sql"
	"testing"

	"github.com/ordershift/ormshift/dialects/sqlite"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestInteroperateSQLCommandWithNamedArgs(t *testing.T) {
	lDriver := &sqlite.SQLiteDriver{}
	lReturnedSQL, lReturnedValue := lDriver.SQLBuilder().InteroperateSQLCommandWithNamedArgs("select * from table where id = @id", sql.Named("id", 1))
	lExpectedSQL := "select * from table where id = @id"
	testutils.AssertEqualWithLabel(t, lExpectedSQL, lReturnedSQL, lDriver.Name()+".InteroperateSQLCommandWithNamedArgs.SQL")
	testutils.AssertEqualWithLabel(t, "id", lReturnedValue[0].(sql.NamedArg).Name, lDriver.Name()+".InteroperateSQLCommandWithNamedArgs.Name")
	testutils.AssertEqualWithLabel(t, 1, lReturnedValue[0].(sql.NamedArg).Value.(int), lDriver.Name()+".InteroperateSQLCommandWithNamedArgs.Value")
}
