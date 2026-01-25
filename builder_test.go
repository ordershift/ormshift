package ormshift_test

import (
	"testing"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/internal/testutils"
)

func TestColumnsValuesToNamedArgs(t *testing.T) {
	lColumnsValues := ormshift.ColumnsValues{"id": 1, "sku": "ABC1234", "active": true}
	lNamedArgs := lColumnsValues.ToNamedArgs()
	testutils.AssertEqualWithLabel(t, 3, len(lNamedArgs), "ColumnsValues.ToNamedArgs")
	testutils.AssertEqualWithLabel(t, lNamedArgs[0].Name, "active", "ColumnsValues.ToNamedArgs[0].Name")
	testutils.AssertEqualWithLabel(t, lNamedArgs[0].Value, true, "ColumnsValues.ToNamedArgs[0].Value")
	testutils.AssertEqualWithLabel(t, lNamedArgs[1].Name, "id", "ColumnsValues.ToNamedArgs[1].Name")
	testutils.AssertEqualWithLabel(t, lNamedArgs[1].Value, 1, "ColumnsValues.ToNamedArgs[1].Value")
	testutils.AssertEqualWithLabel(t, lNamedArgs[2].Name, "sku", "ColumnsValues.ToNamedArgs[2].Name")
	testutils.AssertEqualWithLabel(t, lNamedArgs[2].Value, "ABC1234", "ColumnsValues.ToNamedArgs[2].Value")
}

func TestColumnsValuesToColumns(t *testing.T) {
	lColumnsValues := ormshift.ColumnsValues{"id": 1, "sku": "ABC1234"}
	lColumns := lColumnsValues.ToColumns()
	testutils.AssertEqualWithLabel(t, 2, len(lColumns), "ColumnsValues.ToColumns")
	testutils.AssertEqualWithLabel(t, lColumns[0], "id", "ColumnsValues.ToColumns[0]")
	testutils.AssertEqualWithLabel(t, lColumns[1], "sku", "ColumnsValues.ToColumns[1]")
}
