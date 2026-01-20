package schema

import (
	"fmt"
	"regexp"
)

var regexValidColumnName = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)

type ColumnName struct {
	columnName string
}

func NewColumnName(pName string) (*ColumnName, error) {
	if !regexValidColumnName.MatchString(pName) {
		return nil, fmt.Errorf("invalid column name: %q", pName)
	}
	return &ColumnName{pName}, nil
}

func (tn ColumnName) String() string {
	return tn.columnName
}

type ColumnType int

const (
	Integer ColumnType = iota
	Varchar
	Monetary
	DateTime
	Decimal
	Boolean
	Binary
)

type NewColumnParams struct {
	Name          string
	Type          ColumnType
	Size          uint
	PrimaryKey    bool
	NotNull       bool
	Autoincrement bool
}

type Column struct {
	name       ColumnName
	columnType ColumnType
	size       uint
	pk         bool
	notNull    bool
	autoInc    bool
}

func NewColumn(pParams NewColumnParams) (*Column, error) {
	lColumnName, lError := NewColumnName(pParams.Name)
	if lError != nil {
		return nil, lError
	}
	return &Column{
		name:       *lColumnName,
		columnType: pParams.Type,
		size:       pParams.Size,
		pk:         pParams.PrimaryKey,
		notNull:    pParams.NotNull,
		autoInc:    pParams.Autoincrement,
	}, nil
}

func (c Column) Name() ColumnName {
	return c.name
}

func (c Column) Type() ColumnType {
	return c.columnType
}

func (c Column) Size() uint {
	return c.size
}

func (c Column) PrimaryKey() bool {
	return c.pk
}

func (c Column) NotNull() bool {
	return c.notNull
}

func (c Column) Autoincrement() bool {
	return c.autoInc
}
