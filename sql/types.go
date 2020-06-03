package sql

//Type is an internal SQL type.
type Type string

//HasType is anything that has a Type.
type HasType interface {
	GetType() Type
}

//GetType implements HasType
func (t Type) GetType() Type {
	return t
}

//Column is the name of a column within a database table.
type Column struct {
	Table
	Name string
}

//HasColumn is anything that has an embedded Column.
type HasColumn interface {
	GetColumn() Column
}

type settableColumn interface {
	setColumn(Column)
}

//GetColumn implements HasColumn
func (c Column) GetColumn() Column {
	return c
}

//setColumn implements HasColumn
func (c *Column) setColumn(to Column) {
	*c = to
}

type Value interface {
	HasType
	HasColumn
	Interface() interface{}
}

//Int is an SQL representation of a Go int.
type Int struct {
	Column
	int
}

func NewInt(i int) Int {
	return Int{
		int: i,
	}
}

//Set the Int to the given int value.
func (i *Int) Set(v int) {
	i.int = v
}

//Get the Int's int value.
func (i Int) Get() int {
	return i.int
}

//Interface gets the Int's int value as an interface.
func (i Int) Interface() interface{} {
	return i.int
}

//GetType implements AnyType.
func (Int) GetType() Type {
	return "int"
}

//String is an SQL representation of a Go string.
type String struct {
	Column
	string
}

func NewString(s string) String {
	return String{
		string: s,
	}
}

//Set the String to the given string value.
func (s *String) Set(v string) {
	s.string = v
}

//Get the String's string value.
func (s String) Get() string {
	return s.string
}

func (s String) Equals(literal string) Condition {
	var c Condition
	c.WriteByte('"')
	c.WriteString(s.Column.Name)
	c.WriteByte('"')
	c.WriteByte('=')
	c.WriteString(c.value(literal))
	return c
}

func (s String) To(literal string) Update {
	return Update{
		Column: s.Column,
		Value:  literal,
	}
}

//Interface gets the Int's int value as an interface.
func (s String) Interface() interface{} {
	return s.string
}

//GetType implements AnyType.
func (String) GetType() Type {
	return "text"
}

//Column is a sql column.
/*type Column interface {
	Type
}

//Type is any sql type.
type Type interface {
	Type() NewType
	Name() string
	String() string
	Default() string
}

//NewType can be used as an embedding to create new types.
type NewType struct {
	string
}

//Name returns the name of the Entry.
func (t NewType) Name() string {
	return t.string
}

//Null a null value.
func (t NewType) Null() Value {
	return Value{
		key:   t.string,
		value: "NULL",
	}
}

//Default value.
func (t NewType) Default() string {
	return "NULL"
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

//Default string.
func (s String) Default() string {
	return `''`
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

//Default string.
func (t Text) Default() string {
	return `''`
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

//Like returns an equality condition on this column.
func (t Text) Like(b string) Condition {
	var c Condition
	fmt.Fprintf(&c, "%v LIKE %v ", strconv.Quote(t.string), c.value(b))
	return c
}

//Boolean is an sql 'boolean'
type Boolean struct {
	NewType
}

func (Boolean) String() string {
	return "boolean"
}

func (Boolean) Default() string {
	return "FALSE"
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
	if v {
		fmt.Fprintf(&c, "%v=TRUE ", strconv.Quote(b.string))
	} else {
		fmt.Fprintf(&c, "%v=FALSE ", strconv.Quote(b.string))
	}
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

//NotEquals returns an inequality condition on this column.
func (s Serial) NotEquals(b int) Condition {
	var c Condition
	fmt.Fprintf(&c, "%v!=%v ", strconv.Quote(s.string), b)
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
}*/
