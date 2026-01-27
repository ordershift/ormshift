package sqlserver_test

import (
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/sqlserver"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestHasTable(t *testing.T) {
	lDB, lError := ormshift.OpenDatabase(sqlserver.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNotNilResultAndNilError(t, lDB, lError, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = lDB.Close() }()

	if lDB.DBSchema().HasTable("user") {
		t.Errorf("Expected HasTable('user') to be false")
	}
}

func TestHasColumn(t *testing.T) {
	lDB, lError := ormshift.OpenDatabase(sqlserver.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNotNilResultAndNilError(t, lDB, lError, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = lDB.Close() }()

	if lDB.DBSchema().HasColumn("user", "id") {
		t.Errorf("Expected HasColumn('user', 'id') to be false")
	}
}
