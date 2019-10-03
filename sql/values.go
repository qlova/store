package sql

//Values is a collection of SQL values.
type Values []Value

//Value is a SQL value packed and ready for a query.
type Value struct {
	key, value string
	arg        interface{}
}

//Column returns the column of this value.
func (value Value) Column() string {
	return value.key
}

func (value Value) get(q Query) string {
	if value.arg != nil {
		return q.value(value.arg)
	}
	return value.value
}
