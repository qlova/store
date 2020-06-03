package sql

import (
	"strings"
)

type Condition struct {
	strings.Builder
	Values []interface{}
}

func (c *Condition) value(i interface{}) string {
	c.Values = append(c.Values, i)
	return "$%v"
}

//Condition can be used in where queries.
/*type Condition struct {
	bytes.Buffer
	args []interface{}
}

//False is a condition that evaluates to false.
var False = Condition{
	Buffer: *bytes.NewBuffer([]byte("false")),
}

//True is a condition that evaluates to false.
var True = Condition{
	Buffer: *bytes.NewBuffer([]byte("true")),
}

func (c *Condition) value(v interface{}) string {
	c.args = append(c.args, v)
	return fmt.Sprintf("$%v", len(c.args))
}

func (c Condition) writeTo(q Query) {
	var b = c.Bytes()
	for i, arg := range c.args {
		b = bytes.Replace(b, []byte(fmt.Sprintf("$%v", i+1)), []byte(q.value(arg)), 1)
	}
	q.Write(b)
}

//And does an and operation on the two conditions.
func (c Condition) And(other Condition) Condition {
	var r Condition
	r.Write(c.Bytes())
	r.WriteString(" AND ")
	r.args = c.args

	var b = other.Bytes()
	for i, arg := range other.args {
		b = bytes.Replace(b, []byte(fmt.Sprintf("$%v", i+1)), []byte(r.value(arg)), 1)
	}
	r.Write(b)

	return r
}

//Or does an or operation on the two conditions.
func (c Condition) Or(other Condition) Condition {
	var r Condition
	r.Write(c.Bytes())
	r.WriteString(" Or ")
	r.args = c.args

	var b = other.Bytes()
	for i, arg := range other.args {
		b = bytes.Replace(b, []byte(fmt.Sprintf("$%v", i+1)), []byte(r.value(arg)), 1)
	}
	r.Write(b)

	return r
}

//Not does a not operation on the two conditions.
func (c Condition) Not() Condition {
	var r Condition
	r.WriteString("NOT ")
	r.Write(c.Bytes())
	r.args = c.args
	return r
}*/
