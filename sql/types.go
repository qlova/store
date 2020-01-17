package sql

import (
	"fmt"
	"strconv"
)

//Column is a sql column.
type Column interface {
	Type
}

//Type is any sql type.
type Type interface {
	Type() NewType
	Name() string
	String() string
}

//NewType can be used as an embedding to create new types.
type NewType struct {
	string
}

//Name returns the name of the Entry.
func (t NewType) Name() string {
	return t.string
}

//Type is
func (t NewType) Type() NewType {
	return t
}

var star = starType{NewType{"*"}}

type starType struct {
	NewType
}

func (starType) String() string {
	return "*"
}

//Int is a sql 'int'.
type Int struct {
	NewType
}

func (Int) String() string {
	return "int"
}

//Value returns the int as a value.
func (i Int) Value(v int) Value {
	return Value{
		key:   i.string,
		value: strconv.Itoa(v),
	}
}

//Equals returns an equality condition on this column.
func (i Int) Equals(b int) Condition {
	var c Condition
	fmt.Fprintf(&c, "%v=%v", i.string, b)
	return c
}

//String is a sql 'varchar(255)'
type String struct {
	NewType
	length int
}

func (String) String() string {
	return "varchar(255)"
}

//Value returns the string as a value.
func (i String) Value(v string) Value {
	return Value{
		key: i.string,
		arg: v,
	}
}

//Orderable strings are orderable.
func (s String) Orderable() string {
	return s.string
}

//Equals returns an equality condition on this column.
func (s String) Equals(b string) Condition {
	var c Condition
	fmt.Fprintf(&c, "%v=%v", s.string, c.value(b))
	return c
}
