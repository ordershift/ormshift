package schema

import (
	"fmt"
	"strings"
)

type ConstraintType int

const (
	_ ConstraintType = iota
	ConstraintTypePrimaryKey
	ConstraintTypeForeignKey
	ConstraintTypeUnique
)

type Constraint struct {
	name           string
	constraintType ConstraintType
	columns        []string
}

func NewPrimaryKey(pTableName string, pColumns []string) Constraint {
	return Constraint{
		name:           fmt.Sprintf("PK_%s", pTableName),
		constraintType: ConstraintTypePrimaryKey,
		columns:        pColumns,
	}
}

func NewForeignKey(pFromTableName, pToTableName string, pColumns []string) Constraint {
	return Constraint{
		name:           fmt.Sprintf("FK_%s_%s", pFromTableName, pToTableName),
		constraintType: ConstraintTypeForeignKey,
		columns:        pColumns,
	}
}

func NewUniqueConstraint(pTableName string, pColumns []string) Constraint {
	return Constraint{
		name:           fmt.Sprintf("UQ_%s_%s", pTableName, strings.Join(pColumns, "_")),
		constraintType: ConstraintTypeUnique,
		columns:        pColumns,
	}
}

func (c *Constraint) Name() string {
	return c.name
}

func (c *Constraint) Type() ConstraintType {
	return c.constraintType
}

func (c *Constraint) Columns() []string {
	return c.columns
}
