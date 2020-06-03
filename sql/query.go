package sql

import (
	"fmt"
	"strconv"
	"strings"
)

type Query struct {
	db Database
	strings.Builder
	Values []interface{}
}

func (db Database) NewQuery() *Query {
	return &Query{db: db}
}

func (db Database) Where(condition Condition) *Query {
	var q Query
	q.db = db
	q.WriteString(`WHERE `)
	q.WriteCondition(condition)
	return &q
}

func (q *Query) value(i interface{}) string {
	q.Values = append(q.Values, i)
	return "$%v"
}

func (q *Query) String() string {
	var indicies = make([]interface{}, len(q.Values))
	for i := range q.Values {
		indicies[i] = i + 1
	}
	return fmt.Sprintf(q.Builder.String(), indicies...)
}

func (q *Query) WriteCondition(condition Condition) {
	q.WriteString(condition.Builder.String())
	q.Values = append(q.Values, condition.Values...)
}

func (q *Query) WriteQuery(other *Query) {
	q.WriteString(other.Builder.String())
	q.Values = append(q.Values, other.Values...)
}

func (q *Query) WriteColumn(column HasColumn) {
	q.WriteString(strconv.Quote(string(column.GetColumn().Name)))
}

//Query is a sql query.
/*type Query struct {
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

	if Debug {
		fmt.Println(q.Buffer.String())
	}

	if q.error != nil {
		return Result{error: q.error}
	}

	result, err := q.Query(q.Buffer.String(), q.args...)

	if result != nil {
		runtime.SetFinalizer(result, func(rows *sql.Rows) {
			rows.Close()
		})
	}

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

//Limit places a limit on the number of results.
func (q Query) Limit(n int) Query {
	fmt.Fprintf(q, "LIMIT %v", n)
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
	fmt.Fprintf(q, `"%v"`, orders[0].Orderable())

	for _, order := range orders[1:] {
		fmt.Fprintf(q, `, "%v"`, order.Orderable())
	}

	return q
}*/
