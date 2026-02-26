package ormshift_test

import (
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestColumnsValuesToNamedArgs(t *testing.T) {
	values := ormshift.ColumnsValues{"id": 1, "sku": "ABC1234", "active": true}
	args := values.ToNamedArgs()
	testutils.AssertEqualWithLabel(t, 3, len(args), "ColumnsValues.ToNamedArgs")
	testutils.AssertEqualWithLabel(t, args[0].Name, "active", "ColumnsValues.ToNamedArgs[0].Name")
	testutils.AssertEqualWithLabel(t, args[0].Value, true, "ColumnsValues.ToNamedArgs[0].Value")
	testutils.AssertEqualWithLabel(t, args[1].Name, "id", "ColumnsValues.ToNamedArgs[1].Name")
	testutils.AssertEqualWithLabel(t, args[1].Value, 1, "ColumnsValues.ToNamedArgs[1].Value")
	testutils.AssertEqualWithLabel(t, args[2].Name, "sku", "ColumnsValues.ToNamedArgs[2].Name")
	testutils.AssertEqualWithLabel(t, args[2].Value, "ABC1234", "ColumnsValues.ToNamedArgs[2].Value")
}

func TestColumnsValuesToColumns(t *testing.T) {
	values := ormshift.ColumnsValues{"id": 1, "sku": "ABC1234"}
	columns := values.ToColumns()
	testutils.AssertEqualWithLabel(t, 2, len(columns), "ColumnsValues.ToColumns")
	testutils.AssertEqualWithLabel(t, columns[0], "id", "ColumnsValues.ToColumns[0]")
	testutils.AssertEqualWithLabel(t, columns[1], "sku", "ColumnsValues.ToColumns[1]")
}

func TestNewColumnsValues(t *testing.T) {
	columns := []string{"id", "sku"}
	values := []any{1, "ABC1234"}
	cv, err := ormshift.NewColumnsValues(columns, values)
	if !testutils.AssertNotNilResultAndNilError(t, cv, err, "ormshift.NewColumnsValues") {
		return
	}
	testutils.AssertEqualWithLabel(t, 2, len(*cv), "len(ColumnsValues)")
	testutils.AssertEqualWithLabel(t, (*cv)["id"], 1, "id")
	testutils.AssertEqualWithLabel(t, (*cv)["sku"], "ABC1234", "sku")
}

func TestNewColumnsValuesFailsWhenIncompatibleLengths(t *testing.T) {
	columns := []string{"id", "sku"}
	values := []any{1}
	cv, err := ormshift.NewColumnsValues(columns, values)
	if !testutils.AssertNilResultAndNotNilError(t, cv, err, "ormshift.NewColumnsValues") {
		return
	}
	testutils.AssertErrorMessage(t, "failed to get columns values: columns len must be equal to values len", err, "ormshift.NewColumnsValues")
}
