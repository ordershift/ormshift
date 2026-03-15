package schema_test

import (
	"fmt"
	"testing"

	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/schema"
)

func TestAddColumnFailsWhenAlreadyExists(t *testing.T) {
	productAttributeTable := testutils.FakeProductAttributeTable(t)
	err := productAttributeTable.AddColumns(schema.NewColumnParams{
		Name: "value",
		Type: schema.Integer,
	})
	if !testutils.AssertNotNilError(t, err, "Table.AddColumns") {
		return
	}
	testutils.AssertErrorMessage(t, fmt.Sprintf("failed to add column %q in table %q: column already exists", "value", "product_attribute"), err, "Table.AddColumns")
}

func TestPrimaryKeyFailsWhenAlreadySet(t *testing.T) {
	tbl := schema.NewTable("t")
	_ = tbl.AddColumns(schema.NewColumnParams{Name: "id", Type: schema.Integer})
	_ = tbl.PrimaryKey("id")
	err := tbl.PrimaryKey("id")
	if !testutils.AssertNotNilError(t, err, "Table.PrimaryKey when already set") {
		return
	}
	testutils.AssertErrorMessage(t, "primary key already set for table \"t\"", err, "Table.PrimaryKey")
}

func TestPrimaryKeyFailsWhenColumnDoesNotExist(t *testing.T) {
	tbl := schema.NewTable("t")
	_ = tbl.AddColumns(schema.NewColumnParams{Name: "id", Type: schema.Integer})
	err := tbl.PrimaryKey("missing")
	if !testutils.AssertNotNilError(t, err, "Table.PrimaryKey when column missing") {
		return
	}
	testutils.AssertErrorMessage(t, "primary key column \"missing\" does not exist in table \"t\"", err, "Table.PrimaryKey")
}

func TestAddForeignKeyFailsWhenColumnDoesNotExist(t *testing.T) {
	tbl := schema.NewTable("t")
	_ = tbl.AddColumns(schema.NewColumnParams{Name: "id", Type: schema.Integer})
	err := tbl.AddForeignKey([]string{"ref_id"}, "other", []string{"id"})
	if !testutils.AssertNotNilError(t, err, "Table.AddForeignKey when column missing") {
		return
	}
	testutils.AssertErrorMessage(t, "foreign key column \"ref_id\" does not exist in table \"t\"", err, "Table.AddForeignKey")
}

func TestAddUniqueConstraintFailsWhenColumnDoesNotExist(t *testing.T) {
	tbl := schema.NewTable("t")
	_ = tbl.AddColumns(schema.NewColumnParams{Name: "id", Type: schema.Integer})
	err := tbl.AddUniqueConstraint("missing")
	if !testutils.AssertNotNilError(t, err, "Table.AddUniqueConstraint when column missing") {
		return
	}
	testutils.AssertErrorMessage(t, "unique constraint column \"missing\" does not exist in table \"t\"", err, "Table.AddUniqueConstraint")
}

func TestAddUniqueConstraintAndUCs(t *testing.T) {
	tbl := schema.NewTable("user")
	err := tbl.AddColumns(
		schema.NewColumnParams{Name: "id", Type: schema.Integer},
		schema.NewColumnParams{Name: "email", Type: schema.Varchar, Size: 80},
	)
	if !testutils.AssertNilError(t, err, "Table.AddColumns") {
		return
	}
	err = tbl.AddUniqueConstraint("email")
	if !testutils.AssertNilError(t, err, "Table.AddUniqueConstraint") {
		return
	}
	ucs := tbl.UCs()
	if len(ucs) != 1 {
		t.Fatalf("expected 1 unique constraint, got %d", len(ucs))
	}
	uc := &ucs[0]
	testutils.AssertEqualWithLabel(t, "UC_user_email", uc.Name(), "UniqueConstraint.Name")
	testutils.AssertEqualWithLabel(t, 1, len(uc.Columns()), "UniqueConstraint.Columns length")
	testutils.AssertEqualWithLabel(t, "email", uc.Columns()[0], "UniqueConstraint.Columns[0]")
}
