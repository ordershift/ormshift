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
	_ = tbl.HasPrimaryKey("id")
	err := tbl.HasPrimaryKey("id")
	if !testutils.AssertNotNilError(t, err, "Table.HasPrimaryKey when already set") {
		return
	}
	testutils.AssertErrorMessage(t, "primary key already set for table \"t\"", err, "Table.PrimaryKey")
}

func TestPrimaryKeyFailsWhenColumnDoesNotExist(t *testing.T) {
	tbl := schema.NewTable("t")
	_ = tbl.AddColumns(schema.NewColumnParams{Name: "id", Type: schema.Integer})
	err := tbl.HasPrimaryKey("missing")
	if !testutils.AssertNotNilError(t, err, "Table.HasPrimaryKey when column missing") {
		return
	}
	testutils.AssertErrorMessage(t, "primary key column \"missing\" does not exist in table \"t\"", err, "Table.PrimaryKey")
}

func TestAddForeignKeyFailsWhenColumnDoesNotExist(t *testing.T) {
	tbl := schema.NewTable("t")
	_ = tbl.AddColumns(schema.NewColumnParams{Name: "id", Type: schema.Integer})
	err := tbl.HasForeignKey([]string{"ref_id"}, "other", []string{"id"})
	if !testutils.AssertNotNilError(t, err, "Table.HasForeignKey when column missing") {
		return
	}
	testutils.AssertErrorMessage(t, "foreign key column \"ref_id\" does not exist in table \"t\"", err, "Table.HasForeignKey")
}

func TestAddUniqueConstraintFailsWhenColumnDoesNotExist(t *testing.T) {
	tbl := schema.NewTable("t")
	_ = tbl.AddColumns(schema.NewColumnParams{Name: "id", Type: schema.Integer})
	err := tbl.HasUniqueConstraint("missing")
	if !testutils.AssertNotNilError(t, err, "Table.HasUniqueConstraint when column missing") {
		return
	}
	testutils.AssertErrorMessage(t, "unique constraint column \"missing\" does not exist in table \"t\"", err, "Table.HasUniqueConstraint")
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
	err = tbl.HasUniqueConstraint("email")
	if !testutils.AssertNilError(t, err, "Table.HasUniqueConstraint") {
		return
	}
	ucs := tbl.UniqueConstraints()
	if len(ucs) != 1 {
		t.Fatalf("expected 1 unique constraint, got %d", len(ucs))
	}
	uc := &ucs[0]
	testutils.AssertEqualWithLabel(t, "UC_user_email", uc.Name(), "UniqueConstraint.Name")
	testutils.AssertEqualWithLabel(t, 1, len(uc.Columns()), "UniqueConstraint.Columns length")
	testutils.AssertEqualWithLabel(t, "email", uc.Columns()[0], "UniqueConstraint.Columns[0]")
}
