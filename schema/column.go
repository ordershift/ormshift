package schema

type ColumnType int

const (
	_ ColumnType = iota
	Integer
	Varchar
	Monetary
	DateTime
	Decimal
	Boolean
	Binary
	DateTimeOffSet
)

type NewColumnParams struct {
	Name          string
	Type          ColumnType
	Size          uint
	NotNull       bool
	AutoIncrement bool
	Default       string // SQL expression for DEFAULT (e.g. "0", "'hello'", "now()"). Empty means no default.
	Check         string // SQL expression for CHECK (e.g. "price >= 0"). Empty means no check.
}

type Column struct {
	name          string
	columnType    ColumnType
	size          uint
	notNull       bool
	autoIncrement bool
	defaultExpr   string
	checkExpr     string
}

func NewColumn(params NewColumnParams) Column {
	return Column{
		name:          params.Name,
		columnType:    params.Type,
		size:          params.Size,
		notNull:       params.NotNull,
		autoIncrement: params.AutoIncrement,
		defaultExpr:   params.Default,
		checkExpr:     params.Check,
	}
}

func (c *Column) Name() string {
	return c.name
}

func (c *Column) Type() ColumnType {
	return c.columnType
}

func (c *Column) Size() uint {
	return c.size
}

func (c *Column) NotNull() bool {
	return c.notNull
}

func (c *Column) AutoIncrement() bool {
	return c.autoIncrement
}

func (c *Column) Default() string {
	return c.defaultExpr
}

func (c *Column) Check() string {
	return c.checkExpr
}
