package sql

import (
	"fmt"
	"strconv"
	"time"
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
	fmt.Fprintf(&c, "%v=%v", strconv.Quote(i.string), b)
	return c
}

//NotEquals returns an equality condition on this column.
func (i Int) NotEquals(b int) Condition {
	var c Condition
	fmt.Fprintf(&c, "%v!=%v", strconv.Quote(i.string), b)
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
func (s String) Value(v string) Value {
	return Value{
		key: s.string,
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
	fmt.Fprintf(&c, "%v=%v ", strconv.Quote(s.string), c.value(b))
	return c
}

//NotEquals returns an equality condition on this column.
func (s String) NotEquals(b string) Condition {
	var c Condition
	fmt.Fprintf(&c, "%v!=%v ", strconv.Quote(s.string), c.value(b))
	return c
}

//NotNull returns a null condition on this column.
func (s String) NotNull() Condition {
	var c Condition
	fmt.Fprintf(&c, "%v IS NOT NULL ", strconv.Quote(s.string))
	return c
}

//Text is an sql 'text'
type Text struct {
	NewType
}

func (Text) String() string {
	return "text"
}

//Value returns the string as a value.
func (t Text) Value(v string) Value {
	return Value{
		key: t.string,
		arg: v,
	}
}

//Orderable strings are orderable.
func (t Text) Orderable() string {
	return t.string
}

//Equals returns an equality condition on this column.
func (t Text) Equals(b string) Condition {
	var c Condition
	fmt.Fprintf(&c, "%v=%v ", strconv.Quote(t.string), c.value(b))
	return c
}

//NotNull returns a null condition on this column.
func (t Text) NotNull() Condition {
	var c Condition
	fmt.Fprintf(&c, "%v IS NOT NULL ", strconv.Quote(t.string))
	return c
}

//NotEquals returns an equality condition on this column.
func (t Text) NotEquals(b string) Condition {
	var c Condition
	fmt.Fprintf(&c, "%v!=%v ", strconv.Quote(t.string), c.value(b))
	return c
}

//Boolean is an sql 'boolean'
type Boolean struct {
	NewType
}

func (Boolean) String() string {
	return "boolean"
}

//Value returns the string as a value.
func (b Boolean) Value(v bool) Value {
	return Value{
		key: b.string,
		arg: v,
	}
}

//Orderable booleans are orderable.
func (b Boolean) Orderable() string {
	return b.string
}

//Equals returns an equality condition on this column.
func (b Boolean) Equals(v bool) Condition {
	var c Condition
	fmt.Fprintf(&c, "%v=%v ", strconv.Quote(b.string), c.value(v))
	return c
}

//Serial is a sql 'serial'.
type Serial struct {
	NewType
}

func (Serial) String() string {
	return "serial"
}

//Equals returns an equality condition on this column.
func (s Serial) Equals(b int) Condition {
	var c Condition
	fmt.Fprintf(&c, "%v=%v ", strconv.Quote(s.string), b)
	return c
}

//Timestamp is a sql 'timestamp'.
type Timestamp struct {
	NewType
}

func (Timestamp) String() string {
	return "timestamp"
}

//Value returns the times.Time as a timestamp value.
func (t Timestamp) Value(v time.Time) Value {
	return Value{
		key: t.string,
		arg: v,
	}
}

//Equals returns an equality condition on this column.
func (t Timestamp) Equals(b time.Time) Condition {
	var c Condition
	fmt.Fprintf(&c, "%v=%v ", strconv.Quote(t.string), b)
	return c
}
