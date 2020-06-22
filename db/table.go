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

type Column struct {
	Table
	Name string

	Field      int16
	references *Column

	nullable, null bool
}

func (c Column) GetColumn() Column {
	return c
}

func (c *Column) SetColumn(to Column) {
	*c = to
}

type Update struct {
	Column
	Value interface{}
}
