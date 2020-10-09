package db

//Expression containing a column.
type Expression struct {
	Table, Column string
	Operator
	Value interface{}

	View Table

	driver Driver
}
