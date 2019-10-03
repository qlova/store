package sql

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

//Query is a sql query.
type Query struct {
	*query
}

type query struct {
	Database
	bytes.Buffer

	args  []interface{}
	error error
}

//NewQuery returns a new query.
func (db Database) NewQuery() Query {
	return Query{
		&query{
			Database: db,
		},
	}
}

//QueryError returns a new query that is invalid.
func (db Database) QueryError(msg string) Query {
	return Query{
		&query{
			error: errors.New(msg),
		},
	}
}

//Do executes the query.
func (q Query) Do() Result {
	if q.error != nil {
		return Result{error: q.error}
	}

	result, err := q.Query(q.Buffer.String(), q.args...)
	return Result{q, result, err}
}

func (q Query) String() string {
	var result string = q.Buffer.String()
	for i, arg := range q.args {
		result = strings.Replace(result, fmt.Sprintf("$%v", i+1), fmt.Sprintf("`%v`", arg), 1)
	}
	return result
}

func (q *Query) value(v interface{}) string {
	q.args = append(q.args, v)
	return fmt.Sprintf("$%v", len(q.args))
}

//Where filters on the condition.
func (q Query) Where(condition Condition) Query {
	fmt.Fprint(q, "WHERE ")
	condition.writeTo(q)
	return q
}

//Orderable is an orderable column.
type Orderable interface {
	Orderable() string
}

//OrderBy orders the results.
func (q Query) OrderBy(orders ...Orderable) Query {
	fmt.Fprint(q, "ORDER BY ")
	if len(orders) == 0 {
		return q
	}
	fmt.Fprint(q, orders[0].Orderable())

	for _, order := range orders[1:] {
		fmt.Fprintf(q, "%v,", order.Orderable())
	}

	return q
}
