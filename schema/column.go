package schema

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
	AutoIncrement bool
}

type Column struct {
	name          string
	columnType    ColumnType
	size          uint
	primaryKey    bool
	notNull       bool
	autoIncrement bool
}

func NewColumn(params NewColumnParams) Column {
	return Column{
		name:          params.Name,
		columnType:    params.Type,
		size:          params.Size,
		primaryKey:    params.PrimaryKey,
		notNull:       params.NotNull,
		autoIncrement: params.AutoIncrement,
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

func (c *Column) PrimaryKey() bool {
	return c.primaryKey
}

func (c *Column) NotNull() bool {
	return c.notNull
}

func (c *Column) AutoIncrement() bool {
	return c.autoIncrement
}
