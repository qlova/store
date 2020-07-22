package db

type Join int

const (
	Inner Join = iota
)

type Linker struct {
	From Column
	Join
	To Column
}

type Connection string

func (conn Connection) String() string {
	if len(conn) == 0 {
		return "Table"
	}
	return string(conn) + " Table"
}

type Table struct {
	Name string

	Connection Connection
}

type SortMode int

const (
	Increasing SortMode = iota
	Decreasing
)

type Column struct {
	Table
	Name string

	SortMode

	Field      int16
	references *Column

	Key bool

	nullable, null bool
}

func (c Column) GetColumn() Column {
	return c
}

func (c *Column) SetColumn(to Column) {
	*c = to
}

func (c Column) Decreasing() Column {
	c.SortMode = Decreasing
	return c
}

func (c Column) FieldName() string {
	return c.Name
}

type Update struct {
	Column
	Value interface{}
}
