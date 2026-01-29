package sqlite_test

import (
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/dialects/sqlite"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestHasTable(t *testing.T) {
	db, err := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNotNilResultAndNilError(t, db, err, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = db.Close() }()

	if db.DBSchema().HasTable("user") {
		t.Errorf("Expected HasTable('user') to be false")
	}
}

func TestHasColumn(t *testing.T) {
	db, err := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
	if !testutils.AssertNotNilResultAndNilError(t, db, err, "ormshift.OpenDatabase") {
		return
	}
	defer func() { _ = db.Close() }()

	if db.DBSchema().HasColumn("user", "id") {
		t.Errorf("Expected HasColumn('user', 'id') to be false")
	}
}
